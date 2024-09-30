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

type TenantConfigApiParam struct {
	DataBaseUri    string `json:"data_base_uri"`
	ArchiveBaseUri string `json:"archive_base_uri"`
	baseBackupConfigParam
}

type TenantBackupConfigRequest struct {
	*request.BaseRequest
	canPatchUri bool
	param       *TenantConfigApiParam
}

// NewTenantBackupConfigPostRequest creates a POST request to configure tenant backup settings.
func (c *Client) NewTenantBackupConfigPostRequest(tenantName, dataBaseUri, archiveBaseUri string) *TenantBackupConfigRequest {
	return c.newTenantBackupConfigRequest(tenantName, dataBaseUri, archiveBaseUri, "POST")
}

// NewTenantBackupConfigPatchRequest creates a PATCH request to update tenant backup settings.
func (c *Client) NewTenantBackupConfigPatchRequest(tenantName string) *TenantBackupConfigRequest {
	return c.newTenantBackupConfigRequest(tenantName, "", "", "PATCH")
}

func (c *Client) newTenantBackupConfigRequest(tenantName, dataBaseUri, archiveBaseUri, method string) *TenantBackupConfigRequest {
	req := &TenantBackupConfigRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &TenantConfigApiParam{
			DataBaseUri:    dataBaseUri,
			ArchiveBaseUri: archiveBaseUri,
		},
	}
	if method == "PATCH" {
		req.canPatchUri = true
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/backup/config", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), method)
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *TenantBackupConfigRequest) SetDataBaseUri(dataBaseUri string) *TenantBackupConfigRequest {
	if r.canPatchUri {
		r.param.DataBaseUri = dataBaseUri
		r.SetBody(r.param)
	}
	return r
}

func (r *TenantBackupConfigRequest) SetArchiveBaseUri(archiveBaseUri string) *TenantBackupConfigRequest {
	if r.canPatchUri {
		r.param.ArchiveBaseUri = archiveBaseUri
		r.SetBody(r.param)
	}
	return r
}

func (r *TenantBackupConfigRequest) SetLogArchiveConcurrency(logArchiveConcurrency int) *TenantBackupConfigRequest {
	r.param.LogArchiveConcurrency = &logArchiveConcurrency
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupConfigRequest) SetBinding(binding string) *TenantBackupConfigRequest {
	r.param.Binding = &binding
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupConfigRequest) SetHaLowThreadScore(haLowThreadScore int) *TenantBackupConfigRequest {
	r.param.HaLowThreadScore = &haLowThreadScore
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupConfigRequest) SetPieceSwitchInterval(pieceSwitchInterval string) *TenantBackupConfigRequest {
	r.param.PieceSwitchInterval = &pieceSwitchInterval
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupConfigRequest) SetArchiveLagTarget(archiveLagTarget string) *TenantBackupConfigRequest {
	r.param.ArchiveLagTarget = &archiveLagTarget
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupConfigRequest) SetDeletePolicy(policy, recoveryWindow string) *TenantBackupConfigRequest {
	r.param.DeletePolicy = &DeletePolicy{
		Policy:         policy,
		RecoveryWindow: recoveryWindow,
	}
	r.SetBody(r.param)
	return r
}

// TenantBackupConfig sends a POST request to configure tenant backup settings.
func (c *Client) TenantBackupConfig(tenantName, dataBaseUri, archiveBaseUri string) (*model.DagDetailDTO, error) {
	req := c.NewTenantBackupConfigPostRequest(tenantName, dataBaseUri, archiveBaseUri)
	return c.TenantBackupConfigSyncWithRequest(req)
}

