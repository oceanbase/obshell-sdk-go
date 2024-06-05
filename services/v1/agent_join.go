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

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type JoinRequest struct {
	*request.BaseRequest
}

type joinApiParam struct {
	AgentInfo model.AgentInfo `json:"AgentInfo" binding:"required"`
	ZoneName  string          `json:"zoneName" binding:"required"`
}

type joinResponse struct {
	*response.TaskResponse
}

func (c *Client) createJoinResponse() *joinResponse {
	return &joinResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewJoinRequest return a JoinRequest, which can be used as the argument for the JoinWithRequest/JoinSyncWithRequest.
// ip: the ip of the agent to be joined.
// port: the port of the agent to be joined.
// zone: the zone name of the observer.
func (c *Client) NewJoinRequest(ip string, port int, zone string) *JoinRequest {
	req := &JoinRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.SetBody(&joinApiParam{
		AgentInfo: model.AgentInfo{
			Ip:   c.GetHost(),
			Port: c.GetPort(),
		},
		ZoneName: zone,
	})
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/agent/join", ip, port, "POST")
	return req
}

// Join returns a DagDetailDTO and an error, when the join task is completed successfully, the error will be nil.
// ip: the ip of the agent to be joined.
// port: the port of the agent to be joined.
// zone: the zone name of the observer.
func (c *Client) Join(ip string, port int, zone string) (*model.DagDetailDTO, error) {
	req := c.NewJoinRequest(ip, port, zone)
	return c.JoinSyncWithRequest(req)
}

// JoinWithRequest returns a DagDetailDTO and an error, when the join task is requested successfully, the error will be nil.
// the parameter is a JoinRequest, which can be created by NewJoinRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) JoinWithRequest(req *JoinRequest) (dag *model.DagDetailDTO, err error) {
	auth := c.GetAuth()
	defer auth.ResetMethod()
	auth.ResetMethod()
	response := c.createJoinResponse()
	targetClient, err := NewClientWithServer(req.GetHost(), req.GetPort())
	targetClient.SetAuth(auth)
	if err != nil {
		return nil, errors.Wrap(err, "create target client failed")
	}
	err = targetClient.Execute(req, response)
	dag = response.DagDetailDTO
	return
}

// JoinSyncWithRequest returns a DagDetailDTO and an error, when the join task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a JoinRequest, which can be created by NewJoinRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) JoinSyncWithRequest(req *JoinRequest) (*model.DagDetailDTO, error) {
	dag, err := c.JoinWithRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
