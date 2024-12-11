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
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type ScaleInRequest struct {
	*request.BaseRequest
	param scaleInParam
}

type scaleInParam struct {
	AgentInfo model.AgentInfo `json:"agent_info" binding:"required"`
	ForceKill bool            `json:"force_kill"`
}

type ScaleInResponse struct {
	*response.TaskResponse
}

func (c *Client) createScaleInResponse() *ScaleInResponse {
	return &ScaleInResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewScaleInRequest returns a ScaleInRequest, which can be used as the argument for the ScaleInWithRequest/ScaleInSyncWithRequest.
// ip: the agent ip to be scaled in.
// port: the agent port to be scaled in.
func (c *Client) NewScaleInRequest(ip string, port int) *ScaleInRequest {
	req := &ScaleInRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.param.AgentInfo = model.AgentInfo{
		Ip:   ip,
		Port: port,
	}
	req.SetBody(&req.param)

	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/scale_in", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetForceKill sets ForceKill to true.
func (r *ScaleInRequest) SetForceKill() *ScaleInRequest {
	r.param.ForceKill = true
	r.SetBody(r.param)
	return r
}

// ScaleIn returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// ip: the agent ip to be scaled in.
// port: the agent port to be scaled in.
func (c *Client) ScaleIn(ip string, port int) (*model.DagDetailDTO, error) {
	request := c.NewScaleInRequest(ip, port)
	return c.ScaleInSyncWithRequest(request)
}

// ScaleInWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleInWithRequest(request *ScaleInRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createScaleInResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ScaleInSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleInSyncWithRequest(request *ScaleInRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ScaleInWithRequest(request)
	if err != nil {
		return nil, err
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
