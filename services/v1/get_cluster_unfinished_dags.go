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

type GetClusterUnfinishedDagsRequest struct {
	*request.BaseRequest
}

// NewGetClusterUnfinishedDagsRequest return a GetGetClusterUnfinishedDagsRequest, which can be used as the argument for the GetClusterUnfinishedDagsWithRequest.
// You can set whether show detail by calling SetShowDetail.
func (c *Client) NewGetClusterUnfinishedDagsRequest() *GetClusterUnfinishedDagsRequest {
	req := &GetClusterUnfinishedDagsRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/dag/ob/unfinish", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

// SetShowDetail set whether show detail.

func (r *GetClusterUnfinishedDagsRequest) SetShowDetail(showDetail bool) *GetClusterUnfinishedDagsRequest {
	r.SetQueryParam("show_details", fmt.Sprintf("%t", showDetail))
	return r
}

// GetClusterUnfinishedDags returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't want to show detail of the dag, you can need to create a GetClusterUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetClusterUnfinishedDags() (dag []*model.DagDetailDTO, err error) {
	req := c.NewGetClusterUnfinishedDagsRequest()
	return c.GetClusterUnfinishedDagsWithRequest(req)
}

// GetClusterUnfinishedDagsWithRequest returns a []*DagDetailDTO and an error.
// If the error is non-nil, the []*DagDetailDTO will be empty.
// If you don't want to show detail of the dag, you can need to create a GetClusterUnfinishedDagsRequest and call SetShowDetail(false).
func (c *Client) GetClusterUnfinishedDagsWithRequest(req *GetClusterUnfinishedDagsRequest) (dag []*model.DagDetailDTO, err error) {
	resp := response.NewMutilTaskReponse()
	err = c.Execute(req, resp)
	return resp.Contents, err
}
