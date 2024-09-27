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

type GetTenantInfoRequest struct {
	*request.BaseRequest
}

// NewGetTenantInfoRequest return a GetTenantInfoRequest, which can be used as the argument for the GetTenantInfoWithRequest.
func (c *Client) NewGetTenantInfoRequest(tenantName string) *GetTenantInfoRequest {
	req := &GetTenantInfoRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s", tenantName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getTenantInfoResponse struct {
	*response.OcsAgentResponse
	*model.TenantInfo
}

func (c *Client) createGetTenantInfoResponse() *getTenantInfoResponse {
	resp := &getTenantInfoResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		TenantInfo:       &model.TenantInfo{},
	}
	resp.Data = resp.TenantInfo
	return resp
}

// GetTenantInfo returns a *TenantInfoInfo and an error.
func (c *Client) GetTenantInfo(tenantName string) (tenants *model.TenantInfo, err error) {
	req := c.NewGetTenantInfoRequest(tenantName)
	return c.GetTenantInfoWithRequest(req)
}

// GetTenantInfoWithRequest returns a *TenantInfoInfo and an error.
func (c *Client) GetTenantInfoWithRequest(req *GetTenantInfoRequest) (tenants *model.TenantInfo, err error) {
	response := c.createGetTenantInfoResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.TenantInfo, err
}
