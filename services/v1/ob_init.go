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

type InitRequest struct {
	*request.BaseRequest
}

type InitResponse struct {
	*response.TaskResponse
}

func (c *Client) createInitResponse() *InitResponse {
	return &InitResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NwqInitRequest return a InitRequest, which can be used as the argument for the InitWithRequest/InitSyncWithRequest.
func (c *Client) NewInitRequest() *InitRequest {
	req := &InitRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.InitApiInfo("/api/v1/ob/init", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	return req
}

// InitWithRequest returns a DagDetailDTO and an error, when the init task is requested successfully, the error will be nil.
// the parameter is a InitRequest, which can be created by NewInitRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) InitWithRequest(req *InitRequest) (*model.DagDetailDTO, error) {
	response := c.createInitResponse()
	err := c.Execute(req, response)
	if err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// Init returns a DagDetailDTO and an error, when the init task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) Init() (*model.DagDetailDTO, error) {
	req := c.NewInitRequest()
	return c.InitSyncWithRequest(req)
}

// InitSyncWithRequest returns a DagDetailDTO and an error, when the init task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) InitSyncWithRequest(req *InitRequest) (*model.DagDetailDTO, error) {
	dag, err := c.InitWithRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
