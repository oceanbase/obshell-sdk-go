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
	"time"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
)

type RestoreApiParam struct {
	DataBackupUri string  `json:"data_backup_uri" binding:"required"`
	ArchiveLogUri *string `json:"archive_log_uri"`

	TenantName string `json:"restore_tenant_name" binding:"required"`

	Timestamp *time.Time `json:"timestamp" time_format:"2006-01-02T15:04:05.000Z07:00"`
	SCN       *int64     `json:"scn"`

	ZoneList             []ZoneParam `json:"zone_list" binding:"required"` // Tenant zone list with unit config.
	PrimaryZone          *string     `json:"primary_zone"`
	Concurrency          *int        `json:"concurrency"`
	HaHighThreadScore    *int        `json:"ha_high_thread_score"`
	FullBackupDecryption *string     `json:"full_backup_decryption"`
	IncBackupDecryption  *string     `json:"inc_backup_decryption"`
	KmsEncryptInfo       *string     `json:"kms_encrypt_info"`
}

type RestoreRequest struct {
	*request.BaseRequest
	param *RestoreApiParam
}

// NewRestoreRequest creates a new restore request with the specified parameters.
func (c *Client) NewRestoreRequest(tenantName, dataBackupUri string, zoneList []ZoneParam) *RestoreRequest {
	req := &RestoreRequest{
		BaseRequest: request.NewBaseRequest(),
		param: &RestoreApiParam{
			DataBackupUri: dataBackupUri,
			TenantName:    tenantName,
			ZoneList:      zoneList,
		},
	}
	req.InitApiInfo("/api/v1/tenant/restore", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *RestoreRequest) SetArchiveLogUri(archiveLogUri string) *RestoreRequest {
	r.param.ArchiveLogUri = &archiveLogUri
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetTimestamp(timestamp time.Time) *RestoreRequest {
	r.param.Timestamp = &timestamp
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetSCN(scn int64) *RestoreRequest {
	r.param.SCN = &scn
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetPrimaryZone(primaryZone string) *RestoreRequest {
	r.param.PrimaryZone = &primaryZone
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetConcurrency(concurrency int) *RestoreRequest {
	r.param.Concurrency = &concurrency
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetHaHighThreadScore(haHighThreadScore int) *RestoreRequest {
	r.param.HaHighThreadScore = &haHighThreadScore
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetFullBackupDecryption(fullBackupDecryption string) *RestoreRequest {
	r.param.FullBackupDecryption = &fullBackupDecryption
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetIncBackupDecryption(incBackupDecryption string) *RestoreRequest {
	r.param.IncBackupDecryption = &incBackupDecryption
	r.SetBody(r.param)
	return r
}

func (r *RestoreRequest) SetKmsEncryptInfo(kmsEncryptInfo string) *RestoreRequest {
	r.param.KmsEncryptInfo = &kmsEncryptInfo
	r.SetBody(r.param)
	return r
}

// Restore initiates a restore operation with the given parameters.
func (c *Client) Restore(dataBackupUri, tenantName string, zoneList []ZoneParam) (*model.DagDetailDTO, error) {
	req := c.NewRestoreRequest(dataBackupUri, tenantName, zoneList)
	return c.RestoreSyncWithRequest(req)
}

// RestoreWithRequest executes the restore operation with the specified request.
func (c *Client) RestoreWithRequest(req *RestoreRequest) (*model.DagDetailDTO, error) {
	response := c.createRestoreResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil

}

// RestoreSyncWithRequest executes the restore operation synchronously.
func (c *Client) RestoreSyncWithRequest(req *RestoreRequest) (*model.DagDetailDTO, error) {
	dag, err := c.RestoreWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}

func (c *Client) createRestoreResponse() *RestoreResponse {
	return &RestoreResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

type RestoreResponse struct {
	*response.TaskResponse
}

type TenantRestoreOverviewRequest struct {
	*request.BaseRequest
}

// NewTenantRestoreOverviewRequest creates a request to retrieve the overview of a tenant restore process.
func (c *Client) NewTenantRestoreOverviewRequest(tenantName string) *TenantRestoreOverviewRequest {
	req := &TenantRestoreOverviewRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/restore/overview", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type TenantRestoreOverviewResponse struct {
	*response.OcsAgentResponse
	model.RestoreOverview
}

func (c *Client) createTenantRestoreOverviewResponse() *TenantRestoreOverviewResponse {
	resp := &TenantRestoreOverviewResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
	resp.Data = &resp.RestoreOverview
	return resp
}

// GetTenantRestoreOverview retrieves an overview of the tenant restore process.
func (c *Client) GetTenantRestoreOverview(tenantName string) (*model.RestoreOverview, error) {
	req := c.NewTenantRestoreOverviewRequest(tenantName)
	return c.GetTenantRestoreOverviewWithRequest(req)
}

// GetTenantRestoreOverviewWithRequest executes the request to get the tenant restore overview.
func (c *Client) GetTenantRestoreOverviewWithRequest(req *TenantRestoreOverviewRequest) (*model.RestoreOverview, error) {
	response := c.createTenantRestoreOverviewResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return &response.RestoreOverview, nil
}

type CancelRestoreRequest struct {
	*request.BaseRequest
}

// NewCancelRestoreRequest creates a request to cancel a restore operation.
func (c *Client) NewCancelRestoreRequest(tenantName string) *CancelRestoreRequest {
	req := &CancelRestoreRequest{
		BaseRequest: request.NewBaseRequest(),
	}
	uri := fmt.Sprintf("/api/v1/tenant/%s/restore", tenantName)
	req.InitApiInfo(uri, c.GetHost(), c.GetPort(), "DELETE")
	req.SetAuthentication()
	return req
}

type cancelRestoreResponse struct {
	*response.TaskResponse
}

func (c *Client) createCancelRestoreResponse() *cancelRestoreResponse {
	return &cancelRestoreResponse{
		response.NewTaskResponse(),
	}
}

// CancelRestore cancels a restore operation.
func (c *Client) CancelRestore(tenantName string) (dag *model.DagDetailDTO, err error) {
	req := c.NewCancelRestoreRequest(tenantName)
	return c.CancelRestoreWithRequest(req)
}

// CancelRestoreWithRequest cancels a restore operation with the specified request.
func (c *Client) CancelRestoreWithRequest(req *CancelRestoreRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createCancelRestoreResponse()
	if err = c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, err
}

// CancelRestoreSyncWithRequest cancels a restore operation synchronously.
func (c *Client) CancelRestoreSyncWithRequest(req *CancelRestoreRequest) (dag *model.DagDetailDTO, err error) {
	dag, err = c.CancelRestoreWithRequest(req)
	if err != nil {
		return nil, err
	}
	if dag == nil || dag.GenericDTO == nil {
		return nil, nil
	}
	return c.WaitDagSucceed(dag.GenericID)
}
