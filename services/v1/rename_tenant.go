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

type RenameTenantRequest struct {
	*request.BaseRequest
}

type renameTenantResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createRenameTenantResponse() *renameTenantResponse {
	return &renameTenantResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewRenameTenantRequest return a RenameTenantRequest, which can be used as the argument for the RenameTenantWithRequest.
// tenantName: the name of the tenant.
// newName: the new name of the tenant.
func (c *Client) NewRenameTenantRequest(tenantName string, newName string) *RenameTenantRequest {
	req := &RenameTenantRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(
		map[string]string{
			"new_name": newName,
		},
	)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/name", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// RenameTenant creates a resource unit config.
// tenantName: the name of the tenant.
// newName: the new name of the tenant.
func (c *Client) RenameTenant(tenantName, newName string) error {
	request := c.NewRenameTenantRequest(tenantName, newName)
	return c.RenameTenantWithRequest(request)
}

// RenameTenantWithRequest renames a tenant with a RenameTenantRequest.
func (c *Client) RenameTenantWithRequest(request *RenameTenantRequest) error {
	response := c.createRenameTenantResponse()
	return c.Execute(request, response)
}
