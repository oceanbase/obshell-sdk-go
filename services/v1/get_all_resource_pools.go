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

type GetAllResourcePoolsRequest struct {
	*request.BaseRequest
}

// NewGetAllResourcePoolsRequest return a GetAllResourcePoolsRequest, which can be used as the argument for the GetAllResourcePoolsWithRequest.
func (c *Client) NewGetAllResourcePoolsRequest() *GetAllResourcePoolsRequest {
	req := &GetAllResourcePoolsRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	req.InitApiInfo("/api/v1/resource-pools", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type IterablePools struct {
	Contents []model.ResourcePoolInfo `json:"contents"`
}

type getAllResourcePoolsResponse struct {
	*response.OcsAgentResponse
	*IterablePools
}

func (r *getAllResourcePoolsResponse) Init() {
	r.Data = r.IterablePools
}

func (c *Client) createGetAllResourcePoolsResponse() *getAllResourcePoolsResponse {
	return &getAllResourcePoolsResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
		IterablePools:    &IterablePools{},
	}
}

// GetAllResourcePools returns a []ResourcePoolInfo and an error.
// If the error is non-nil, the []ResourcePoolInfo will be empty.
func (c *Client) GetAllResourcePools() (dags []model.ResourcePoolInfo, err error) {
	req := c.NewGetAllResourcePoolsRequest()
	return c.GetAllResourcePoolsWithRequest(req)
}

// GetAllResourcePoolsWithRequest returns a []ResourcePoolInfo and an error.
// If the error is non-nil, the []*ResourcePoolInfo will be empty.
func (c *Client) GetAllResourcePoolsWithRequest(req *GetAllResourcePoolsRequest) (dags []model.ResourcePoolInfo, err error) {
	response := c.createGetAllResourcePoolsResponse()
	err = c.Execute(req, response)
	return response.Contents, err
}
