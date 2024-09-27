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

type SetTenantVariablesRequest struct {
	*request.BaseRequest
}

type setTenantVariablesResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createSetTenantVariablesResponse() *setTenantVariablesResponse {
	return &setTenantVariablesResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewSetTenantVariablesRequest return a SetTenantVariablesRequest, which can be used as the argument for the SetTenantVariablesWithRequest.
// tenantName: the name of the tenant.
// variables: the variables of the tenant.
func (c *Client) NewSetTenantVariablesRequest(tenantName string, variables map[string]interface{}) *SetTenantVariablesRequest {
	req := &SetTenantVariablesRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(map[string]map[string]interface{}{
		"variables": variables,
	})
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/variables", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// SetTenantVariables sets the variables of a tenant.
func (c *Client) SetTenantVariables(tenantName string, variables map[string]interface{}) error {
	request := c.NewSetTenantVariablesRequest(tenantName, variables)
	return c.SetTenantVariablesWithRequest(request)
}

// SetTenantVariablesWithRequest sets the variables of a tenant.
func (c *Client) SetTenantVariablesWithRequest(request *SetTenantVariablesRequest) error {
	response := c.createSetTenantVariablesResponse()
	return c.Execute(request, response)
}
