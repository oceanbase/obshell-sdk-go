/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sdk

import (
	"fmt"
	"reflect"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/internal/util"
	"github.com/oceanbase/obshell-sdk-go/sdk/auth"
	"github.com/oceanbase/obshell-sdk-go/sdk/option"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	responselib "github.com/oceanbase/obshell-sdk-go/sdk/response"
)

// Client is not thread-safe
type Client struct {
	httpClient *resty.Client
	host       string
	port       int

	auth          auth.Auther
	candidateAuth auth.Auther

	TaskQueue chan func() error // Task queue
	IsSyncing bool              // Is syncing
}

// NewClient creates a new client with the given host and port.
//
// The client will use the given options to configure the client.
// You can use the sdk.WithPasswordAuth option to set the password for the client.
//
// AS: sdk.NewClient("127.0.0.1", 8080, sdk.WithPasswordAuth("password"))
func NewClient(host string, port int, options ...option.Optioner) (c *Client, err error) {
	c, err = NewClientWithServer(host, port)
	if err != nil {
		return
	}

	for _, opt := range options {
		switch opt.Type() {
		case option.AUTH_OPT:
			c.auth = opt.Value().(auth.Auther)
		}
	}
	return
}

func NewClientWithServer(host string, port int) (*Client, error) {
	c := &Client{
		httpClient: resty.New(),
		TaskQueue:  make(chan func() error, 10),
		auth:       auth.NewPasswordAuth(""),
		host:       host,
		port:       port,
	}
	return c, nil
}

func NewClientWithPassword(host string, port int, password string) (c *Client, err error) {
	return NewClient(host, port, WithPasswordAuth(password))
}

func (c *Client) GetServer() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *Client) GetAuth() auth.Auther {
	return c.auth
}

func (c *Client) SetAuth(auth auth.Auther) {
	auth.ResetMethod()
	c.setAuth(auth)
}

func (c *Client) setAuth(auth auth.Auther) {
	c.auth = auth
	c.candidateAuth = nil // Clear candidate auth when an new auth is set
}

func (c *Client) SetCandidateAuth(auth auth.Auther) {
	if auth.IsAutoSelectVersion() && auth.Type() == c.auth.Type() {
		if c.auth.IsAutoSelectVersion() {
			auth.AutoSelectVersion(c.auth.GetVersion())
		} else {
			auth.SetVersion(c.auth.GetVersion())
		}
	}
	c.candidateAuth = auth
}

func (c *Client) AdoptCandidateAuth() {
	if c.candidateAuth == nil {
		return
	}
	c.auth = c.candidateAuth
	c.candidateAuth = nil
}

func (c *Client) DiscardCandidateAuth() {
	c.candidateAuth = nil
}

func (c *Client) GetHost() string { // MASTER or CLUSTER AGENT
	return c.host
}

func (c *Client) GetPort() int {
	return c.port
}

func (c *Client) confirmAuthVersion() error {
	agentInfo, err := util.GetInfo(c.GetServer())
	if err != nil {
		return errors.Wrap(err, "get version failed")
	}

	if !c.auth.IsAutoSelectVersion() {
		if !c.auth.IsSupported(c.auth.GetVersion()) {
			return auth.ErrNotSupportedAuthVersion
		}
		if agentInfo.SupportedAuth == nil {
			// Check agent version and auth version compatibility, if not compatible, return error.
			// 4.2.2 only support v1, 4.2.3 only support v2
			if !(c.auth.GetVersion() == auth.AUTH_V1 && auth.VERSION_4_2_2.Equals(agentInfo.Version) ||
				c.auth.GetVersion() == auth.AUTH_V2 && auth.VERSION_4_2_3.BeforeOrEquals(agentInfo.Version)) {
				return errors.New("unsupported auth version of obshell")
			}
		} else {
			for _, v := range agentInfo.SupportedAuth {
				if v == c.auth.GetVersion() {
					return nil
				}
			}
			return errors.New("unsupported auth version of obshell")
		}
		return nil
	}

	var supportedAuth []string
	if agentInfo.SupportedAuth != nil {
		supportedAuth = agentInfo.SupportedAuth
	} else if auth.VERSION_4_2_2.Equals(agentInfo.Version) {
		supportedAuth = append(supportedAuth, auth.AUTH_V1)
	} else if auth.VERSION_4_2_3.BeforeOrEquals(agentInfo.Version) {
		supportedAuth = append(supportedAuth, auth.AUTH_V2)
	} else {
		return fmt.Errorf("unsupported obshell version: %s", agentInfo.Version) // Unexpected error
	}

	if !c.auth.AutoSelectVersion(supportedAuth...) {
		return fmt.Errorf("there is no supprt auth version for target obshell(version: %s)", agentInfo.Version)
	}
	return nil
}

