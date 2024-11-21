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
	"strconv"
	"strings"
)

// cmpVersionString compare two version string, return 1 if ver1 > ver2, -1 if ver1 < ver2, 0 if ver1 == ver2
func CmpVersionString(ver1, ver2 string) int {
	if ver1 == ver2 {
		return 0
	}
	ver1Arr := strings.Split(ver1, ".")
	ver2Arr := strings.Split(ver2, ".")
	for i := 0; i < len(ver1Arr) && i < len(ver2Arr); i++ {
		v1, _ := strconv.Atoi(ver1Arr[i])
		v2, _ := strconv.Atoi(ver2Arr[i])
		if v1 > v2 {
			return 1
		} else if v1 < v2 {
			return -1
		}
	}
	if len(ver1Arr) > len(ver2Arr) {
		return 1
	} else {
		return -1
	}
}
