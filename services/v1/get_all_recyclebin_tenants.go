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

type GetAllRecyclebinTenantsRequest struct {
	*request.BaseRequest
}

// NewGetAllRecyclebinTenantsRequest return a GetAllRecyclebinTenantsRequest, which can be used as the argument for the GetAllRecyclebinTenantsWithRequest.
func (c *Client) NewGetAllRecyclebinTenantsRequest() *GetAllRecyclebinTenantsRequest {
	req := &GetAllRecyclebinTenantsRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/recyclebin/tenants", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type IteratorRecyclebinTenants struct {
	Contents []model.RecycledTenantOverView `json:"contents"`
}

type getAllRecyclebinTenantsResponse struct {
	*response.OcsAgentResponse
	*IteratorRecyclebinTenants
}

func (r *getAllRecyclebinTenantsResponse) Init() {
	r.Data = r.IteratorRecyclebinTenants
}
func (c *Client) createGetAllRecyclebinTenantsResponse() *getAllRecyclebinTenantsResponse {
	return &getAllRecyclebinTenantsResponse{
		OcsAgentResponse:          response.NewOcsAgentResponse(),
		IteratorRecyclebinTenants: &IteratorRecyclebinTenants{},
	}
}

// GetAllRecyclebinTenants returns a []RecyclebinTenantInfo and an error.
// If the error is non-nil, the []RecyclebinTenantInfo will be empty.
func (c *Client) GetAllRecyclebinTenants() (tenants []model.RecycledTenantOverView, err error) {
	req := c.NewGetAllRecyclebinTenantsRequest()
	return c.GetAllRecyclebinTenantsWithRequest(req)
}

// GetAllRecyclebinTenantsWithRequest returns a []RecyclebinTenantInfo and an error.
// If the error is non-nil, the []RecycledTenantOverView will be empty.
func (c *Client) GetAllRecyclebinTenantsWithRequest(req *GetAllRecyclebinTenantsRequest) (tenants []model.RecycledTenantOverView, err error) {
	response := c.createGetAllRecyclebinTenantsResponse()
	err = c.Execute(req, response)
	return response.Contents, err
}
