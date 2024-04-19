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
)

// getPublicKey function retrieves the public key from the API
func GetPublicKey(server string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/api/v1/secret", server))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var response struct {
		Data struct {
			PublicKey string `json:"public_key"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	return response.Data.PublicKey, nil
}
