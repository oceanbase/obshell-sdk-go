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

type GetAgentUnfinishedDagsRequest struct {
	*request.BaseRequest
}

// NewGetAgentUnfinishedDagsRequest return a GetAgentUnfinishedDagsRequest, which can be used as the argument for the GetAgentUnfinishedDagsWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetAgentUnfinishedDagsRequest() *GetAgentUnfinishedDagsRequest {
	req := &GetAgentUnfinishedDagsRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/agent/unfinish", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetAgentUnfinishedDagsRequest) SetShowDetail(showDetail bool) *GetAgentUnfinishedDagsRequest {
	r.SetBody(map[string]bool{"showDetail": showDetail})
	return r
}

type getAgentUnfinishedDagsResponse struct {
	*response.OcsAgentResponse
	Dags []*model.DagDetailDTO `json:"contents"`
}

func (c *Client) creatGetAgentUnfinishedDagsResponse() *getAgentUnfinishedDagsResponse {
	return &getAgentUnfinishedDagsResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		Dags:             make([]*model.DagDetailDTO, 0),
	}
}

// GetAgentUnfinishedDags returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't want to show detail of the dag, you can need to create a GetAgentUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetAgentUnfinishedDags() (dags []*model.DagDetailDTO, err error) {
	req := c.NewGetAgentUnfinishedDagsRequest()
	return c.GetAgentUnfinishedDagsWithRequest(req)
}

// GetAgentUnfinishedDagsWithRequest returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't to show detail of the dag, you can need to create a GetAgentUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetAgentUnfinishedDagsWithRequest(req *GetAgentUnfinishedDagsRequest) (dags []*model.DagDetailDTO, err error) {
	response := c.creatGetAgentUnfinishedDagsResponse()
	err = c.Execute(req, response)
	return response.Dags, err
}
