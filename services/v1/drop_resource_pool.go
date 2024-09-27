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

type DropResourcePoolRequest struct {
	*request.BaseRequest
}

type DropResourcePoolResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createDropResourcePoolResponse() *DropResourcePoolResponse {
	return &DropResourcePoolResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewDropResourcePoolRequest return a DropResourcePoolRequest, which can be used as the argument for the DropResourcePoolWithRequest.
// poolName: the name of the resource pool.
func (c *Client) NewDropResourcePoolRequest(poolName string) *DropResourcePoolRequest {
	req := &DropResourcePoolRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/resource-pool/%s", poolName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

// DropResourcePool drops a resource pool.
// poolName: the name of the resource pool.
func (c *Client) DropResourcePool(poolName string) error {
	request := c.NewDropResourcePoolRequest(poolName)
	return c.DropResourcePoolWithRequest(request)
}

// DropResourcePoolWithRequest drops a resource pool with a DropResourcePoolRequest.
func (c *Client) DropResourcePoolWithRequest(request *DropResourcePoolRequest) error {
	response := c.createDropResourcePoolResponse()
	return c.Execute(request, response)
}
