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

type GetAllAgentLastMaintenanceDagRequest struct {
	*request.BaseRequest
}

// NewGetAllAgentLastMaintenanceDagRequest return a GetAllAgentLastMaintenanceDagRequest, which can be used as the argument for the GetAllAgentLastMaintenanceDagWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetAllAgentLastMaintenanceDagRequest() *GetAllAgentLastMaintenanceDagRequest {
	req := &GetAllAgentLastMaintenanceDagRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/maintain/agents", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.
func (r *GetAllAgentLastMaintenanceDagRequest) SetShowDetail(showDetail bool) *GetAllAgentLastMaintenanceDagRequest {
	r.SetQueryParam("show_details", fmt.Sprintf("%t", showDetail))
	return r
}

// GetAllAgentLastMaintenanceDag returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't want to show detail of the dag, you can need to create a GetAllAgentLastMaintenanceDagRequest and call SetShowDetail(false).
// Notice: versions of obshell prior to 4.2.3 do not support this method.
func (c *Client) GetAllAgentLastMaintenanceDag() (dags []*model.DagDetailDTO, err error) {
	req := c.NewGetAllAgentLastMaintenanceDagRequest()
	return c.GetAllAgentLastMaintenanceDagWithRequest(req)
}

// GetAllAgentLastMaintenanceDagWithRequest returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't want to show detail of the dag, you can need to create a GetAllAgentLastMaintenanceDagRequest and call SetShowDetail(false).
// Notice: versions of obshell prior to 4.2.3 do not support this method.
func (c *Client) GetAllAgentLastMaintenanceDagWithRequest(req *GetAllAgentLastMaintenanceDagRequest) (dags []*model.DagDetailDTO, err error) {
	resp := response.NewMutilTaskReponse()
	err = c.Execute(req, resp)
	return resp.Contents, err
}
