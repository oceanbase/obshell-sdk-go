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

package request

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Request interface {
	GetMethod() string
	GetBody() interface{}
	SetBody(body interface{})
	SetHeader(k, v string)
	BuildUrl() string
	GetUri() string
	GetServer() string
	Authentication() bool
	IsAsync() bool
	BuildHttpRequest(context *Context) *resty.Request
}

type BaseRequest struct {
	host           string
	port           int
	uri            string
	Protocol       string
	method         string
	authentication bool
	Version        string
	header         map[string]string
	body           interface{}
	files          map[string]string
	isAsync        bool
}

func NewBaseRequest() *BaseRequest {
	return &BaseRequest{}
}

func NewAsyncBaseRequest() *BaseRequest {
	return &BaseRequest{
		isAsync: true,
	}
}

func (r *BaseRequest) IsAsync() bool {
	return r.isAsync
}

func (r *BaseRequest) GetBody() interface{} {
	return r.body
}

func (r *BaseRequest) SetBody(body interface{}) {
	r.body = body
}

func (r *BaseRequest) SetHeader(k, v string) {
	r.header[k] = v
}

func (r *BaseRequest) GetServer() string {
	return fmt.Sprintf("%s:%d", r.host, r.port)
}

func (r *BaseRequest) GetHost() string {
	return r.host
}

func (r *BaseRequest) GetPort() int {
	return r.port
}

func (r *BaseRequest) BuildUrl() string {
	return fmt.Sprintf("%s://%s%s", r.Protocol, r.GetServer(), r.uri)
}

func (r *BaseRequest) GetUri() string {
	return r.uri
}

func (r *BaseRequest) GetMethod() string {
	return r.method
}

func (r *BaseRequest) SetAuthentication() {
	r.authentication = true
}

func (r *BaseRequest) Authentication() bool {
	return r.authentication
}

func (r *BaseRequest) InitApiInfo(uri, host string, port int, method string) {
	r.uri = uri
	r.host = host
	r.port = port
	r.header = map[string]string{
		"Content-Type": "application/json",
	}
	r.files = map[string]string{}
	r.Protocol = "http"
	r.method = method
}

func (r *BaseRequest) BuildHttpRequest(context *Context) *resty.Request {
	// The default format of the request is JSON.
	req := resty.New().R()

	// Set headers which are not in context, set by service.
	for k, v := range r.header {
		req.SetHeader(k, v)
	}

	for k, v := range context.headers {
		req.SetHeader(k, v)
	}

	if context.body != nil {
		req.SetBody(context.body)
	} else {
		req.SetBody(r.body)
	}

	return req
}
