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

type ScaleInReplicasRequest struct {
	*request.BaseRequest
}

type ScaleInReplicasParam struct {
	Zones []string `json:"zones" binding:"required"`
}

type ScaleInReplicasResponse struct {
	*response.TaskResponse
}

func (c *Client) createScaleInReplicasResponse() *ScaleInReplicasResponse {
	return &ScaleInReplicasResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewScaleInReplicasRequest return a ScaleInReplicasRequest, which can be used as the argument for the ScaleInReplicasWithRequest/ScaleInReplicasSyncWithRequest.
// tenantName: the name of the tenant.
// zones: the zones to be scaled in.
func (c *Client) NewScaleInReplicasRequest(tenantName string, zones []string) *ScaleInReplicasRequest {
	req := &ScaleInReplicasRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.SetBody(&ScaleInReplicasParam{
		Zones: zones,
	})
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/replicas", tenantName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

// ScaleInReplicas returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// tenantName: the name of the tenant.
// zones: the zones to be scaled in.
func (c *Client) ScaleInReplicas(tenantName string, zones []string) (*model.DagDetailDTO, error) {
	request := c.NewScaleInReplicasRequest(tenantName, zones)
	return c.ScaleInReplicasSyncWithRequest(request)
}

// ScaleInReplicasWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a ScaleInReplicasRequest, which can be created by NewScaleInReplicasRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleInReplicasWithRequest(request *ScaleInReplicasRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createScaleInReplicasResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ScaleInReplicasSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ScaleInReplicasRequest, which can be created by NewScaleInReplicasRequest.
// You can check or operater the task through the DagDetailDTO.
// If the replica does not exist, the DagDetailDTO will be nil.
func (c *Client) ScaleInReplicasSyncWithRequest(request *ScaleInReplicasRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ScaleInReplicasWithRequest(request)
	if err != nil {
		return nil, err
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
