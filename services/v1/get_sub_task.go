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

type GetSubTaskRequest struct {
	*request.BaseRequest
}

// NewGetSubTaskRequest return a GetSubTaskRequest, which can be used as the argument for the GetSubTaskWithRequest.
func (c *Client) NewGetSubTaskRequest(taskId string) *GetSubTaskRequest {
	req := &GetSubTaskRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/task/sub_task/"+taskId, c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type GetSubTaskResponse struct {
	*response.OcsAgentResponse
	*model.TaskDetailDTO
}

func (c *Client) createGetSubTaskResponse() *GetSubTaskResponse {
	return &GetSubTaskResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// GetSubTask returns a TaskDetailDTO and an error.
// If the error is non-nil, the TaskDetailDTO will be nil.
// id is the id of the sub_task.
func (c *Client) GetSubTask(id string) (*model.TaskDetailDTO, error) {
	req := c.NewGetSubTaskRequest(id)
	return c.GetSubTaskWithRequest(req)
}

// GetSubTaskWithRequest returns a TaskDetailDTO and an error.
// The parameter is a GetSubTaskRequest, which can be created by NewGetSubTaskRequest.
// If the error is non-nil, the TaskDetailDTO will be nil.
func (c *Client) GetSubTaskWithRequest(req *GetSubTaskRequest) (SubTask *model.TaskDetailDTO, err error) {
	response := c.createGetSubTaskResponse()
	err = c.Execute(req, response)
	return response.TaskDetailDTO, err
}
