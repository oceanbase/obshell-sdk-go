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

package model

import "fmt"

type AgentInfo struct {
	Ip   string `json:"ip" form:"ip" binding:"required"`
	Port int    `json:"port" form:"port" binding:"required"`
}

func (agent *AgentInfo) String() string {
	return fmt.Sprintf("%s:%d", agent.Ip, agent.Port)
}

type ObInfoResp struct {
	Agents []AgentInstance `json:"agent_info"`
	Config ClusterConfig   `json:"obcluster_info"`
}

type AgentRunStatus struct {
	Pid      int    `json:"pid"`
	State    int32  `json:"state"`
	StartAt  int64  `json:"startAt"`
	HomePath string `json:"homePath"`
	AgentInstance
	SupportedAuth []string `json:"supportedAuth"`
}

type AgentInstance struct {
	AgentInfo
	Zone     string `json:"zone" binding:"required"`
	Identity string `json:"identity" binding:"required"`
	Version  string `json:"version"`
}

type AgentStatus struct {
	Agent AgentInfoWithIdentity `json:"agent"`
	// service state
	State int32 `json:"state"`
	// service version
	Version string `json:"version"`
	// service pid
	Pid int `json:"pid"`
	// timestamp when service started
	StartAt int64 `json:"startAt"`
	// Ports process occupied ports
	OBState          int  `json:"obState"`
	UnderMaintenance bool `json:"underMaintenance"`
}

type AgentInfoWithIdentity struct {
	AgentInfo
	Identity string `json:"identity" binding:"required"`
}

type ClusterConfig struct {
	ClusterID   int                        `json:"id"`
	ClusterName string                     `json:"name"`
	Version     string                     `json:"version"`
	ZoneConfig  map[string][]*ServerConfig `json:"topology"`
}

type ServerConfig struct {
	SvrIP        string `json:"svr_ip"`
	SvrPort      int    `json:"svr_port"`
	SqlPort      int    `json:"sql_port"`
	AgentPort    int    `json:"agent_port"`
	WithRootSvr  string `json:"with_rootserver"`
	Status       string `json:"status"`
	BuildVersion string `json:"build_version"`
}

type Scope struct {
	Type   string   `json:"type"`
	Target []string `json:"target"`
}

type GitInfo struct {
	GitBranch        string `json:"branch"`
	GitCommitId      string `json:"commitId"`
	GitShortCommitId string `json:"shortCommitId"`
	GitCommitTime    string `json:"commitTime"`
}
