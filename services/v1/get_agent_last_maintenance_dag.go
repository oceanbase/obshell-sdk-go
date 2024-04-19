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
	"github.com/obshell-sdk-go/model"
	"github.com/obshell-sdk-go/sdk/request"
	"github.com/obshell-sdk-go/sdk/response"
)

type GetAgentLastMaintenanceDagRequest struct {
	*request.BaseRequest
}

// NewGetAgentLastMaintenanceDagRequest return a GetAgentLastMaintenanceDagRequest, which can be used as the argument for the GetAgentLastMaintenanceDagWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetAgentLastMaintenanceDagRequest() *GetAgentLastMaintenanceDagRequest {
	req := &GetAgentLastMaintenanceDagRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/maintain/agent", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetAgentLastMaintenanceDagRequest) SetShowDetail(showDetail bool) *GetAgentLastMaintenanceDagRequest {
	r.SetBody(map[string]bool{"showDetail": showDetail})
	return r
}

type getAgentLastMaintenanceDagResponse struct {
	*response.TaskResponse
}

func (c *Client) createGetAgentLastMaintenanceDagResponse() *getAgentLastMaintenanceDagResponse {
	return &getAgentLastMaintenanceDagResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// GetAgentLastMaintenanceDag returns a DagDetailDTO and an error.
// If the error is non-nil, the DagDetailDTO will be nil.
// If you don't want to show detail of the dag, you can need to create a GetAgentUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetAgentLastMaintenanceDag() (dag *model.DagDetailDTO, err error) {
	req := c.NewGetAgentLastMaintenanceDagRequest()
	return c.GetAgentLastMaintenanceDagWithRequest(req)
}

// GetAgentLastMaintenanceDagWithRequest returns a DagDetailDTO and an error.
// If the error is non-nil, the DagDetailDTO will be nil.
// If you don't want to show detail of the dag, you can need to create a GetAgentUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetAgentLastMaintenanceDagWithRequest(req *GetAgentLastMaintenanceDagRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createGetAgentLastMaintenanceDagResponse()
	err = c.Execute(req, response)
	return response.DagDetailDTO, err
}
