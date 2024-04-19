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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/obshell-sdk-go/model"
)

func isValidIp(ip string) bool {
	ipRegexp := regexp.MustCompile(`^(?:(?:25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(?:25[0-5]|2[0-4]\d|[01]?\d\d?)$`)
	return ipRegexp.MatchString(ip)
}

func isValidPort(port string) bool {
	if port == "" {
		return true
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return false
	}
	return p > 1024 && p < 65536
}

func ParseAddr(addr string) (*model.AgentInfo, error) {
	args := strings.Split(addr, ":")
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid addr")
	}
	if !isValidIp(args[0]) {
		return nil, fmt.Errorf("invalid ip: %s", args[0])
	}
	if !isValidPort(args[1]) {
		return nil, fmt.Errorf("invalid port: %s", args[1])
	}
	port, _ := strconv.Atoi(args[1])
	return &model.AgentInfo{Ip: args[0], Port: port}, nil
}
