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
	"github.com/oceanbase/obshell-sdk-go/log"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/util"
)

// PasswordAuth is a struct that implements Auther interface.
// It is used to authenticate with password.
type PasswordAuth struct {
	*BaseAuth
	pwd string
}

func NewPasswordAuth(pwd string) *PasswordAuth {
	return &PasswordAuth{
		BaseAuth: newBaseAuth(AUTH_TYPE_PASSWORD, AUTH_V1, AUTH_V2),
		pwd:      pwd,
	}
}

func (auth *PasswordAuth) Auth(request request.Request) error {
	method := auth.method
	if method != nil {
		return method.Auth(request)
	}

	switch auth.GetVersion() {
	case AUTH_V1:
		auth.method = newPasswordAuthV1(auth.pwd)
	case AUTH_V2:
		auth.method = newPasswordAuthV2(auth.pwd)
	default:
		return ErrNotSupportedAuthVersion
	}

	return auth.method.Auth(request)
}

type PasswordAuthMethod struct {
	pwd           string
	pk            string
	identityCheck bool
}

func (auth *PasswordAuthMethod) Reset() {
	auth.pk = ""
}

func (auth *PasswordAuthMethod) checkIdentity(req request.Request) error {
	if !auth.identityCheck {
		identity, err := util.GetIdentity(req.GetServer())
		if err != nil {
			return err
		}
		if identity == util.SINGLE {
			auth.pwd = ""
			log.Warn("Identity is single, password is not needed.")
		}
		auth.identityCheck = true
	}
	return nil
}
func newPasswordAuthMethod(pwd string) *PasswordAuthMethod {
	return &PasswordAuthMethod{
		pwd: pwd,
	}
}
