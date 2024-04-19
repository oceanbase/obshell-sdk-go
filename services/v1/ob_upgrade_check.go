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

type UpgradeObCheckRequest struct {
	*request.BaseRequest
	param *UpgradeObCheckApiParam
}

type UpgradeObCheckApiParam struct {
	Version    string `json:"version"`
	Release    string `json:"release"`
	UpgradeDir string `json:"upgradeDir"`
}

type UpgradeObCheckResponse struct {
	*response.TaskResponse
}

func (c *Client) createUpgradeObCheckResponse() *UpgradeObCheckResponse {
	return &UpgradeObCheckResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewUpgradeObCheckRequest return a UpgradeObCheckRequest, which can be used as the argument for the UpgradeObCheckWithRequest/UpgradeObCheckSyncWithRequest.
// version: the version of the ob to be upgraded to.
// release: the release of the ob to be upgraded to.
// You can set the upgradeDir by calling SetUpgradeDir.
func (c *Client) NewUpgradeObCheckRequest(version, release string) *UpgradeObCheckRequest {
	req := &UpgradeObCheckRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &UpgradeObCheckApiParam{
			Version: version,
			Release: release,
		},
	}
	req.SetBody(req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/ob/upgrade/check", c.GetHost(), c.GetPort(), "POST")
	return req
}

// SetUpgradeDir set the upgradeDir of the UpgradeObCheckRequest.
func (r *UpgradeObCheckRequest) SetUpgradeDir(upgradeDir string) *UpgradeObCheckRequest {
	r.param.UpgradeDir = upgradeDir
	r.SetBody(r.param)
	return r
}

// UpgradeObCheck returns a DagDetailDTO and an error, when the upgrade check task is completed successfully, the error will be nil.
// version: the version of the ob to be upgraded to.
// release: the release of the ob to be upgraded to.
// if you want to set the upgradeDir, you need to use NewUpgradeObCheckRequest and call SetUpgradeDir.
func (c *Client) UpgradeObCheck(version, release string) (*model.DagDetailDTO, error) {
	req := c.NewUpgradeObCheckRequest(version, release)
	return c.UpgradeObCheckSyncWithRequest(req)
}

// UpgradeObCheckWithRequest returns a DagDetailDTO and an error, when the upgrade check task is requested successfully, the error will be nil.
// the parameter is a UpgradeObCheckRequest, which can be created by NewUpgradeObCheckRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeObCheckWithRequest(req *UpgradeObCheckRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createUpgradeObCheckResponse()
	err = c.Execute(req, response)
	dag = response.DagDetailDTO
	return
}

// UpgradeObCheckSyncWithRequest returns a DagDetailDTO and an error, when the upgrade check task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a UpgradeObCheckRequest, which can be created by NewUpgradeObCheckRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) UpgradeObCheckSyncWithRequest(req *UpgradeObCheckRequest) (*model.DagDetailDTO, error) {
	dag, err := c.UpgradeObCheckWithRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
