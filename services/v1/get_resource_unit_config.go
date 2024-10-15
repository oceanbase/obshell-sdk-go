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

type GetUnitConfigRequest struct {
	*request.BaseRequest
}

// NewGetUnitConfigRequest return a GetUnitConfigRequest, which can be used as the argument for the GetUnitConfigWithRequest.
func (c *Client) NewGetUnitConfigRequest(unitConfigName string) *GetUnitConfigRequest {
	req := &GetUnitConfigRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo(fmt.Sprintf("/api/v1/unit/config/%s", unitConfigName), c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type getUnitConfigResponse struct {
	*response.OcsAgentResponse
	*model.ResourceUnitConfig
}

func (c *Client) createGetUnitConfigResponse() *getUnitConfigResponse {
	resp := &getUnitConfigResponse{
		OcsAgentResponse:   response.NewOcsAgentResponse(),
		ResourceUnitConfig: &model.ResourceUnitConfig{},
	}
	resp.Data = resp.ResourceUnitConfig
	return resp
}

// GetUnitConfig returns a []ResourceUnitConfig and an error.
// If the error is non-nil, the []ResourceUnitConfig will be empty.
func (c *Client) GetUnitConfig(unitConfigName string) (*model.ResourceUnitConfig, error) {
	req := c.NewGetUnitConfigRequest(unitConfigName)
	return c.GetUnitConfigWithRequest(req)
}

// GetUnitConfigWithRequest returns a []ResourceUnitConfig and an error.
// If the error is non-nil, the []ResourceUnitConfig will be empty.
func (c *Client) GetUnitConfigWithRequest(req *GetUnitConfigRequest) (unitConfig *model.ResourceUnitConfig, err error) {
	response := c.createGetUnitConfigResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.ResourceUnitConfig, nil
}
