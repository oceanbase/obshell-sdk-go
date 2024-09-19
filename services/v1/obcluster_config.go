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

	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/auth"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
	"github.com/oceanbase/obshell-sdk-go/util"
)

type ConfigObclusterRequest struct {
	*request.BaseRequest
	body     map[string]interface{}
	password string
}

type ConfigObclusterResponse struct {
	*response.TaskResponse
}

func (c *Client) createConfigObclusterResponse() *ConfigObclusterResponse {
	return &ConfigObclusterResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewConfigObclusterRequest return a ConfigObclusterRequest, which can be used as the argument for the ConfigObclusterWithRequest/ConfigObclusterSyncWithRequest.
// clusterName: the name of the cluster.
// clusterId: the id of the cluster.
// You can set the root password by SetRootPwd.
func (c *Client) NewConfigObclusterRequest(clusterName string, clusterId int) *ConfigObclusterRequest {
	req := &ConfigObclusterRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		body: map[string]interface{}{
			"clusterName": clusterName,
			"clusterId":   clusterId,
		},
	}
	req.SetBody(req.body)
	req.InitApiInfo("/api/v1/obcluster/config", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	return req
}

// SetRootPwd set the root password of the observer.
func (r *ConfigObclusterRequest) SetRootPwd(rootPwd string) *ConfigObclusterRequest {
	r.body["rootPwd"] = rootPwd
	r.password = rootPwd
	r.SetBody(r.body)
	return r
}

// ConfigObcluster returns a DagDetailDTO and an error, when the config obcluster task is completed successfully, the error will be nil.
// clusterName: the name of the cluster.
// clusterId: the id of the cluster.
// If you want to set the root password, you need to use NewConfigObclusterRequest and call SetRootPwd.
func (c *Client) ConfigObcluster(clusterName string, clusterId int) (*model.DagDetailDTO, error) {
	req := c.NewConfigObclusterRequest(clusterName, clusterId)
	return c.ConfigObclusterWithRequest(req)
}

func (r *ConfigObclusterRequest) encryptPassword() error {
	pwd, exist := r.body["rootPwd"]
	if exist {
		agentVserion, _, err := util.GetVersion(r.GetServer())
		if err != nil {
			return fmt.Errorf("get agent version error: %v", err)
		}
		if auth.VERSION_4_2_4.After(agentVserion) {
			pk, _ := util.GetPublicKey(r.GetServer())
			r.body["rootPwd"], err = auth.RSAEncrypt([]byte(pwd.(string)), pk)
			if err != nil {
				return fmt.Errorf("encrypt password error: %v", err)
			}
		}
	}
	return nil
}

// ConfigObclusterWithRequest returns a DagDetailDTO and an error, when the config obcluster task is requested successfully, the error will be nil.
// the parameter is a ConfigObclusterRequest, which can be created by NewConfigObclusterRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ConfigObclusterWithRequest(req *ConfigObclusterRequest) (dag *model.DagDetailDTO, err error) {
	c.setPasswordCandidateAuth(req.password)
	response := c.createConfigObclusterResponse()
	if err = req.encryptPassword(); err != nil {
		return
	}
	err = c.Execute(req, response)
	dag = response.DagDetailDTO
	go func() {
		defer func() {
			c.DiscardCandidateAuth()
		}()

		for {
			dag, err := c.GetDag(dag.GenericID)
			if err != nil {
				return
			}
			if dag.IsFinished() {
				c.AdoptCandidateAuth()
				return
			}
			time.Sleep(2 * time.Second)
		}
	}()
	return
}

// ConfigObclusterSyncWithRequest returns a DagDetailDTO and an error, when the config obcluster task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ConfigObclusterRequest, which can be created by NewConfigObclusterRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ConfigObclusterSyncWithRequest(request *ConfigObclusterRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ConfigObclusterWithRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
