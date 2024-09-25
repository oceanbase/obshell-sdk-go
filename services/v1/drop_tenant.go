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

	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type DropTenantRequest struct {
	*request.BaseRequest
	param DropTenantParam
}

type DropTenantParam struct {
	NeedRecycle bool `json:"need_recycle"` // Whether to recycle tenant(can be flashback).
}

type DropTenantResponse struct {
	*response.TaskResponse
}

func (c *Client) createDropTenantResponse() *DropTenantResponse {
	return &DropTenantResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewDropTenantRequest return a DropTenantRequest, which can be used as the argument for the DropTenantWithRequest/DropTenantSyncWithRequest.
// tenantName: the name of the tenant.
// You can use SetNeedRecycle to set whether need recycle.
func (c *Client) NewDropTenantRequest(tenantName string) *DropTenantRequest {
	req := &DropTenantRequest{
		BaseRequest: request.NewBaseRequest(),
		param: DropTenantParam{
			NeedRecycle: false,
		},
	}
	req.SetBody(&req.param)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s", tenantName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

func (r *DropTenantRequest) SetNeedRecycle(needRecycle bool) *DropTenantRequest {
	r.param.NeedRecycle = needRecycle
	r.SetBody(r.param)
	return r
}

// DropTenant returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// tenantName: the name of the tenant.
// The tenant will be dropped purgely no matter whether recyclebin is enabled.
func (c *Client) DropTenant(tenantName string) (*model.DagDetailDTO, error) {
	request := c.NewDropTenantRequest(tenantName)
	return c.DropTenantSyncWithRequest(request)
}

// DropTenantWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a DropTenantRequest, which can be created by NewDropTenantRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
// The tenant will be dropped purgely no matter whether recyclebin is enabled.
// If the tenant is not exist, the DagDetailDTO will be nil.
func (c *Client) DropTenantWithRequest(request *DropTenantRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createDropTenantResponse()
	err = c.Execute(request, response)
	dag = response.DagDetailDTO
	return
}

// DropTenantSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a DropTenantRequest, which can be created by NewDropTenantRequest.
// You can check or operater the task through the DagDetailDTO.
// The tenant will be dropped purgely no matter whether recyclebin is enabled.
// If the tenant is not exist, the DagDetailDTO will be nil.
func (c *Client) DropTenantSyncWithRequest(request *DropTenantRequest) (dag *model.DagDetailDTO, err error) {
	if dag, err = c.DropTenantWithRequest(request); err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
