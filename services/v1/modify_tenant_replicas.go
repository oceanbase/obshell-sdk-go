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

type ModifyTenantReplicasRequest struct {
	*request.BaseRequest
}

type ModifyTenantReplicasParam struct {
	ZoneList []modifyReplicaZoneParam `json:"zone_list" binding:"required"` // Tenant zone list with unit config.
}

type modifyReplicaZoneParam struct {
	Name           string  `json:"zone_name" binding:"required"`
	ReplicaType    *string `json:"replica_type"`     // Replica type, "FULL"(default) or "READONLY". optional
	UnitConfigName *string `json:"unit_config_name"` // optional
	UnitNum        *int    `json:"unit_num"`         // optional
}

type ModifyTenantReplicasResponse struct {
	*response.TaskResponse
}

func (c *Client) createModifyTenantReplicasResponse() *ModifyTenantReplicasResponse {
	return &ModifyTenantReplicasResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewModifyTenantReplicasRequest return a ModifyTenantReplicasRequest, which can be used as the argument for the ModifyTenantReplicasWithRequest/ModifyTenantReplicasSyncWithRequest.
// tenantName: the name of the tenant.
// param: the zone list with replica properties to be modified.
func (c *Client) NewModifyTenantReplicasRequest(tenantName string, param []ZoneParam) *ModifyTenantReplicasRequest {
	req := &ModifyTenantReplicasRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}

	var body []modifyReplicaZoneParam
	for i := range param {
		var t modifyReplicaZoneParam
		t.Name = param[i].Name
		if param[i].ReplicaType == "" {
			t.ReplicaType = nil
		} else {
			t.ReplicaType = &param[i].ReplicaType
		}
		if param[i].UnitConfigName == "" {
			t.UnitConfigName = nil
		} else {
			t.UnitConfigName = &param[i].UnitConfigName
		}
		if param[i].UnitNum == 0 {
			t.UnitNum = nil
		} else {
			t.UnitNum = &param[i].UnitNum
		}
		body = append(body, t)
	}
	req.SetBody(map[string]interface{}{
		"zone_list": body,
	})
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/replicas", tenantName), c.GetHost(), c.GetPort(), "PATCH")
	return req
}

// ModifyTenantReplicas returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// tenantName: the name of the tenant.
// param: the zone list with replica properties to be modified.
func (c *Client) ModifyTenantReplicas(tenantName string, param []ZoneParam) (*model.DagDetailDTO, error) {
	request := c.NewModifyTenantReplicasRequest(tenantName, param)
	return c.ModifyTenantReplicasSyncWithRequest(request)
}

// ModifyTenantReplicasWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a ModifyTenantReplicasRequest, which can be created by NewModifyTenantReplicasRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ModifyTenantReplicasWithRequest(request *ModifyTenantReplicasRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createModifyTenantReplicasResponse()
	err = c.Execute(request, response)
	dag = response.DagDetailDTO
	return
}

// ModifyTenantReplicasSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ModifyTenantReplicasRequest, which can be created by NewModifyTenantReplicasRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ModifyTenantReplicasSyncWithRequest(request *ModifyTenantReplicasRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ModifyTenantReplicasWithRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