func (c *Client) reconfirmAuthVersion() error {
	c.auth.Reset()
	return c.confirmAuthVersion()
}

func (c *Client) tryCandidateAuth(request request.Request, response responselib.Response) bool {
	if c.candidateAuth == nil {
		return false
	}

	agentInfo, err := util.GetInfo(c.GetServer())
	if err != nil {
		return false
	}

	// Check if the agent version is less than or equal to 4.2.4.
	if auth.VERSION_4_2_4.BeforeOrEquals(agentInfo.Version) {
		// For this version, when an UnauthorizedError is returned, it may indicate issues other than just a unauthorized error.
		// Therefore, we need to reconfirm the authentication version and attempt the request again instead of immediately using the candidate.
		if err = c.reconfirmAuthVersion(); err != nil {
			return false
		}
		if err = c.realExecute(request, response); err == nil {
			return true
		}
	}

	c.AdoptCandidateAuth()
	err = c.realExecute(request, response)
	if err == nil {
		return true
	} else {
		return false
	}
}

func (c *Client) Execute(request request.Request, response responselib.Response) (err error) {
	if c.auth.GetVersion() == "" {
		if err = c.confirmAuthVersion(); err != nil {
			return err
		}
	}

	err = c.realExecute(request, response)
	if err != nil {
		apiError, ok := err.(*responselib.ApiError)
		if !ok {
			// Network error
			return err
		}

		if !c.auth.IsAutoSelectVersion() {
			// Auth version is not auto select, can't reconfirm auth version
			return err
		}

		if apiError.IsError(responselib.DecryptError) {
			c.auth.ResetMethod()
		} else {
			if apiError.IsError(responselib.UnauthorizedError) {
				if c.tryCandidateAuth(request, response) {
					return nil
				}
				// If the current auth version greater than v2, or not auto select version, return error. Because UnauthorizedError means the certificate is invalid when the auth version greater than v2
				if c.auth.GetVersion() > auth.AUTH_V2 {
					return err
				}
			} else if !apiError.IsError(responselib.IncompatibleError) {
				return err
			}

			// Maybe agent upgrade, reconfirm auth version
			if err = c.reconfirmAuthVersion(); err != nil {
				return err
			}
		}
		// Re-execute request
		return c.realExecute(request, response)
	}
	return
}

func (c *Client) realExecute(req request.Request, response responselib.Response) (err error) {
	if req == nil || reflect.ValueOf(req).IsNil() {
		return errors.New("request is nil")
	}

	requestContext := request.NewContext()
	if req.Authentication() {
		if err = c.auth.Auth(req, requestContext); err != nil {
			return err
		}
	}

	r := req.BuildHttpRequest(requestContext).
		SetError(response).
		SetResult(response)

	var resp *resty.Response
	targetUrl, err := req.BuildUrl()
	if err != nil {
		return errors.Wrap(err, "build url failed")
	}
	switch req.GetMethod() {
	case "GET":
		resp, err = r.Get(targetUrl)
	case "PUT":
		resp, err = r.Put(targetUrl)
	case "POST":
		resp, err = r.Post(targetUrl)
	case "PATCH":
		resp, err = r.Patch(targetUrl)
	case "DELETE":
		resp, err = r.Delete(targetUrl)
	default:
		return fmt.Errorf("%s method not support", req.GetMethod())
	}
	if err != nil {
		return errors.Wrap(err, "request failed")
	}
	if resp.IsError() {
		if response != nil {
			return response.GetError()
		}
		return errors.New("http response error")
	}
	return nil
}
