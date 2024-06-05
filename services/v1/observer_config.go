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

const (
	SCOPE_GLOBAL = "GLOBAL"
	SCOPE_ZONE   = "ZONE"
	SCOPE_SERVER = "SERVER"
)

type ConfigObserverRequest struct {
	*request.BaseRequest
}

type obServerConfigParams struct {
	ObServerConfig map[string]string `json:"observerConfig" binding:"required"`
	Restart        bool              `json:"restart"`
	Scope          model.Scope       `json:"scope" binding:"required"`
}

type ConfigObserverResponse struct {
	*response.TaskResponse
}

func (c *Client) createConfigObserverResponse() *ConfigObserverResponse {
	return &ConfigObserverResponse{
		TaskResponse: response.NewTaskResponse(),
	}
}

// NewConfigObserverRequest return a ConfigObserverRequest, which can be used as the argument for the ConfigObserverWithRequest/ConfigObserverSyncWithRequest.
// configs: the configs of the observer.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL.
// targets is the target to be started, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
func (c *Client) NewConfigObserverRequest(configs map[string]string, level string, targets ...string) *ConfigObserverRequest {
	req := &ConfigObserverRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
	}
	obServerConfigParams := &obServerConfigParams{
		ObServerConfig: configs,
		Restart:        true,
		Scope: model.Scope{
			Type:   level,
			Target: targets,
		},
	}
	req.SetBody(obServerConfigParams)
	req.InitApiInfo("/api/v1/observer/config", c.GetHost(), c.GetPort(), "POST")
	req.SetAuthentication()
	return req
}

// ConfigObserver returns a DagDetailDTO and an error, when the config observer task is completed successfully, the error will be nil.
// configs: the configs of the observer.
// level: the level of the scope of the task, can be v1.SCOPE_SERVER, v1.SCOPE_ZONE, v1.SCOPE_GLOBAL.
// targets is the target to be started, can be zone name or server 'ip:port', when level is SCOPE_GLOBAL, targets is not needed.
func (c *Client) ConfigObserver(configs map[string]string, level string, targets ...string) (dag *model.DagDetailDTO, err error) {
	request := c.NewConfigObserverRequest(configs, level, targets...)
	return c.ConfigObserverWithRequest(request)
}

// ConfigObserverWithRequest returns a DagDetailDTO and an error, when the config observer task is requested successfully, the error will be nil.
// the parameter is a ConfigObserverRequest, which can be created by NewConfigObserverRequest.
// You can use WaitDagSucceed to wait for the task to complete.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ConfigObserverWithRequest(request *ConfigObserverRequest) (dag *model.DagDetailDTO, err error) {
	response := c.createConfigObserverResponse()
	err = c.Execute(request, response)
	dag = response.DagDetailDTO
	return
}

// ConfigObserverSyncWithRequest returns a DagDetailDTO and an error, when the config observer task is completed successfully, the error will be nil.
// the DagDetailDTO is the final status of the task.
// the parameter is a ConfigObserverRequest, which can be created by NewConfigObserverRequest.
// You can check or operater the task through the DagDetailDTO.
func (c *Client) ConfigObserverSyncWithRequest(request *ConfigObserverRequest) (*model.DagDetailDTO, error) {
	dag, err := c.ConfigObserverWithRequest(request)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	return c.WaitDagSucceed(dag.GenericID)
}
