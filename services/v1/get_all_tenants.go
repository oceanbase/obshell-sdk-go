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
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type GetAllTenantOverviewRequest struct {
	*request.BaseRequest
}

// NewGetAllTenantOverviewRequest return a GetAllTenantOverviewRequest, which can be used as the argument for the GetAllTenantOverviewWithRequest.
func (c *Client) NewGetAllTenantOverviewRequest() *GetAllTenantOverviewRequest {
	req := &GetAllTenantOverviewRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/tenants/overview", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getAllTenantOverviewResponse struct {
	*response.OcsAgentResponse
	*IteratorTenantOverview
}

type IteratorTenantOverview struct {
	Tenants []model.TenantOverview `json:"contents"`
}

func (c *Client) createGetAllTenantOverviewResponse() *getAllTenantOverviewResponse {
	resp := &getAllTenantOverviewResponse{
		OcsAgentResponse:       response.NewOcsAgentResponse(),
		IteratorTenantOverview: &IteratorTenantOverview{},
	}
	resp.Data = resp.IteratorTenantOverview
	return resp
}

// GetAllTenantOverview returns a []TenantOverview and an error.
// If the error is non-nil, the []TenantOverview will be empty.
func (c *Client) GetAllTenantOverview() (tenants []model.TenantOverview, err error) {
	req := c.NewGetAllTenantOverviewRequest()
	return c.GetAllTenantOverviewWithRequest(req)
}

// GetAllTenantOverviewWithRequest returns a []TenantOverview and an error.
// If the error is non-nil, the []TenantOverview will be empty.
func (c *Client) GetAllTenantOverviewWithRequest(req *GetAllTenantOverviewRequest) (tenants []model.TenantOverview, err error) {
	response := c.createGetAllTenantOverviewResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.Tenants, err
}
