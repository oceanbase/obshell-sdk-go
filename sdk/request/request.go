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

import "fmt"

type Request interface {
	GetMethod() string
	GetHeader() map[string]string
	Body() interface{}
	SetBody(body interface{})
	SetOriginalBody(body interface{})
	OriginalBody() interface{}
	BuildUrl() string
	GetUri() string
	GetServer() string
	SetFile(string, string)
	GetFiles() map[string]string
	SetHeaderByKey(string, string)
	Authentication() bool
	IsAsync() bool
}

type BaseRequest struct {
	host           string
	port           int
	uri            string
	Protocol       string
	Method         string
	authentication bool
	Version        string
	Header         map[string]string
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

func (r *BaseRequest) Body() interface{} {
	return r.body
}

func (r *BaseRequest) SetBody(body interface{}) {
	r.body = body
}

func (r *BaseRequest) SetHeader(k, v string) {
	r.Header[k] = v
}

func (r *BaseRequest) OriginalBody() interface{} {
	return r.body
}

func (r *BaseRequest) SetOriginalBody(body interface{}) {
	r.body = body
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

func (r *BaseRequest) SetFile(name, dir string) {
	r.files[name] = dir
}

func (r *BaseRequest) GetFiles() map[string]string {
	return r.files
}

func (r *BaseRequest) GetMethod() string {
	return r.Method
}

func (r *BaseRequest) SetAuthentication() {
	r.authentication = true
}

func (r *BaseRequest) Authentication() bool {
	return r.authentication
}

func (r *BaseRequest) GetHeader() map[string]string {
	return r.Header
}

func (r *BaseRequest) SetHeaderByKey(key, value string) {
	r.Header[key] = value
}

func (r *BaseRequest) InitApiInfo(uri, host string, port int, method string) {
	r.uri = uri
	r.host = host
	r.port = port
	r.Header = map[string]string{}
	r.files = map[string]string{}
	r.Protocol = "http"
	r.Method = method
}
