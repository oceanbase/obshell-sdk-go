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

type ScaleOutRequest struct {
	*request.BaseRequest
}

type ScaleOutApiParam struct {
	AgentInfo model.AgentInfo   `json:"AgentInfo" binding:"required"`
	Zone      string            `json:"zone" binding:"required"`
	ObConfigs map[string]string `json:"obConfigs"`
}

type ScaleOutResponse struct {
	*response.TaskResponse
}

func (c *Client) createScaleOutResponse() *ScaleOutResponse {
	return &ScaleOutResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewScaleOutRequest return a ScaleOutRequest, which can be used as the argument for the ScaleOutWithRequest/ScaleOutSyncWithRequest.
// ip: the ip of the agent to be scale_out.
// port: the port of the agent to be scale_out.
// zone: the zone name of the observer.
// obConfigs: the observer configs.
func (c *Client) NewScaleOutRequest(ip string, port int, zone string, obConfigs map[string]string) *ScaleOutRequest {
	req := &ScaleOutRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.SetBody(&ScaleOutApiParam{
		AgentInfo: model.AgentInfo{
			Ip:   ip,
			Port: port,
		},
		Zone:      zone,
		ObConfigs: obConfigs,
	})
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/scale_out", c.GetHost(), c.GetPort(), "POST")
	return req
}

// ScaleOut returns a DagDetailDTO and an error, when the scale_out task is completed successfully, the error will be nil.
// ip: the ip of the agent to be scaled out.
// port: the port of the agent to be scaled out.
// zone: the zone name of the observer.
// obConfigs: the observer configs.
func (c *Client) ScaleOut(ip string, port int, zone string, obConfigs map[string]string) (*model.DagDetailDTO, error) {
	req := c.NewScaleOutRequest(ip, port, zone, obConfigs)
	return c.ScaleOutSyncWithRequest(req)
}

// ScaleOutWithRequest returns a DagDetailDTO and an error, when the scale_out task is requested successfully, the error will be nil.
// the parameter is a ScaleOutRequest, which can be created by NewScaleOutRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleOutWithRequest(request *ScaleOutRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createScaleOutResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ScaleOutSyncWithRequest returns a DagDetailDTO and an error, when the scale_out task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ScaleOutRequest, which can be created by NewScaleOutRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ScaleOutSyncWithRequest(request *ScaleOutRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ScaleOutWithRequest(request)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}
