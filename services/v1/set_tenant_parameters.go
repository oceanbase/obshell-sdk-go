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

type SetTenantParametersRequest struct {
	*request.BaseRequest
}

type setTenantParametersResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createSetTenantParametersResponse() *setTenantParametersResponse {
	return &setTenantParametersResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewSetTenantParametersRequest return a SetTenantParametersRequest, which can be used as the argument for the SetTenantParametersWithRequest.
// tenantName: the name of the tenant.
// parameters: the parameters of the tenant.
func (c *Client) NewSetTenantParametersRequest(tenantName string, parameters map[string]interface{}) *SetTenantParametersRequest {
	req := &SetTenantParametersRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetBody(map[string]map[string]interface{}{
		"parameters": parameters,
	})
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/parameters", tenantName), c.GetHost(), c.GetPort(), "PUT")
	return req
}

// SetTenantParameters sets the parameters of a tenant.
func (c *Client) SetTenantParameters(tenantName string, parameters map[string]interface{}) error {
	request := c.NewSetTenantParametersRequest(tenantName, parameters)
	return c.SetTenantParametersWithRequest(request)
}

// SetTenantParametersWithRequest sets the parameters of a tenant with a SetTenantParametersRequest.
func (c *Client) SetTenantParametersWithRequest(request *SetTenantParametersRequest) error {
	response := c.createSetTenantParametersResponse()
	return c.Execute(request, response)
}