// TenantBackupConfigWithRequest executes the request for configuring tenant backup settings.
func (c *Client) TenantBackupConfigWithRequest(req *TenantBackupConfigRequest) (*model.DagDetailDTO, error) {
	response := c.createPostTenantBackupConfigResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// TenantBackupConfigSyncWithRequest synchronously executes the tenant backup configuration request.
func (c *Client) TenantBackupConfigSyncWithRequest(req *TenantBackupConfigRequest) (*model.DagDetailDTO, error) {
	dag, err := c.TenantBackupConfigWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}

type TenantBackupConfigResponse struct {
	*response.TaskResponse
}

func (c *Client) createPostTenantBackupConfigResponse() *TenantBackupConfigResponse {
	return &TenantBackupConfigResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

type TenantBackupRequest struct {
	*request.BaseRequest
	param *TenantBackupApiParam
}

type TenantBackupApiParam struct {
	baseBackupApiParam
}

// NewTenantBackupRequest creates a new request for initiating a tenant backup.
func (c *Client) NewTenantBackupRequest(tenantName string) *TenantBackupRequest {
	req := &TenantBackupRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &TenantBackupApiParam{},
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/backup", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *TenantBackupRequest) SetBackupMode(mode string) *TenantBackupRequest {
	r.param.Mode = &mode
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupRequest) SetPlusArchive(plusArchive bool) *TenantBackupRequest {
	r.param.PlusArchive = &plusArchive
	r.SetBody(r.param)
	return r
}

func (r *TenantBackupRequest) SetEncryption(encryption string) *TenantBackupRequest {
	r.param.Encryption = &encryption
	r.SetBody(r.param)
	return r
}

// PostTenantBackup sends a POST request to start the tenant backup process.
func (c *Client) PostTenantBackup(tenantName string) (*model.DagDetailDTO, error) {
	req := c.NewTenantBackupRequest(tenantName)
	return c.TenantBackupSyncWithRequest(req)
}

// TenantBackupWithRequest executes the tenant backup request.
func (c *Client) TenantBackupWithRequest(req *TenantBackupRequest) (*model.DagDetailDTO, error) {
	response := c.createPostTenantBackupResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// TenantBackupSyncWithRequest synchronously executes the tenant backup request.
func (c *Client) TenantBackupSyncWithRequest(req *TenantBackupRequest) (*model.DagDetailDTO, error) {
	dag, err := c.TenantBackupWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}

type TenantBackupResponse struct {
	*response.TaskResponse
}

func (c *Client) createPostTenantBackupResponse() *TenantBackupResponse {
	return &TenantBackupResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

type TenantBackupStatusPatchParam struct {
	baseBackupStatusParam
}

type TenantBackupStatusPatchRequest struct {
	*request.BaseRequest
	param *TenantBackupStatusPatchParam
}

// NewTenantBackupStatusPatchRequest creates a PATCH request to update the tenant backup status.
func (c *Client) NewTenantBackupStatusPatchRequest(tenantName string) *TenantBackupStatusPatchRequest {
	req := &TenantBackupStatusPatchRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &TenantBackupStatusPatchParam{},
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/backup", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "PATCH")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *TenantBackupStatusPatchRequest) SetStatus(status string) *TenantBackupStatusPatchRequest {
	r.param.Status = &status
	r.SetBody(r.param)
	return r
}

// PatchTenantBackupStatus sends a PATCH request to update the tenant backup status.
func (c *Client) PatchTenantBackupStatus(tenantName string) error {
	req := c.NewTenantBackupStatusPatchRequest(tenantName)
	return c.TenantBackupStatusWithPatchRequest(req)
}

// TenantBackupStatusWithPatchRequest executes the PATCH request to update tenant backup status.
func (c *Client) TenantBackupStatusWithPatchRequest(req *TenantBackupStatusPatchRequest) error {
	response := c.createTenantBackupStatusResponse()
	return c.Execute(req, response)
}

type TenantBackupStatusResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createTenantBackupStatusResponse() *TenantBackupStatusResponse {
	return &TenantBackupStatusResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

type TenantLogStatusPatchRequest struct {
	*request.BaseRequest
	param *PatchLogStatusParam
}

// NewTenantLogStatusPatchRequest creates a PATCH request to update tenant log status.
func (c *Client) NewTenantLogStatusPatchRequest(tenantName string) *TenantLogStatusPatchRequest {
	req := &TenantLogStatusPatchRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &PatchLogStatusParam{},
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/backup/log", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "PATCH")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *TenantLogStatusPatchRequest) SetStatus(status string) *TenantLogStatusPatchRequest {
	r.param.Status = &status
	r.SetBody(r.param)
	return r
}

// PatchTenantLogStatus sends a PATCH request to update the tenant log status.
func (c *Client) PatchTenantLogStatus(tenantName string) error {
	req := c.NewTenantLogStatusPatchRequest(tenantName)
	return c.TenantLogStatusWithPatchRequest(req)
}

// TenantLogStatusWithPatchRequest executes the PATCH request to update tenant log status.
func (c *Client) TenantLogStatusWithPatchRequest(req *TenantLogStatusPatchRequest) error {
	response := c.createTenantLogStatusResponse()
	return c.Execute(req, response)
}

type TenantLogStatusResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createTenantLogStatusResponse() *TenantLogStatusResponse {
	return &TenantLogStatusResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

type TenantBackupOverviewRequest struct {
	*request.BaseRequest
}

// NewTenantBackupOverviewRequest creates a request to retrieve an overview of tenant backups.
func (c *Client) NewTenantBackupOverviewRequest(tenantName string) *TenantBackupOverviewRequest {
	req := &TenantBackupOverviewRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/backup/overview", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type TenantBackupOverview struct {
	Status model.CdbObBackupTask `json:"status"`
}

type TenantBackupOverviewResponse struct {
	*response.OcsAgentResponse
	TenantBackupOverview
}

func (c *Client) createTenantBackupOverviewResponse() *TenantBackupOverviewResponse {
	resp := &TenantBackupOverviewResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
	resp.Data = &resp.TenantBackupOverview
	return resp
}

// GetTenantBackupOverview fetches an overview of tenant backups.
func (c *Client) GetTenantBackupOverview(tenantName string) (*model.CdbObBackupTask, error) {
	req := c.NewTenantBackupOverviewRequest(tenantName)
	return c.GetTenantBackupOverviewWithRequest(req)
}

// GetTenantBackupOverviewWithRequest executes the request to fetch the tenant backup overview.
func (c *Client) GetTenantBackupOverviewWithRequest(req *TenantBackupOverviewRequest) (*model.CdbObBackupTask, error) {
	response := c.createTenantBackupOverviewResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return &response.TenantBackupOverview.Status, nil
}
