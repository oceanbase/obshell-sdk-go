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
	"errors"
	"fmt"

	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/sdk/response"
	"github.com/oceanbase/obshell-sdk-go/util"
)

type CreateClusterRequest struct {
	targetHost string
	targetPort int
	password   string
	server     map[string]string // server -> zone
	requests   []request.Request // only support ConfigClusterRequest and ConfigObserverRequest
}

// AddServer add a server to the cluster later.
func (req *CreateClusterRequest) AddServer(host string, port int, zone string) *CreateClusterRequest {
	if req.server == nil {
		return nil
	}
	req.server[fmt.Sprintf("%s:%d", host, port)] = zone
	return req
}

// ConfigObserver can configure observer at different level (SCOPE_SERVER, SCOPE_ZONE and SCOPE_GLOBAL).
func (req *CreateClusterRequest) ConfigObserver(configs map[string]string, level string, targets ...string) *CreateClusterRequest {
	if req.requests == nil {
		return nil
	}
	configClusterRequest := &ConfigObserverRequest{
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
	configClusterRequest.SetBody(obServerConfigParams)
	configClusterRequest.InitApiInfo("/api/v1/observer/config", req.targetHost, req.targetPort, "POST")
	configClusterRequest.SetAuthentication()
	req.requests = append(req.requests, configClusterRequest)
	return req
}

// ConfigCluster can configure obcluster's configuration, password is not required (optional).
func (req *CreateClusterRequest) ConfigCluster(clusterName string, clusterId int, password ...string) *CreateClusterRequest {
	if req.requests == nil {
		return nil
	}
	configObclusterRequest := &ConfigObclusterRequest{
		BaseRequest: request.NewAsyncBaseRequest(),
		body: map[string]interface{}{
			"clusterName": clusterName,
			"clusterId":   clusterId,
		},
	}
	configObclusterRequest.InitApiInfo("/api/v1/obcluster/config", req.targetHost, req.targetPort, "POST")
	configObclusterRequest.SetAuthentication()
	req.requests = append(req.requests, configObclusterRequest)
	return req
}

// SetPassword set the password of the cluster.
func (req *CreateClusterRequest) SetPassword(pwd string) *CreateClusterRequest {
	req.password = pwd
	return req
}

// CreateCreateClusterRequest create a CreateClusterRequest, which can be used to create a cluster.
func (c *Client) NewCreateClusterRequest() *CreateClusterRequest {
	req := &CreateClusterRequest{
		targetHost: c.GetHost(),
		targetPort: c.GetPort(),
		server:     make(map[string]string),
		requests:   make([]request.Request, 0),
	}
	return req
}

func (c *Client) join(server, zone string) error {
	agentInfo, err := util.ParseAddr(server)
	if err != nil {
		return errors.New("The format of server only can be 'ip:port' at present ")
	}
	if _, err := c.Join(agentInfo.Ip, agentInfo.Port, zone); err != nil {
		return err
	}
	return nil
}

// CreateClusterWithRequest recieve a CreateClusterRequest, and send mutilple requests to OBShell to create a cluster.
// CreateClusterWithRequest is a synchronous method, it will return an error if any task is failed.
func (c *Client) CreateClusterWithRequest(req *CreateClusterRequest) (err error) {
	if len(req.server) == 0 {
		return fmt.Errorf("There is no servers to be joined")
	}
	// find the master
	masterZone, ok := req.server[c.GetServer()]
	if !ok {
		return fmt.Errorf("The master server is not in the server list")
	}
	// join master
	if err := c.join(c.GetServer(), masterZone); err != nil {
		return err
	}
	delete(req.server, c.GetServer())

	// join follower
	for server, zone := range req.server {
		if err := c.join(server, zone); err != nil {
			return err
		}
	}

	response := response.NewTaskResponse()
	for _, subReq := range req.requests {
		if configObclusterReq, ok := subReq.(*ConfigObclusterRequest); ok {
			configObclusterReq.SetRootPwd(req.password)
			c.setPasswordCandidateAuth(req.password)
			if err = configObclusterReq.encryptPassword(); err != nil {
				return err
			}
		}
		err := c.Execute(subReq, response)
		if err != nil {
			return err
		}
		dag := response.DagDetailDTO

		if _, err := c.WaitDagSucceed(dag.GenericID); err != nil {
			return err
		}
	}
	// init
	InitRequest := c.NewInitRequest()
	if _, err := c.InitSyncWithRequest(InitRequest); err != nil {
		return err
	}
	return nil
}
