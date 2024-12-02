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

type ClusterBackupConfigRequest struct {
	*request.BaseRequest
	canPatchUri bool
	param       *ClusterBackupConfigApiParam
}

type ClusterBackupConfigApiParam struct {
	BackupBaseUri string `json:"backup_base_uri"`
	baseBackupConfigParam
}

// baseBackupConfigParam contains the basic configuration parameters for backup operations.
type baseBackupConfigParam struct {
	LogArchiveConcurrency *int          `json:"log_archive_concurrency"`
	Binding               *string       `json:"binding"`
	HaLowThreadScore      *int          `json:"ha_low_thread_score"`
	PieceSwitchInterval   *string       `json:"piece_switch_interval"`
	ArchiveLagTarget      *string       `json:"archive_lag_target"`
	DeletePolicy          *DeletePolicy `json:"delete_policy"`
}

// DeletePolicy defines the policy for deleting backup data.
type DeletePolicy struct {
	Policy         string `json:"policy"`
	RecoveryWindow string `json:"recovery_window"`
}

// NewClusterBackupConfigPostRequest creates a new POST request for cluster backup configuration.
func (c *Client) NewClusterBackupConfigPostRequest(backupBaseUri string) *ClusterBackupConfigRequest {
	return c.newClusterBackupConfigRequest(backupBaseUri, "POST")
}

// NewClusterBackupConfigPatchRequest creates a new PATCH request for modifying cluster backup configuration.
func (c *Client) NewClusterBackupConfigPatchRequest() *ClusterBackupConfigRequest {
	return c.newClusterBackupConfigRequest("", "PATCH")
}

// newClusterBackupConfigRequest initializes a new request for cluster backup configuration.
func (c *Client) newClusterBackupConfigRequest(backupBaseUri, method string) *ClusterBackupConfigRequest {
	req := &ClusterBackupConfigRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param: &ClusterBackupConfigApiParam{
			BackupBaseUri: backupBaseUri,
		},
	}
	if method == "PATCH" {
		req.canPatchUri = true
	}
	req.InitApiInfo("/api/v1/obcluster/backup/config", c.GetHost(), c.GetPort(), method)
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *ClusterBackupConfigRequest) SetBackupBaseUri(backupBaseUri string) *ClusterBackupConfigRequest {
	if r.canPatchUri {
		r.param.BackupBaseUri = backupBaseUri
		r.SetBody(r.param)
	}
	return r
}

func (r *ClusterBackupConfigRequest) SetLogArchiveConcurrency(logArchiveConcurrency int) *ClusterBackupConfigRequest {
	r.param.LogArchiveConcurrency = &logArchiveConcurrency
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupConfigRequest) SetBinding(binding string) *ClusterBackupConfigRequest {
	r.param.Binding = &binding
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupConfigRequest) SetHaLowThreadScore(haLowThreadScore int) *ClusterBackupConfigRequest {
	r.param.HaLowThreadScore = &haLowThreadScore
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupConfigRequest) SetPieceSwitchInterval(pieceSwitchInterval string) *ClusterBackupConfigRequest {
	r.param.PieceSwitchInterval = &pieceSwitchInterval
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupConfigRequest) SetArchiveLagTarget(archiveLagTarget string) *ClusterBackupConfigRequest {
	r.param.ArchiveLagTarget = &archiveLagTarget
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupConfigRequest) SetDeletePolicy(policy, recoveryWindow string) *ClusterBackupConfigRequest {
	r.param.DeletePolicy = &DeletePolicy{
		Policy:         policy,
		RecoveryWindow: recoveryWindow,
	}
	r.SetBody(r.param)
	return r
}

// PostClusterBackupConfig submits a POST request to configure cluster backup settings.
func (c *Client) PostClusterBackupConfig(dataBaseUri string) (*model.DagDetailDTO, error) {
	req := c.NewClusterBackupConfigPostRequest(dataBaseUri)
	return c.ClusterBackupConfigSyncWithRequest(req)
}

