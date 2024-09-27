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

const (
	PARAM_ROLLING_UPGRADE      = "ROLLING"
	PARAM_STOP_SERVICE_UPGRADE = "STOPSERVICE"
)

type ObUpgradeParam struct {
	Version    string `json:"version" binding:"required"`
	Release    string `json:"release" binding:"required"`
	Mode       string `json:"mode" binding:"required"`
	UpgradeDir string `json:"upgradeDir" `
}

type UpgradeObRequest struct {
	*request.BaseRequest
	param *ObUpgradeParam
}

type UpgradeObResponse struct {
	*response.TaskResponse
}

func (c *Client) createUpgradeObResponse() *UpgradeObResponse {
	return &UpgradeObResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewUpgradeObRequest return a UpgradeObRequest, which can be used as the argument for the UpgradeObWithRequest/UpgradeObSyncWithRequest.
// version: the version of the ob to be upgraded to.
// release: the release of the ob to be upgraded to.
// You can set the upgradeDir by calling SetUpgradeDir.
func (c *Client) NewUpgradeObRequest(version, release, mode string) *UpgradeObRequest {
	req := &UpgradeObRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &ObUpgradeParam{
			Version: version,
			Release: release,
			Mode:    mode,
		},
	}
	req.SetBody(req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/upgrade", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetUpgradeDir set the upgradeDir of the UpgradeObRequest.

func (r *UpgradeObRequest) SetUpgradeDir(upgradeDir string) *UpgradeObRequest {
	r.param.UpgradeDir = upgradeDir
	r.SetBody(r.param)
	return r
}

// UpgradeOb returns a DagDetailDTO and an error, when the upgrade ob task is completed successfully, the error will be nil.
// version: the version of the ob to be upgraded to.
// release: the release of the ob to be upgraded to.
// if you want to set the upgradeDir, you need to use NewUpgradeObRequest and call SetUpgradeDir.
func (c *Client) UpgradeOb(version, release, mode string) (*model.DagDetailDTO, error) {
	req := c.NewUpgradeObRequest(version, release, mode)
	return c.UpgradeObSyncWithRequest(req)
}

// UpgradeObWithRequest returns a DagDetailDTO and an error, when the upgrade ob task is requested successfully, the error will be nil.
// the parameter is a UpgradeObRequest, which can be created by NewUpgradeObRequest.
// You need to call WaitDagSucceedWithRetry instead of WaitDagSucceed to query the task status.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeObWithRequest(req *UpgradeObRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createUpgradeObResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// UpgradeObSyncWithRequest returns a DagDetailDTO and an error, when the upgrade ob task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a UpgradeObRequest, which can be created by NewUpgradeObRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeObSyncWithRequest(req *UpgradeObRequest) (*model.DagDetailDTO, error) {
	dag, err := c.UpgradeObWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceedWithRetry(dag.GenericID, 600)
}
