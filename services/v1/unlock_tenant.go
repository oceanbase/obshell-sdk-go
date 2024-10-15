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

type UnlockTenantRequest struct {
	*request.BaseRequest
}

type UnlockTenantResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createUnlockTenantResponse() *UnlockTenantResponse {
	return &UnlockTenantResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// NewUnlockTenantRequest return a UnlockTenantRequest, which can be used as the argument for the UnlockTenantWithRequest.
func (c *Client) NewUnlockTenantRequest(tenantName string) *UnlockTenantRequest {
	req := &UnlockTenantRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/lock", tenantName), c.GetHost(), c.GetPort(), "DELETE")
	req.SetAuthentication()
	return req
}

// UnlockTenantWithRequest unlocks a tenant.
// the parameter is a UnlockTenantRequest, which can be created by NewUnlockTenantRequest.
func (c *Client) UnlockTenantWithRequest(req *UnlockTenantRequest) error {
	response := c.createUnlockTenantResponse()
	if err := c.Execute(req, response); err != nil {
		return err
	}
	return nil
}

// UnlockTenant unlocks a tenant.
func (c *Client) UnlockTenant(tenantName string) error {
	req := c.NewUnlockTenantRequest(tenantName)
	return c.UnlockTenantWithRequest(req)
}
