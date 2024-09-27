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

type StopRequest struct {
	*request.BaseRequest
	param stopApiParam
}

type stopApiParam struct {
	Scope             model.Scope       `json:"scope" binding:"required"`
	Force             bool              `json:"force"`
	Terminate         bool              `json:"terminate"`
	ForcePassDagParam ForcePassDagParam `json:"forcePassDag"`
}

// SetForcePassDag set the force pass dag.
func (r *StopRequest) SetForcePassDag(id []string) *StopRequest {
	r.param.ForcePassDagParam = ForcePassDagParam{
		ID: id,
	}
	r.SetBody(r.param)
	return r
}

// SetForce set the force.
// when force is true, the target observer will be forcely killed.
func (r *StopRequest) SetForce() *StopRequest {
	r.param.Force = true
	r.param.Terminate = false
	r.SetBody(r.param)
	return r
}

func (r *StopRequest) SetTerminate() *StopRequest {
	r.param.Terminate = true
	r.param.Force = false
	r.SetBody(r.param)
	return r
}

type StopResponse struct {
	*response.TaskResponse
}

func (c *Client) createStopResponse() *StopResponse {
	return &StopResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewStopRequest return a StopRequest, which can be used as the argument for the StopWithRequest/StopSyncWithRequest.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL, when level is SCOPE_GLOBAL, you need to set force by calling SetForce.
// targets is the target to be stopped, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
// You can set the force pass dag by calling SetForcePassDag.
// You can set the force by calling SetForce.
func (c *Client) NewStopRequest(level string, targets ...string) *StopRequest {
	req := &StopRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: stopApiParam{
			Scope: model.Scope{
				Type:   level,
				Target: targets,
			},
		},
	}
	req.SetBody(&req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/stop", c.GetHost(), c.GetPort(), "POST")
	return req
}

// Stop returns a DagDetailDTO and an error, when the stop task is completed successfully, the error will be nil.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL.
// targets is the target to be stopped, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
// If you want to set the force pass dag, you need to use NewStopRequest and call SetForcePassDag.
// If you want to kill observer forcely, you need to use NewStopRequest and call SetForce.
func (c *Client) Stop(level string, targets ...string) (*model.DagDetailDTO, error) {
	request := c.NewStopRequest(level, targets...)
	return c.StopSyncWithRequest(request)
}

// StopWithRequest returns a DagDetailDTO and an error, when the stop task is requested successfully, the error will be nil.
// the parameter is a StopRequest, which can be created by CreateStopRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) StopWithRequest(request *StopRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createStopResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// StopSyncWithRequest returns a DagDetailDTO and an error, when the stop task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a StopRequest, which can be created by NewStopRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) StopSyncWithRequest(request *StopRequest) (*model.DagDetailDTO, error) {
	dag, err := c.StopWithRequest(request)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}
