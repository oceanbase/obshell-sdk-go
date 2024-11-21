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

func GetInfo(server string) (*model.AgentRunStatus, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/info", server))
	if err != nil {
		log.Warn("Failed to get version: Network error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("Failed to get version: Read error: %v", err)
		return nil, err
	}

	var response struct {
		Data model.AgentRunStatus `json:"data"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		log.Warn("Failed to get version: Unmarshal error: %v", err)
		return nil, err
	}

	return &response.Data, nil
}

func GetIdentity(server string) (model.AgentIdentity, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/info", server))
	if err != nil {
		return model.UNIDENTIFIED, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.UNIDENTIFIED, err
	}

	var response struct {
		Data struct {
			Identity model.AgentIdentity `json:"identity"`
		} `json:"data"`
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return model.UNIDENTIFIED, err
	}
	return response.Data.Identity, nil
}
