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

type UpgradeAgentRequest struct {
	*request.BaseRequest
	param *UpgradeAgentCheckApiParam
}

type upgradeAgentResponse struct {
	*response.TaskResponse
}

func (c *Client) createUpgradeAgentResponse() *upgradeAgentResponse {
	return &upgradeAgentResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewUpgradeAgentRequest return a UpgradeAgentRequest, which can be used as the argument for the UpgradeAgentWithRequest/UpgradeAgentSyncWithRequest.
// version: the version of the agent to be upgraded to.
// release: the release of the agent to be upgraded to.
// You can set the upgradeDir by calling SetUpgradeDir.
func (c *Client) NewUpgradeAgentRequest(version, release string) *UpgradeAgentRequest {
	req := &UpgradeAgentRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &UpgradeAgentCheckApiParam{
			Version: version,
			Release: release,
		},
	}
	req.SetBody(req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/agent/upgrade", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetUpgradeDir set the upgradeDir of the UpgradeAgentRequest.
func (r *UpgradeAgentRequest) SetUpgradeDir(upgradeDir string) *UpgradeAgentRequest {
	r.param.UpgradeDir = upgradeDir
	r.SetBody(r.param)
	return r
}

// UpgradeAgent returns a DagDetailDTO and an error, when the upgrade agent task is completed successfully, the error will be nil.
// version: the version of the agent to be upgraded to.
// release: the release of the agent to be upgraded to.
// if you want to set the upgradeDir, you need to use NewUpgradeAgentRequest and call SetUpgradeDir.
func (c *Client) UpgradeAgent(version, release string) (*model.DagDetailDTO, error) {
	req := c.NewUpgradeAgentRequest(version, release)
	return c.UpgradeAgentSyncWithRequest(req)
}

// UpgradeAgentWithRequest returns a DagDetailDTO and an error, when the upgrade agent task is requested successfully, the error will be nil.
// the parameter is a UpgradeAgentRequest, which can be created by NewUpgradeAgentRequest.
// You may need to call WaitDagSucceedWithRetry instead of WaitDagSucceed to query the task status.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeAgentWithRequest(req *UpgradeAgentRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createUpgradeAgentResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// UpgradeAgentSyncWithRequest returns a DagDetailDTO and an error, when the upgrade agent task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a UpgradeAgentRequest, which can be created by NewUpgradeAgentRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeAgentSyncWithRequest(req *UpgradeAgentRequest) (*model.DagDetailDTO, error) {
	dag, err := c.UpgradeAgentWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceedWithRetry(dag.GenericID, 600)
}
