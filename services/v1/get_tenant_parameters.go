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

type GetTenantParametersRequest struct {
	*request.BaseRequest
}

// NewGetTenantParametersRequest return a GetTenantParametersRequest, which can be used as the argument for the GetTenantParametersWithRequest.
// tenantName: the name of the tenant.
// filter: the filter of the parameters. If not set, the default value is "%".
func (c *Client) NewGetTenantParametersRequest(tenantName string, filter ...string) *GetTenantParametersRequest {
	req := &GetTenantParametersRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	if len(filter) > 0 {
		req.SetQueryParam("filter", filter[0])
	} else {
		req.SetQueryParam("filter", "%")
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/parameters", tenantName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type IterableParameters struct {
	Contents []model.ParameterInfo `json:"contents"`
}

type getTenantParametersResponse struct {
	*response.OcsAgentResponse
	*IterableParameters
}

func (c *Client) createGetTenantParametersResponse() *getTenantParametersResponse {
	resp := &getTenantParametersResponse{
		OcsAgentResponse:   response.NewOcsAgentResponse(),
		IterableParameters: &IterableParameters{},
	}
	resp.Data = resp.IterableParameters
	return resp
}

// GetTenantParameters returns a []ParameterInfo and an error.
// If the error is non-nil, the []ParameterInfo will be empty.
// tenantName: the name of the tenant.
// filter: the filter of the parameters. If not set, the default value is "%".
func (c *Client) GetTenantParameters(tenantName string, filter ...string) (dags []model.ParameterInfo, err error) {
	req := c.NewGetTenantParametersRequest(tenantName, filter...)
	return c.GetTenantParametersWithRequest(req)
}

// GetTenantParametersWithRequest returns a []ParameterInfo and an error.
// If the error is non-nil, the []ParameterInfo will be empty.
func (c *Client) GetTenantParametersWithRequest(req *GetTenantParametersRequest) (unitConfigs []model.ParameterInfo, err error) {
	response := c.createGetTenantParametersResponse()
	err = c.Execute(req, response)
	return response.Contents, err
}
