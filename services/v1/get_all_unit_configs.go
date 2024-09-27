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

type GetAllUnitConfigsRequest struct {
	*request.BaseRequest
}

// NewGetAllUnitConfigsRequest return a GetAllUnitConfigsRequest, which can be used as the argument for the GetAllUnitConfigsWithRequest.
func (c *Client) NewGetAllUnitConfigsRequest() *GetAllUnitConfigsRequest {
	req := &GetAllUnitConfigsRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/units/config", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type IterableData struct {
	Contents []model.ResourceUnitConfig `json:"contents"`
}

type getAllUnitConfigsResponse struct {
	*response.OcsAgentResponse
	*IterableData
}

func (c *Client) createGetAllUnitConfigsResponse() *getAllUnitConfigsResponse {
	resp := &getAllUnitConfigsResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		IterableData:     &IterableData{},
	}
	resp.Data = resp.IterableData
	return resp
}

// GetAllUnitConfigs returns a []ResourceUnitConfig and an error.
// If the error is non-nil, the []ResourceUnitConfig will be empty.
func (c *Client) GetAllUnitConfigs() (dags []model.ResourceUnitConfig, err error) {
	req := c.NewGetAllUnitConfigsRequest()
	return c.GetAllUnitConfigsWithRequest(req)
}

// GetAllUnitConfigsWithRequest returns a []ResourceUnitConfig and an error.
// If the error is non-nil, the []ResourceUnitConfig will be empty.
func (c *Client) GetAllUnitConfigsWithRequest(req *GetAllUnitConfigsRequest) (unitConfigs []model.ResourceUnitConfig, err error) {
	response := c.createGetAllUnitConfigsResponse()
	err = c.Execute(req, response)
	return response.IterableData.Contents, err
}
