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

type GetTenantVariablesRequest struct {
	*request.BaseRequest
}

// NewGetTenantVariablesRequest return a GetTenantVariablesRequest, which can be used as the argument for the GetTenantVariablesWithRequest.
// tenantName: the name of the tenant.
// filter: the filter of the tenant variables. If not set, the default value is "%".
func (c *Client) NewGetTenantVariablesRequest(tenantName string, filter ...string) *GetTenantVariablesRequest {
	req := &GetTenantVariablesRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/variables", tenantName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type IterableVariables struct {
	Contents []model.VariableInfo `json:"contents"`
}

type getTenantVariablesResponse struct {
	*response.OcsAgentResponse
	*IterableVariables
}

func (c *Client) createGetTenantVariablesResponse() *getTenantVariablesResponse {
	resp := &getTenantVariablesResponse{
		OcsAgentResponse:  response.NewOcsAgentResponse(),
		IterableVariables: &IterableVariables{},
	}
	resp.Data = resp.IterableVariables
	return resp
}

// GetTenantVariables returns a []VariableInfo and an error.
// If the error is non-nil, the []VariableInfo will be empty.
// tenantName: the name of the tenant.
// filter: the filter of the tenant variables. If not set, the default value is "%".
func (c *Client) GetTenantVariables(tenantName string, filter ...string) (variables []model.VariableInfo, err error) {
	req := c.NewGetTenantVariablesRequest(tenantName, filter...)
	return c.GetTenantVariablesWithRequest(req)
}

// GetTenantVariablesWithRequest returns a []VariableInfo and an error.
// If the error is non-nil, the []VariableInfo will be empty.
func (c *Client) GetTenantVariablesWithRequest(req *GetTenantVariablesRequest) (variables []model.VariableInfo, err error) {
	response := c.createGetTenantVariablesResponse()
	err = c.Execute(req, response)
	return response.Contents, err
}
