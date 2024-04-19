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

type UpgradeAgentCheckRequest struct {
	*request.BaseRequest
	param *UpgradeAgentCheckApiParam
}

type UpgradeAgentCheckApiParam struct {
	Version    string `json:"version"`
	Release    string `json:"release"`
	UpgradeDir string `json:"upgradeDir"`
}

type upgradeAgentCheckResponse struct {
	*response.TaskResponse
}

func (c *Client) createUpgradeAgentCheckResponse() *upgradeAgentCheckResponse {
	return &upgradeAgentCheckResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewUpgradeAgentCheckRequest return a UpgradeAgentCheckRequest, which can be used as the argument for the UpgradeAgentCheckWithRequest/UpgradeAgentCheckSyncWithRequest.
// version: the version of the agent to be upgraded to.
// release: the release of the agent to be upgraded to.
// You can set the upgradeDir by calling SetUpgradeDir.
func (c *Client) NewUpgradeAgentCheckRequest(version, release string) *UpgradeAgentCheckRequest {
	req := &UpgradeAgentCheckRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &UpgradeAgentCheckApiParam{
			Version: version,
			Release: release,
		},
	}
	req.SetBody(req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/agent/upgrade/check", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetUpgradeDir set the upgradeDir of the UpgradeAgentCheckRequest.
func (r *UpgradeAgentCheckRequest) SetUpgradeDir(upgradeDir string) *UpgradeAgentCheckRequest {
	r.param.UpgradeDir = upgradeDir
	r.SetBody(r.param)
	return r
}

// UpgradeAgentCheck returns a DagDetailDTO and an error, when the upgrade check task is completed successfully, the error will be nil.
// version: the version of the agent to be upgraded to.
// release: the release of the agent to be upgraded to.
// if you want to set the upgradeDir, you need to use NewUpgradeAgentCheckRequest and call SetUpgradeDir.
func (c *Client) UpgradeAgentCheck(version, release string) (*model.DagDetailDTO, error) {
	req := c.NewUpgradeAgentCheckRequest(version, release)
	return c.UpgradeAgentCheckSyncWithRequest(req)
}

// UpgradeAgentCheckWithRequest returns a DagDetailDTO and an error, when the upgrade check task is requested successfully, the error will be nil.
// the parameter is a UpgradeAgentCheckRequest, which can be created by NewUpgradeAgentCheckRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeAgentCheckWithRequest(req *UpgradeAgentCheckRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createUpgradeAgentCheckResponse()
	err = c.Execute(req, response)
	dag = response.DagDetailDTO
	return
}

// UpgradeAgentCheckSyncWithRequest returns a DagDetailDTO and an error, when the upgrade check task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a UpgradeAgentCheckRequest, which can be created by NewUpgradeAgentCheckRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeAgentCheckSyncWithRequest(req *UpgradeAgentCheckRequest) (*model.DagDetailDTO, error) {
	dag, err := c.UpgradeAgentCheckWithRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
