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

type DropResourceUnitConfigRequest struct {
	*request.BaseRequest
}

type DropResourceUnitConfigResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createDropResourceUnitConfigResponse() *DropResourceUnitConfigResponse {
	return &DropResourceUnitConfigResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

// NewDropResourceUnitConfigRequest return a DropResourceUnitConfigRequest, which can be used as the argument for the DropResourceUnitConfigWithRequest.
// unitConfigName: the name of the resource unit config.
func (c *Client) NewDropResourceUnitConfigRequest(unitConfigName string) *DropResourceUnitConfigRequest {
	req := &DropResourceUnitConfigRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.SetAuthentication()
	req.InitApiInfo(fmt.Sprintf("/api/v1/unit/config/%s", unitConfigName), c.GetHost(), c.GetPort(), "DELETE")
	return req
}

// DropResourceUnitConfig creates a resource unit config.
// unitConfigName: the name of the resource unit config.
func (c *Client) DropResourceUnitConfig(unitConfigName string) error {
	request := c.NewDropResourceUnitConfigRequest(unitConfigName)
	return c.DropResourceUnitConfigWithRequest(request)
}

// DropResourceUnitConfigWithRequest creates a resource unit config with the DropResourceUnitConfigRequest.
func (c *Client) DropResourceUnitConfigWithRequest(request *DropResourceUnitConfigRequest) error {
	response := c.createDropResourceUnitConfigResponse()
	return c.Execute(request, response)
}
