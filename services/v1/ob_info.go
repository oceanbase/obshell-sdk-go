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
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type GetObInfoRequest struct {
	*request.BaseRequest
}

// NewGetObInfoRequest return a GetObInfoRequest, which can be used as the argument for the GetObInfoWithRequest.
func (c *Client) NewGetObInfoRequest() *GetObInfoRequest {
	req := &GetObInfoRequest{
		BaseRequest: &request.BaseRequest{},
	}
	req.InitApiInfo("/api/v1/ob/info", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getObInfoResponse struct {
	*response.OcsAgentResponse
	*model.ObInfoResp
}

func (c *Client) createGetObInfoResponse() *getObInfoResponse {
	resp := &getObInfoResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		ObInfoResp:       &model.ObInfoResp{},
	}
	resp.Data = resp.ObInfoResp
	return resp
}

// GetObInfo returns a ObInfoResp and an error.
// If the error is non-nil, the ObInfoResp will be nil.
func (c *Client) GetObInfo() (ObInfoResp *model.ObInfoResp, err error) {
	req := c.NewGetObInfoRequest()
	return c.GetObInfoWithRequest(req)
}

// GetObInfoWithRequest returns a ObInfoResp and an error.
// If the error is non-nil, the ObInfoResp will be nil.
func (c *Client) GetObInfoWithRequest(req *GetObInfoRequest) (ObInfoResp *model.ObInfoResp, err error) {
	response := c.createGetObInfoResponse()
	err = c.Execute(req, response)
	return response.ObInfoResp, err
}
