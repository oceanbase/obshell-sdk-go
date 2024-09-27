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
	"fmt"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type GetObLastMaintenanceDagRequest struct {
	*request.BaseRequest
}

// NewGetObLastMaintenanceDagRequestRequest return a GetObLastMaintenanceDagRequestRequest, which can be used as the argument for the GetObLastMaintenanceDagRequestWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetObLastMaintenanceDagRequest() *GetObLastMaintenanceDagRequest {
	req := &GetObLastMaintenanceDagRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/maintain/ob", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetObLastMaintenanceDagRequest) SetShowDetail(showDetail bool) *GetObLastMaintenanceDagRequest {
	r.SetQueryParam("show_details", fmt.Sprintf("%t", showDetail))
	return r
}

type getObLastMaintenanceDagResponse struct {
	*response.TaskResponse
}

func (c *Client) createGetObLastMaintenanceDagResponse() *getObLastMaintenanceDagResponse {
	return &getObLastMaintenanceDagResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// GetObLastMaintenanceDagRequest returns a DagDetailDTO and an error.
// If the error is non-nil, the DagDetailDTO will be nil.
// If you don't want to show detail of the dag, you can need to create a GetObLastMaintenanceDagRequestRequest and call SetShowDetail(false).
func (c *Client) GetObLastMaintenanceDag() (dag *model.DagDetailDTO, err error) {
	req := c.NewGetObLastMaintenanceDagRequest()
	return c.GetObLastMaintenanceDagWithRequest(req)
}

// GetObLastMaintenanceDagWithRequest returns a DagDetailDTO and an error.
// If the error is non-nil, the DagDetailDTO will be nil.
// If you don't want to show detail of the dag, you can need to create a GetObLastMaintenanceDagRequestRequest and call SetShowDetail(false).
func (c *Client) GetObLastMaintenanceDagWithRequest(req *GetObLastMaintenanceDagRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createGetObLastMaintenanceDagResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, err
}
