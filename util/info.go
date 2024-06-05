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

package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/oceanbase/obshell-sdk-go/log"
	"github.com/oceanbase/obshell-sdk-go/model"
)

// getPublicKey function retrieves the public key from the API
func GetVersion(server string) (string, *model.ObshellAuthInfo, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/info", server))
	if err != nil {
		log.Warn("Failed to get version: Network error: %v", err)
		return "", nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Failed to get version: Read error: %v", err)
		return "", nil, err
	}

	var response struct {
		Data struct {
			Version  string                 `json:"version"`
			AuthInfo *model.ObshellAuthInfo `json:"auth_info"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		log.Warn("Failed to get version: Unmarshal error: %v", err)
		return "", nil, err
	}
	return response.Data.Version, response.Data.AuthInfo, nil
}

type AgentIdentity string

const (
	MASTER             AgentIdentity = "MASTER"
	FOLLOWER           AgentIdentity = "FOLLOWER"
	SINGLE             AgentIdentity = "SINGLE"
	CLUSTER_AGENT      AgentIdentity = "CLUSTER AGENT"
	TAKE_OVER_MASTER   AgentIdentity = "TAKE OVER MASTER"
	TAKE_OVER_FOLLOWER AgentIdentity = "TAKE OVER FOLLOWER"
	SCALING_OUT        AgentIdentity = "SCALING OUT"
	UNIDENTIFIED       AgentIdentity = "UNIDENTIFIED"
)

func GetIdentity(server string) (AgentIdentity, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/info", server))
	if err != nil {
		return UNIDENTIFIED, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UNIDENTIFIED, err
	}

	var response struct {
		Data struct {
			Identity AgentIdentity `json:"identity"`
		} `json:"data"`
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return UNIDENTIFIED, err
	}
	return response.Data.Identity, nil
}
