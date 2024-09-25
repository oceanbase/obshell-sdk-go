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

package auth

import (
	"errors"
	"strconv"
	"strings"

	"github.com/oceanbase/obshell-sdk-go/sdk/request"
)

// auth version
const (
	AUTH_V1 = "v1"
	AUTH_V2 = "v2"
)

// obshell version
var (
	VERSION_4_2_2 = newAuthVersion("4.2.2.0")
	VERSION_4_2_3 = newAuthVersion("4.2.3.0")
	VERSION_4_2_4 = newAuthVersion("4.2.4.0")
)

func init() {
	VERSION_4_2_2.SetVersion("4.2.2.0")
	VERSION_4_2_3.SetVersion("4.2.3.0")
	VERSION_4_2_4.SetVersion("4.2.4.0")
}

// auth type
const (
	AUTH_TYPE_PASSWORD AuthType = iota + 1
)

type AuthType int

// error type
var (
	ErrNotSupportedAuthVersion = errors.New("unsupported auth version")
)

type AuthMethod interface {
	Auth(request request.Request, context *request.Context) error
	Reset()
}

type Versioner interface {
	GetCompatibleList() []string
	IsSupported(version string) bool
	SetVersion(version string) bool
	GetVersion() string
	AutoSelectVersion(version ...string) bool
	IsAutoSelectVersion() bool
	Berfore(version string) bool
	After(version string) bool
	Equals(version string) bool
	BeforeOrEquals(version string) bool
	AfterOrEquals(version string) bool
}

type Auther interface {
	Versioner
	Type() AuthType
	Reset()       // Reset will set method to nil.
	ResetMethod() // ResetMethod will call method.Reset()
	Auth(request request.Request, context *request.Context) error
}

// AuthVersion implements Versioner
type AuthVersion struct {
	version           string
	isAutoSelect      bool
	compatibleVersion []string
}

func newAuthVersion(version ...string) *AuthVersion {
	return &AuthVersion{
		compatibleVersion: version,
		isAutoSelect:      true,
	}
}

func (authVersion *AuthVersion) GetCompatibleList() []string {
	return authVersion.compatibleVersion
}

func (authVersion *AuthVersion) IsSupported(version string) bool {
	version = strings.Split(version, "-")[0]
	for _, v := range authVersion.compatibleVersion {
		if v == version {
			return true
		}
	}
	return false
}

func (authVersion *AuthVersion) SetVersion(version string) bool {
	version = strings.Split(version, "-")[0]
	if !authVersion.IsSupported(version) {
		return false
	}
	authVersion.version = version
	authVersion.isAutoSelect = false
	return true
}

func (authVersion *AuthVersion) GetVersion() string {
	return authVersion.version
}

func (authVersion *AuthVersion) AutoSelectVersion(version ...string) bool {
	for _, v := range version {
		if authVersion.IsSupported(v) {
			authVersion.version = strings.Split(v, "-")[0]
			authVersion.isAutoSelect = true
			return true
		}
	}
	return false
}

func (authVersion *AuthVersion) IsAutoSelectVersion() bool {
	return authVersion.isAutoSelect
}

func (authVersion *AuthVersion) Berfore(version string) bool {
	return cmpVersionString(authVersion.version, strings.Split(version, "-")[0]) < 0
}

func (authVersion *AuthVersion) After(version string) bool {
	return cmpVersionString(authVersion.version, strings.Split(version, "-")[0]) > 0
}

func (authVersion *AuthVersion) Equals(version string) bool {
	return cmpVersionString(authVersion.version, strings.Split(version, "-")[0]) == 0
}

func (authVersion *AuthVersion) BeforeOrEquals(version string) bool {
	return cmpVersionString(authVersion.version, strings.Split(version, "-")[0]) <= 0
}

func (authVersion *AuthVersion) AfterOrEquals(version string) bool {
	return cmpVersionString(authVersion.version, strings.Split(version, "-")[0]) >= 0
}

// cmpVersionString compare two version string, return 1 if ver1 > ver2, -1 if ver1 < ver2, 0 if ver1 == ver2
func cmpVersionString(ver1, ver2 string) int {
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

// BaseAuth implements Versioner and Auther.GetMethod
type BaseAuth struct {
	Versioner
	method   AuthMethod
	authType AuthType
}

func newBaseAuth(authType AuthType, version ...string) *BaseAuth {
	return &BaseAuth{
		authType:  authType,
		Versioner: newAuthVersion(version...),
	}
}

func (auther *BaseAuth) ResetMethod() {
	if auther.method != nil {
		auther.method.Reset()
	}
}

func (auther *BaseAuth) Reset() {
	auther.method = nil
}

func (auther *BaseAuth) Type() AuthType {
	return auther.authType
}
