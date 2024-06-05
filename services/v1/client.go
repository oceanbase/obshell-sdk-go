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

package v1

import (
	"github.com/oceanbase/obshell-sdk-go/sdk"
	"github.com/oceanbase/obshell-sdk-go/sdk/auth"
	"github.com/oceanbase/obshell-sdk-go/sdk/option"
)

type Client struct {
	*sdk.Client
}

// NewClient creates a new client with the given host and port.
//
// The client will use the given options to configure the client.
// You can use the sdk.WithPasswordAuth option to set the password for the client.
//
// AS: sdk.NewClient("127.0.0.1", 8080, sdk.WithPasswordAuth("password"))
func NewClient(host string, port int, options ...option.Optioner) (*Client, error) {
	c, err := sdk.NewClient(host, port, options...)
	return &Client{c}, err
}

func NewClientWithServer(host string, port int) (*Client, error) {
	c, err := sdk.NewClientWithServer(host, port)
	return &Client{c}, err
}

func NewClientWithPassword(host string, port int, password string) (*Client, error) {
	c, err := sdk.NewClientWithPassword(host, port, password)
	return &Client{c}, err
}

func (c *Client) setPasswordCandidateAuth(password string) {
	if c.GetAuth().Type() == auth.AUTH_TYPE_PASSWORD {
		candidateAuth := auth.NewPasswordAuth(password)
		c.SetCandidateAuth(candidateAuth)
	}
}
