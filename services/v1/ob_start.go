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

type StartRequest struct {
	*request.BaseRequest
	param StartApiParam
}

type StartApiParam struct {
	Scope             model.Scope       `json:"scope" binding:"required"`
	ForcePassDagParam ForcePassDagParam `json:"forcePassDag"`
}

type ForcePassDagParam struct {
	ID []string `json:"id"`
}

type StartResponse struct {
	*response.TaskResponse
}

func (c *Client) createStartResponse() *StartResponse {
	return &StartResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewStartRequest return a StartRequest, which can be used as the argument for the StartWithRequest/StartSyncWithRequest.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL.
// targets is the target to be started, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
// You can set the force pass dag by calling SetForcePassDag.
func (c *Client) NewStartRequest(level string, targets ...string) *StartRequest {
	req := &StartRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: StartApiParam{
			Scope: model.Scope{
				Type:   level,
				Target: targets,
			},
		},
	}
	req.SetBody(&req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/start", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetForcePassDag set the force pass dag of the start request.
func (r *StartRequest) SetForcePassDag(id []string) *StartRequest {
	r.param.ForcePassDagParam = ForcePassDagParam{
		ID: id,
	}
	r.SetBody(r.param)
	return r
}

// Start returns a DagDetailDTO and an error, when the start task is completed successfully, the error will be nil.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL.
// targets is the target to be started, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
// If you want to set the force pass dag, you need to use NewStartRequest and call SetForcePassDag.
func (c *Client) Start(level string, targets ...string) (*model.DagDetailDTO, error) {
	request := c.NewStartRequest(level, targets...)
	return c.StartSyncWithRequest(request)
}

// StartWithRequest returns a DagDetailDTO and an error, when the start task is requested successfully, the error will be nil.
// the parameter is a StartRequest, which can be created by NewStartRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) StartWithRequest(request *StartRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createStartResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// StartSyncWithRequest returns a DagDetailDTO and an error, when the start task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a StartRequest, which can be created by NewStartRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) StartSyncWithRequest(request *StartRequest) (*model.DagDetailDTO, error) {
	dag, err := c.StartWithRequest(request)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}
