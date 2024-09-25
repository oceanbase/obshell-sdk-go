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

	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type SetTenantRootPasswordRequest struct {
	*request.BaseRequest
}

type setTenantRootPasswordResponse struct {
	*response.OcsAgentResponse
}

func (r *setTenantRootPasswordResponse) Init() {}

func (c *Client) createSetTenantRootPasswordResponse() *setTenantRootPasswordResponse {
	return &setTenantRootPasswordResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewSetTenantRootPasswordRequest return a SetTenantRootPasswordRequest, which can be used as the argument for the SetTenantRootPasswordWithRequest.
// tenantName: the name of the tenant, could not be sys.
// oldPassword: the old password of the tenant.
// newPassword: the new password of the tenant.
func (c *Client) NewSetTenantRootPasswordRequest(tenantName string, oldPassword string, newPassword string) *SetTenantRootPasswordRequest {
	req := &SetTenantRootPasswordRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(
		map[string]string{
			"old_password": oldPassword,
			"new_password": newPassword,
		},
	)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/password", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// SetTenantRootPassword sets the root password of a tenant.
func (c *Client) SetTenantRootPassword(tenantName, oldPassword, newPassword string) error {
	request := c.NewSetTenantRootPasswordRequest(tenantName, oldPassword, newPassword)
	return c.SetTenantRootPasswordWithRequest(request)
}

// SetTenantRootPasswordWithRequest sets the root password of a tenant with a SetTenantRootPasswordRequest.
func (c *Client) SetTenantRootPasswordWithRequest(request *SetTenantRootPasswordRequest) error {
	response := c.createSetTenantRootPasswordResponse()
	return c.Execute(request, response)
}
