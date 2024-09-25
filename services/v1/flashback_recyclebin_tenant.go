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

type FlashbackRecyclebinTenantRequest struct {
	*request.BaseRequest
}

type FlashbackRecyclebinTenantResponse struct {
	*response.OcsAgentResponse
}

func (r *FlashbackRecyclebinTenantResponse) Init() {}

func (c *Client) createFlashbackRecyclebinTenantResponse() *FlashbackRecyclebinTenantResponse {
	return &FlashbackRecyclebinTenantResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewFlashbackRecyclebinTenantRequest return a FlashbackRecyclebinTenantRequest, which can be used as the argument for the FlashbackRecyclebinTenantWithRequest.
// objectOrOriginalName: the name of the object in recyclebin or the original name of the tenant.
// newName: the new name of the tenant.
func (c *Client) NewFlashbackRecyclebinTenantRequest(objectOrOriginalName string, newName ...string) *FlashbackRecyclebinTenantRequest {
	req := &FlashbackRecyclebinTenantRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetAuthentication()
	if len(newName) > 0 {
		req.SetBody(map[string]string{
			"new_name": newName[0],
		})
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/recyclebin/tenant/%s", objectOrOriginalName), c.GetHost(), c.GetPort(), "POST")
	return req
}

// FlashbackRecyclebinTenant will flashback a tenant from recyclebin.
// objectOrOriginalName: the name of the object in recyclebin or the original name of the tenant.
// newName: the new name of the tenant.
func (c *Client) FlashbackRecyclebinTenant(objectOrOriginalName string, newName ...string) error {
	request := c.NewFlashbackRecyclebinTenantRequest(objectOrOriginalName, newName...)
	return c.FlashbackRecyclebinTenantWithRequest(request)
}

// FlashbackRecyclebinTenantWithRequest will flashback a tenant from recyclebin with a FlashbackRecyclebinTenantRequest.
func (c *Client) FlashbackRecyclebinTenantWithRequest(request *FlashbackRecyclebinTenantRequest) error {
	response := c.createFlashbackRecyclebinTenantResponse()
	return c.Execute(request, response)
}
