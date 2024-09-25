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
	"github.com/pkg/errors"
)

type PurgeRecyclebinTenantRequest struct {
	*request.BaseRequest
}

type purgeRecyclebinTenantResponse struct {
	*response.TaskResponse
}

func (c *Client) createPurgeRecyclebinTenantResponse() *purgeRecyclebinTenantResponse {
	return &purgeRecyclebinTenantResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewPurgeRecyclebinTenantRequest return a PurgeRecyclebinTenantRequest, which can be used as the argument for the PurgeRecyclebinTenantWithRequest/PurgeRecyclebinTenantSyncWithRequest.
// objectOrOriginalName: the name of the object(tenant) in recyclebin or the original name of the tenant.
func (c *Client) NewPurgeRecyclebinTenantRequest(objectOrOriginalName string) *PurgeRecyclebinTenantRequest {
	req := &PurgeRecyclebinTenantRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/recyclebin/tenant/%s", objectOrOriginalName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

// PurgeRecyclebinTenant purges a tenant from recyclebin.
func (c *Client) PurgeRecyclebinTenant(tenantName string) (dag *model.DagDetailDTO, err error) {
	request := c.NewPurgeRecyclebinTenantRequest(tenantName)
	return c.PurgeRecyclebinTenantSyncWithRequest(request)
}

// PurgeRecyclebinTenantWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a PurgeRecyclebinTenantRequest, which can be created by NewPurgeRecyclebinTenantRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) PurgeRecyclebinTenantWithRequest(request *PurgeRecyclebinTenantRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createPurgeRecyclebinTenantResponse()
	err = c.Execute(request, response)
	dag = response.DagDetailDTO
	return
}

// PurgeRecyclebinTenantSyncWithRequest returns a DagDetailDTO and an error.
// the DagDetailDTO is the final status of the task.
// the parameter is a PurgeRecyclebinTenantRequest, which can be created by NewPurgeRecyclebinTenantRequest.
// You can check or operater the task through the DagDetailDTO.
// If the tenant is not exist in recyclebin, the DagDetailDTO will be nil.
func (c *Client) PurgeRecyclebinTenantSyncWithRequest(request *PurgeRecyclebinTenantRequest) (dag *model.DagDetailDTO, err error) {
	if dag, err = c.PurgeRecyclebinTenantWithRequest(request); err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
