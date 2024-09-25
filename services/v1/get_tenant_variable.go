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

type GetTenantVariableRequest struct {
	*request.BaseRequest
}

// NewGetTenantVariableRequest return a GetTenantVariableRequest, which can be used as the argument for the GetTenantVariableWithRequest.
func (c *Client) NewGetTenantVariableRequest(tenantName string, variableName string) *GetTenantVariableRequest {
	req := &GetTenantVariableRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/variable/%s", tenantName, variableName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getTenantVariableResponse struct {
	*response.OcsAgentResponse
	*model.VariableInfo
}

func (r *getTenantVariableResponse) Init() {
	r.Data = r.VariableInfo
}

func (c *Client) createGetTenantVariableResponse() *getTenantVariableResponse {
	return &getTenantVariableResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		VariableInfo:     &model.VariableInfo{},
	}
}

// GetTenantVariable returns a *VariableInfo and an error.
func (c *Client) GetTenantVariable(tenantName string, variableName string) (variable *model.VariableInfo, err error) {
	req := c.NewGetTenantVariableRequest(tenantName, variableName)
	return c.GetTenantVariableWithRequest(req)
}

// GetTenantVariableWithRequest returns a *VariableInfo and an error.
func (c *Client) GetTenantVariableWithRequest(req *GetTenantVariableRequest) (variable *model.VariableInfo, err error) {
	response := c.createGetTenantVariableResponse()
	err = c.Execute(req, response)
	return response.VariableInfo, err
}
