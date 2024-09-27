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

type ScaleOutReplicasRequest struct {
	*request.BaseRequest
	param ScaleOutReplicasParam
}

type ScaleOutReplicasParam struct {
	ZoneList []ZoneParam `json:"zone_list" binding:"required"` // Tenant zone list with unit config.
}

type ScaleOutReplicasResponse struct {
	*response.TaskResponse
}

func (c *Client) createScaleOutReplicasResponse() *ScaleOutReplicasResponse {
	return &ScaleOutReplicasResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewScaleOutReplicasRequest return a ScaleOutReplicasRequest, which can be used as the argument for the ScaleOutReplicasWithRequest/ScaleOutReplicasSyncWithRequest.
// tenantName: the name of the tenant.
// zoneList: the zone list with replicas properties to be scaled out.
func (c *Client) NewScaleOutReplicasRequest(tenantName string, zoneList []ZoneParam) *ScaleOutReplicasRequest {
	req := &ScaleOutReplicasRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: ScaleOutReplicasParam{
			ZoneList: zoneList,
		},
	}
	req.SetBody(&req.param)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/replicas", tenantName), c.GetHost(), c.GetPort(), "POST")
	return req
}

// ScaleOutReplicas returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// tenantName: the name of the tenant.
// zoneList: the zone list with replicas properties to be scaled out.
func (c *Client) ScaleOutReplicas(tenantName string, zoneList []ZoneParam) (*model.DagDetailDTO, error) {
	request := c.NewScaleOutReplicasRequest(tenantName, zoneList)
	return c.ScaleOutReplicasSyncWithRequest(request)
}

// ScaleOutReplicasWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a ScaleOutReplicasRequest, which can be created by NewScaleOutReplicasRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleOutReplicasWithRequest(request *ScaleOutReplicasRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createScaleOutReplicasResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ScaleOutReplicasSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ScaleOutReplicasRequest, which can be created by NewScaleOutReplicasRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleOutReplicasSyncWithRequest(request *ScaleOutReplicasRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ScaleOutReplicasWithRequest(request)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}
