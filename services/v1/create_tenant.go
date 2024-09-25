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
	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type CreateTenantParam struct {
	Name         string                 `json:"name" binding:"required"` // Tenant name.
	Mode         string                 `json:"mode"`                    // Tenant mode, "MYSQL"(default) or "ORACLE".
	PrimaryZone  string                 `json:"primary_zone"`            // Tenant primary_zone.
	Whilelist    string                 `json:"whitelist"`
	RootPassword string                 `json:"root_password"` // Root password.
	Scenario     string                 `json:"scenario"`      // Can be one of 'express_oltp', 'complex_oltp', 'olap', 'htap'(default), 'kv'.
	Charset      string                 `json:"charset"`
	Collation    string                 `json:"collation"`
	ReadOnly     bool                   `json:"read_only"`                    // Default to false.
	Comment      string                 `json:"comment"`                      // Messages of tenant.
	Variables    map[string]interface{} `json:"variables"`                    // Tenant global variables.
	Parameters   map[string]interface{} `json:"parameters"`                   // Tenant parameters.
	ZoneList     []ZoneParam            `json:"zone_list" binding:"required"` // Tenant zone list with unit config.
}

type ZoneParam struct {
	Name           string `json:"name" binding:"required"`
	UnitConfigName string `json:"unit_config_name" binding:"required"`
	UnitNum        int    `json:"unit_num" binding:"required"`
	ReplicaType    string `json:"replica_type"` // Replica type, "FULL"(default) or "READONLY"
}

type CreateTenantRequest struct {
	*request.BaseRequest
	param CreateTenantParam
}

func (req *CreateTenantRequest) SetMode(mode string) *CreateTenantRequest {
	req.param.Mode = mode
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetPrimaryZone(primaryZone string) *CreateTenantRequest {
	req.param.PrimaryZone = primaryZone
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetWhitelist(whitelist string) *CreateTenantRequest {
	req.param.Whilelist = whitelist
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetRootPassword(rootPassword string) *CreateTenantRequest {
	req.param.RootPassword = rootPassword
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetCharset(charset string) *CreateTenantRequest {
	req.param.Charset = charset
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetCollation(collation string) *CreateTenantRequest {
	req.param.Collation = collation
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetReadOnly(readOnly bool) *CreateTenantRequest {
	req.param.ReadOnly = readOnly
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetComment(comment string) *CreateTenantRequest {
	req.param.Comment = comment
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetVariables(variables map[string]interface{}) *CreateTenantRequest {
	req.param.Variables = variables
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetParameters(parameters map[string]interface{}) *CreateTenantRequest {
	req.param.Parameters = parameters
	req.SetBody(&req.param)
	return req
}

func (req *CreateTenantRequest) SetScenario(scenario string) *CreateTenantRequest {
	req.param.Scenario = scenario
	req.SetBody(&req.param)
	return req
}

type createTenantResponse struct {
	*response.TaskResponse
}

func (c *Client) newCreateTenantResponse() *createTenantResponse {
	return &createTenantResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewCreateTenantRequest return a CreateTenantRequest, which can be used as the argument for the CreateTenantWithRequest/CreateTenantSyncWithRequest.
// name: the name of the tenant.
// zoneList: the zone list with replicas properties.
// You can use other SetXXX functions to set other parameters.
func (c *Client) NewCreateTenantRequest(name string, zoneList []ZoneParam) *CreateTenantRequest {
	req := &CreateTenantRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: CreateTenantParam{
			Name:     name,
			ZoneList: zoneList,
		},
	}
	req.SetBody(&req.param)
	req.SetAuthentication()
	req.InitApiInfo("/api/v1/tenant", c.GetHost(), c.GetPort(), "POST")
	return req
}

// CreateTenant returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// name: the name of the tenant.
// zoneList: the zone list with replicas properties.
func (c *Client) CreateTenant(name string, zoneList []ZoneParam) (*model.DagDetailDTO, error) {
	request := c.NewCreateTenantRequest(name, zoneList)
	return c.CreateTenantSyncWithRequest(request)
}

// CreateTenantWithRequest returns a DagDetailDTO and an error, when the task is requested successfully, the error will be nil.
// the parameter is a CreateTenantRequest, which can be created by NewCreateTenantRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) CreateTenantWithRequest(request *CreateTenantRequest) (dag *model.DagDetailDTO, err error) {
	response := c.newCreateTenantResponse()
	err = c.Execute(request, response)
	dag = response.DagDetailDTO
	return
}

// CreateTenantSyncWithRequest returns a DagDetailDTO and an error, when the task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a CreateTenantRequest, which can be created by NewCreateTenantRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) CreateTenantSyncWithRequest(request *CreateTenantRequest) (*model.DagDetailDTO, error) {
	dag, err := c.CreateTenantWithRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
