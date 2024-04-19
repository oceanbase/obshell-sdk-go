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
	"github.com/pkg/errors"

	"github.com/obshell-sdk-go/model"
	"github.com/obshell-sdk-go/sdk/request"
	"github.com/obshell-sdk-go/sdk/response"
)

type RemoveRequest struct {
	*request.BaseRequest
}

type removeResponse struct {
	*response.TaskResponse
}

func (c *Client) createRemoveResponse() *removeResponse {
	return &removeResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewRemoveRequest returns a RemoveRequest, which can be used as the argument for the RemoveWithRequest/RemoveSyncWithRequest
// ip: the ip of the agent to be removed.
// port: the port of the agent to be removed.
func (c *Client) NewRemoveRequest(ip string, port int) *RemoveRequest {
	req := &RemoveRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.SetBody(&model.AgentInfo{
		Ip:   ip,
		Port: port,
	},
	)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/agent/remove", c.GetHost(), c.GetPort(), "POST")
	return req
}

// Remove returns a DagDetailDTO and an error, when the remove task is completed successfully, the error will be nil.
// ip: the ip of the agent to be removed.
// port: the port of the agent to be removed.
func (c *Client) Remove(ip string, port int) (*model.DagDetailDTO, error) {
	req := c.NewRemoveRequest(ip, port)
	return c.RemoveSyncWithRequest(req)
}

// RemoveWithRequest returns a DagDetailDTO and an error, when the remove task is requested successfully, the error will be nil.
// the parameter is a RemoveRequest, which can be created by NewRemoveRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) RemoveWithRequest(req *RemoveRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createRemoveResponse()
	err = c.Execute(req, response)
	dag = response.DagDetailDTO
	return
}

// RemoveSyncWithRequest returns a DagDetailDTO and an error, when the remove task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a RemoveRequest, which can be created by NewRemoveRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) RemoveSyncWithRequest(req *RemoveRequest) (*model.DagDetailDTO, error) {
	dag, err := c.RemoveWithRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
