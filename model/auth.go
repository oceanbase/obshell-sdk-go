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

type ObshellAuthInfo struct {
	AuthVersion   string   `json:"auth_version"`
	SupportedAuth []string `json:"supported_auth"`
}

func (authInfo *ObshellAuthInfo) IsSupported(version string) bool {
	for _, v := range authInfo.SupportedAuth {
		if v == version {
			return true
		}
	}
	return false
}
