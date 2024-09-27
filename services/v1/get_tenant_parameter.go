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

type GetTenantParameterRequest struct {
	*request.BaseRequest
}

// NewGetTenantParameterRequest return a GetTenantParameterRequest, which can be used as the argument for the GetTenantParameterWithRequest.
func (c *Client) NewGetTenantParameterRequest(tenantName string, parameterName string) *GetTenantParameterRequest {
	req := &GetTenantParameterRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/parameter/%s", tenantName, parameterName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getTenantParameterResponse struct {
	*response.OcsAgentResponse
	*model.ParameterInfo
}

func (c *Client) createGetTenantParameterResponse() *getTenantParameterResponse {
	resp := &getTenantParameterResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		ParameterInfo:    &model.ParameterInfo{},
	}
	resp.Data = resp.ParameterInfo
	return resp
}

// GetTenantParameter returns a *ParameterInfo and an error.
// If the error is non-nil, the *ParameterInfo will be nil.
func (c *Client) GetTenantParameter(tenantName string, parameterName string) (paramerter *model.ParameterInfo, err error) {
	req := c.NewGetTenantParameterRequest(tenantName, parameterName)
	return c.GetTenantParameterWithRequest(req)
}

// GetTenantParameterWithRequest returns a *ParameterInfo and an error.
// If the error is non-nil, the *ParameterInfo will be nil.
func (c *Client) GetTenantParameterWithRequest(req *GetTenantParameterRequest) (paramerter *model.ParameterInfo, err error) {
	response := c.createGetTenantParameterResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.ParameterInfo, err
}
