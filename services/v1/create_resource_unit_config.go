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
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type CreateResourceUnitConfigRequest struct {
	*request.BaseRequest
	param CreateResourceUnitConfigParam
}

type CreateResourceUnitConfigParam struct {
	Name        string   `json:"name" binding:"required"`
	MemorySize  string   `json:"memory_size" binding:"required"` // MemorySize should be greater than or equal to '1G'.
	MaxCpu      float64  `json:"max_cpu" binding:"required"`     // MaxCpu should be greater than 0.
	MinCpu      *float64 `json:"min_cpu"`                        // MinCpu should be smaller than or equal MaxCpu.
	MaxIops     *int     `json:"max_iops"`                       // MaxIops should be greater than or equal to 1024.
	MinIops     *int     `json:"min_iops"`                       // MinIops should be smaller than or equal to MaxIops.
	LogDiskSize *string  `json:"log_disk_size"`                  // LogDiskSize should be greater than or equal to '2G'.
}

type createResourceUnitConfigResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createCreateResourceUnitConfigResponse() *createResourceUnitConfigResponse {
	return &createResourceUnitConfigResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

func (r *createResourceUnitConfigResponse) Init() {}

// NewCreateResourceUnitConfigRequest return a CreateResourceUnitConfigRequest, which can be used as the argument for the CreateResourceUnitConfigWithRequest.
// unitConfigName: the name of the resource unit config.
// memorySize: the memory size of the resource unit config.
// maxCpu: the max cpu cores of the resource unit config, greater than 1.
// You can set the minCpu, maxIops, minIops, logDiskSize by calling SetMinCpu, SetMaxIops, SetMinIops, SetLogDiskSize.
func (c *Client) NewCreateResourceUnitConfigRequest(unitConfigName string, memorySize string, maxCpu float64) *CreateResourceUnitConfigRequest {
	req := &CreateResourceUnitConfigRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: CreateResourceUnitConfigParam{
			Name:       unitConfigName,
			MemorySize: memorySize,
			MaxCpu:     maxCpu,
		},
	}
	req.InitApiInfo("/api/v1/unit/config", c.GetHost(), c.GetPort(), "POST")
	req.SetBody(&req.param)
	req.SetAuthentication()
	return req
}

func (r *CreateResourceUnitConfigRequest) SetMinCpu(minCpu float64) *CreateResourceUnitConfigRequest {
	r.param.MinCpu = &minCpu
	r.SetBody(r.param)
	return r
}

func (r *CreateResourceUnitConfigRequest) SetMaxIops(maxIops int) *CreateResourceUnitConfigRequest {
	r.param.MaxIops = &maxIops
	r.SetBody(r.param)
	return r
}

func (r *CreateResourceUnitConfigRequest) SetMinIops(minIops int) *CreateResourceUnitConfigRequest {
	r.param.MinIops = &minIops
	r.SetBody(r.param)
	return r
}

func (r *CreateResourceUnitConfigRequest) SetLogDiskSize(logDiskSize string) *CreateResourceUnitConfigRequest {
	r.param.LogDiskSize = &logDiskSize
	r.SetBody(r.param)
	return r
}

// CreateResourceUnitConfig creates a resource unit config.
// unitConfigName: the name of the resource unit config.
// memorySize: the memory size of the resource unit config.
// maxCpu: the max cpu cores of the resource unit config, greater than 1.
func (c *Client) CreateResourceUnitConfig(unitConfigName string, memorySize string, maxCpu float64) error {
	request := c.NewCreateResourceUnitConfigRequest(unitConfigName, memorySize, maxCpu)
	return c.CreateResourceUnitConfigWithRequest(request)
}

// CreateResourceUnitConfigWithRequest creates a resource unit config with the CreateResourceUnitConfigRequest.
func (c *Client) CreateResourceUnitConfigWithRequest(request *CreateResourceUnitConfigRequest) error {
	response := c.createCreateResourceUnitConfigResponse()
	return c.Execute(request, response)
}
