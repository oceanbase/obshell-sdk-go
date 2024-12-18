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
	"fmt"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type DeleteZoneRequest struct {
	*request.BaseRequest
}

type deleteZoneResponse struct {
	*response.TaskResponse
}

func (c *Client) createDeleteZoneResponse() *deleteZoneResponse {
	return &deleteZoneResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewDeleteZoneRequest returns a DeleteZoneRequest, which can be used as the argument for the DeleteZoneWithRequest/DeleteZoneSyncWithRequest.
// zoneName: the name of zone to be deleted.
func (c *Client) NewDeleteZoneRequest(zoneName string) *DeleteZoneRequest {
	req := &DeleteZoneRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/zone/%s", zoneName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

// DeleteZone deletes a zone without unit from cluster.
func (c *Client) DeleteZone(zoneName string) (dag *model.DagDetailDTO, err error) {
	request := c.NewDeleteZoneRequest(zoneName)
	return c.DeleteZoneSyncWithRequest(request)
}

// DeleteZoneWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) DeleteZoneWithRequest(request *DeleteZoneRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createDeleteZoneResponse()
	if err = c.Execute(request, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// DeleteZoneSyncWithRequest returns a DagDetailDTO and an error.
// You can check or operater the task through the DagDetailDTO.
// If the zone is not exist in cluster, the DagDetailDTO will be nil.
func (c *Client) DeleteZoneSyncWithRequest(request *DeleteZoneRequest) (dag *model.DagDetailDTO, err error) {
	if dag, err = c.DeleteZoneWithRequest(request); err != nil {
		return nil, err
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
