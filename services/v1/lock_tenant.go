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

type LockTenantRequest struct {
	*request.BaseRequest
}

type LockTenantResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createLockTenantResponse() *LockTenantResponse {
	return &LockTenantResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
}

// NewLockTenantRequest return a LockTenantRequest, which can be used as the argument for the LockTenantWithRequest.
func (c *Client) NewLockTenantRequest(tenantName string) *LockTenantRequest {
	req := &LockTenantRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/tenant/%s/lock", tenantName), c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	return req
}

// LockTenantWithRequest locks a tenant.
// the parameter is a LockTenantRequest, which can be created by NewLockTenantRequest.
func (c *Client) LockTenantWithRequest(req *LockTenantRequest) error {
	response := c.createLockTenantResponse()
	if err := c.Execute(req, response); err != nil {
		return err
	}
	return nil
}

// LockTenant locks a tenant.
func (c *Client) LockTenant(tenantName string) error {
	req := c.NewLockTenantRequest(tenantName)
	return c.LockTenantWithRequest(req)
}