// ClusterBackupConfigWithRequest executes the cluster backup configuration request.
func (c *Client) ClusterBackupConfigWithRequest(req *ClusterBackupConfigRequest) (*model.DagDetailDTO, error) {
	response := c.createClusterBackupConfigResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ClusterBackupConfigSyncWithRequest synchronously executes the cluster backup configuration request.
func (c *Client) ClusterBackupConfigSyncWithRequest(req *ClusterBackupConfigRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ClusterBackupConfigWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}

type ClusterBackupConfigResponse struct {
	*response.TaskResponse
}

func (c *Client) createClusterBackupConfigResponse() *ClusterBackupConfigResponse {
	return &ClusterBackupConfigResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

type ClusterBackupRequest struct {
	*request.BaseRequest
	param *ClusterBackupApiParam
}

type baseBackupApiParam struct {
	Mode        *string `json:"mode"`
	PlusArchive *bool   `json:"plus_archive"`
	Encryption  *string `json:"encryption"`
}

type ClusterBackupApiParam struct {
	baseBackupApiParam
}

// NewClusterBackupRequest creates a new request for initiating a cluster backup.
func (c *Client) NewClusterBackupRequest() *ClusterBackupRequest {
	req := &ClusterBackupRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &ClusterBackupApiParam{},
	}
	req.InitApiInfo("/api/v1/obcluster/backup", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *ClusterBackupRequest) SetBackupMode(mode string) *ClusterBackupRequest {
	r.param.Mode = &mode
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupRequest) SetPlusArchive(plusArchive bool) *ClusterBackupRequest {
	r.param.PlusArchive = &plusArchive
	r.SetBody(r.param)
	return r
}

func (r *ClusterBackupRequest) SetEncryption(Encryption string) *ClusterBackupRequest {
	r.param.Encryption = &Encryption
	r.SetBody(r.param)
	return r
}

func (c *Client) PostClusterBackup() (*model.DagDetailDTO, error) {
	req := c.NewClusterBackupRequest()
	return c.ClusterBackupSyncWithRequest(req)
}

// ClusterBackupWithRequest executes the cluster backup request.
func (c *Client) ClusterBackupWithRequest(req *ClusterBackupRequest) (*model.DagDetailDTO, error) {
	response := c.createClusterBackupResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.DagDetailDTO, nil
}

// ClusterBackupSyncWithRequest synchronously executes the cluster backup request.
func (c *Client) ClusterBackupSyncWithRequest(req *ClusterBackupRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ClusterBackupWithRequest(req)
	if err != nil {
		return nil, err
	}
	return c.WaitDagSucceed(dag.GenericID)
}

type ClusterBackupResponse struct {
	*response.TaskResponse
}

func (c *Client) createClusterBackupResponse() *ClusterBackupResponse {
	return &ClusterBackupResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

type baseBackupStatusParam struct {
	Status *string `json:"status"`
}

type ClusterBackupStatusPatchParam struct {
	baseBackupStatusParam
}

type ClusterBackupStatusPatchRequest struct {
	*request.BaseRequest
	param *ClusterBackupStatusPatchParam
}

// NewClusterBackupStatusPatchRequest creates a new PATCH request to update backup status.
func (c *Client) NewClusterBackupStatusPatchRequest() *ClusterBackupStatusPatchRequest {
	req := &ClusterBackupStatusPatchRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &ClusterBackupStatusPatchParam{},
	}
	req.InitApiInfo("/api/v1/obcluster/backup", c.GetHost(), c.GetPort(), "PATCH")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *ClusterBackupStatusPatchRequest) SetStatus(status string) *ClusterBackupStatusPatchRequest {
	r.param.Status = &status
	r.SetBody(r.param)
	return r
}

func (c *Client) PatchClusterBackupStatus() error {
	req := c.NewClusterBackupStatusPatchRequest()
	return c.ClusterBackupStatusWithPatchRequest(req)
}

// ClusterBackupStatusWithPatchRequest synchronously executes the PATCH request to update backup status.
func (c *Client) ClusterBackupStatusWithPatchRequest(req *ClusterBackupStatusPatchRequest) error {
	response := c.createClusterBackupStatusResponse()
	return c.Execute(req, response)
}

type ClusterBackupStatusResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createClusterBackupStatusResponse() *ClusterBackupStatusResponse {
	return &ClusterBackupStatusResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

type PatchLogStatusParam struct {
	Status *string `json:"status"`
}

type ClusterLogStatusPatchRequest struct {
	*request.BaseRequest
	param *PatchLogStatusParam
}

// NewClusterLogStatusPatchRequest creates a new PATCH request to update log status.
func (c *Client) NewClusterLogStatusPatchRequest() *ClusterLogStatusPatchRequest {
	req := &ClusterLogStatusPatchRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		param:       &PatchLogStatusParam{},
	}
	req.InitApiInfo("/api/v1/obcluster/backup/log", c.GetHost(), c.GetPort(), "PATCH")
	req.SetAuthentication()
	req.SetBody(req.param)
	return req
}

func (r *ClusterLogStatusPatchRequest) SetStatus(status string) *ClusterLogStatusPatchRequest {
	r.param.Status = &status
	r.SetBody(r.param)
	return r
}

func (c *Client) PatchClusterLogStatus() error {
	req := c.NewClusterLogStatusPatchRequest()
	return c.ClusterLogStatusWithPatchRequest(req)
}

// NewClusterLogStatusPatchRequest creates a new PATCH request to update log status.
func (c *Client) ClusterLogStatusWithPatchRequest(req *ClusterLogStatusPatchRequest) error {
	response := c.createClusterLogStatusResponse()
	return c.Execute(req, response)
}

type ClusterLogStatusResponse struct {
	*response.OcsAgentResponse
}

func (c *Client) createClusterLogStatusResponse() *ClusterLogStatusResponse {
	return &ClusterLogStatusResponse{
		OcsAgentResponse: response.NewOcsAgentResponseWithoutReturn(),
	}
}

type ClusterBackupOverviewRequest struct {
	*request.BaseRequest
}

// NewClusterBackupOverviewRequest creates a new request to fetch the cluster backup overview.
func (c *Client) NewClusterBackupOverviewRequest() *ClusterBackupOverviewRequest {
	req := &ClusterBackupOverviewRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	req.InitApiInfo("/api/v1/obcluster/backup/overview", c.GetHost(), c.GetPort(), "GET")
	req.SetAuthentication()
	return req
}

type ClusterBackupOverview struct {
	Statuses []model.CdbObBackupTask `json:"statuses"`
}

type ClusterBackupOverviewResponse struct {
	*response.OcsAgentResponse
	ClusterBackupOverview
}

func (c *Client) createClusterBackupOverviewResponse() *ClusterBackupOverviewResponse {
	resp := &ClusterBackupOverviewResponse{
		OcsAgentResponse: response.NewOcsAgentResponse(),
	}
	resp.Data = &resp.ClusterBackupOverview
	return resp
}

// GetClusterBackupOverview fetches the overview of cluster backups.
func (c *Client) GetClusterBackupOverview() ([]model.CdbObBackupTask, error) {
	req := c.NewClusterBackupOverviewRequest()
	return c.GetClusterBackupOverviewWithRequest(req)
}

// GetClusterBackupOverviewWithRequest executes the request to fetch the backup overview.
func (c *Client) GetClusterBackupOverviewWithRequest(req *ClusterBackupOverviewRequest) ([]model.CdbObBackupTask, error) {
	response := c.createClusterBackupOverviewResponse()
	if err := c.Execute(req, response); err != nil {
		return nil, err
	}
	return response.Statuses, nil
}
