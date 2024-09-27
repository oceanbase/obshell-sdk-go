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
	"github.com/pkg/errors"
)

type SetTenantPrimaryZoneRequest struct {
	*request.BaseRequest
}

type setTenantPrimaryZoneResponse struct {
	*response.TaskResponse
}

func (c *Client) createSetTenantPrimaryZoneResponse() *setTenantPrimaryZoneResponse {
	return &setTenantPrimaryZoneResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewSetTenantPrimaryZoneRequest return a SetTenantPrimaryZoneRequest, which can be used as the argument for the SetTenantPrimaryZoneWithRequest.
// tenantName: the name of the tenant.
// primaryZone: the primary zone of the tenant.
func (c *Client) NewSetTenantPrimaryZoneRequest(tenantName string, primaryZone string) *SetTenantPrimaryZoneRequest {
	req := &SetTenantPrimaryZoneRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(
		map[string]string{
			"primary_zone": primaryZone,
		},
	)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/primary-zone", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// SetTenantPrimaryZone sets the primary zone of a tenant.
func (c *Client) SetTenantPrimaryZone(tenantName, newName string) (*model.DagDetailDTO, error) {
	request := c.NewSetTenantPrimaryZoneRequest(tenantName, newName)
	return c.SetTenantPrimaryZoneSyncWithRequest(request)
}

// SetTenantPrimaryZoneSyncWithRequest sets the primary zone of a tenant with a SetTenantPrimaryZoneRequest.
func (c *Client) SetTenantPrimaryZoneSyncWithRequest(request *SetTenantPrimaryZoneRequest) (*model.DagDetailDTO, error) {
	dag, err := c.SetTenantPrimaryZoneWithRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}

// SetTenantPrimaryZoneWithRequest sets the primary zone of a tenant with a SetTenantPrimaryZoneRequest.
func (c *Client) SetTenantPrimaryZoneWithRequest(request *SetTenantPrimaryZoneRequest) (*model.DagDetailDTO, error) {
	response := c.createSetTenantPrimaryZoneResponse()
	err := c.Execute(request, response)
	dag := response.DagDetailDTO
	return dag, err
}
