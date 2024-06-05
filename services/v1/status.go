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

type GetStatusRequest struct {
	*request.BaseRequest
}

// NewGetStatusRequest return a GetStatusRequest, which can be used as the argument for the GetStatusWithRequest.
func (c *Client) NewGetStatusRequest() *GetStatusRequest {
	req := &GetStatusRequest{
		BaseRequest: &request.BaseRequest{},
	}
	req.InitApiInfo("/api/v1/status", c.GetHost(), c.GetPort(), "GET")
	return req
}

type GetStatusResponse struct {
	*response.OcsAgentResponse
	*model.AgentStatus
}

func (c *Client) createGetStatusResponse() *GetStatusResponse {
	return &GetStatusResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// GetStatus returns a AgentStatus and an error.
// If the error is non-nil, the AgentStatus will be nil.
func (c *Client) GetStatus() (ObStatusResp *model.AgentStatus, err error) {
	req := c.NewGetStatusRequest()
	return c.GetStatusWithRequest(req)
}

// GetStatus returns a AgentStatus and an error.
// If the error is non-nil, the AgentStatus will be nil.
func (c *Client) GetStatusWithRequest(req *GetStatusRequest) (ObStatusResp *model.AgentStatus, err error) {
	response := c.createGetStatusResponse()
	err = c.Execute(req, response)
	return response.AgentStatus, err
}
