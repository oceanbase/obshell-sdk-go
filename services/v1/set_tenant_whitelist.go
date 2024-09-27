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

type SetTenantWhitelistRequest struct {
	*request.BaseRequest
}

type setTenantWhitelistResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createSetTenantWhitelistResponse() *setTenantWhitelistResponse {
	return &setTenantWhitelistResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewSetTenantWhitelistRequest return a SetTenantWhitelistRequest, which can be used as the argument for the SetTenantWhitelistWithRequest.
// tenantName: the name of the tenant.
// whitelist: the whitelist of the tenant.
func (c *Client) NewSetTenantWhitelistRequest(tenantName, whitelist string) *SetTenantWhitelistRequest {
	req := &SetTenantWhitelistRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(
		map[string]string{
			"whitelist": whitelist,
		},
	)
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/whitelist", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// SetTenantWhitelist sets the whitelist of a tenant.
// tenantName: the name of the tenant.
// whitelist: the whitelist of the tenant.
func (c *Client) SetTenantWhitelist(tenantName, whitelist string) error {
	request := c.NewSetTenantWhitelistRequest(tenantName, whitelist)
	return c.SetTenantWhitelistWithRequest(request)
}

// SetTenantWhitelistWithRequest sets the whitelist of a tenant with a SetTenantWhitelistRequest.
func (c *Client) SetTenantWhitelistWithRequest(request *SetTenantWhitelistRequest) error {
	response := c.createSetTenantWhitelistResponse()
	return c.Execute(request, response)
}
