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

type GetInfoRequest struct {
	*request.BaseRequest
}

// NewGetInfoRequest return a GetInfoRequest, which can be used as the argument for the GetInfoWithRequest.
func (c *Client) NewGetInfoRequest() *GetInfoRequest {
	req := &GetInfoRequest{
		BaseRequest: &request.BaseRequest{},
	}
	req.InitApiInfo("/api/v1/info", c.GetHost(), c.GetPort(), "GET")
	return req
}

type getInfoResponse struct {
	*response.OcsAgentResponse
	*model.AgentRunStatus
}

func (c *Client) createGetInfoResponse() *getInfoResponse {
	resp := &getInfoResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		AgentRunStatus:   &model.AgentRunStatus{},
	}
	resp.Data = resp.AgentRunStatus
	return resp
}

// GetInfo returns a AgentRunStatus and an error.
// If the error is non-nil, the AgentRunStatus will be nil.
func (c *Client) GetInfo() (ObInfoResp *model.AgentRunStatus, err error) {
	req := c.NewGetInfoRequest()
	return c.GetInfoWithRequest(req)
}

// GetInfoWithRequest returns a AgentRunStatus and an error.
// If the error is non-nil, the AgentRunStatus will be nil.
func (c *Client) GetInfoWithRequest(req *GetInfoRequest) (ObInfoResp *model.AgentRunStatus, err error) {
	response := c.createGetInfoResponse()
	err = c.Execute(req, response)
	return response.AgentRunStatus, err
}
