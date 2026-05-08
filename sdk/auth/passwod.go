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
	"time"

	"github.com/oceanbase/obshell-sdk-go/internal/util"
	"github.com/oceanbase/obshell-sdk-go/log"
	"github.com/oceanbase/obshell-sdk-go/model"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
)

// PasswordAuth is a struct that implements Auther interface.
// It is used to authenticate with password.
type PasswordAuth struct {
	*BaseAuth
	pwd                string
	letftime           time.Duration
	prefetchedInfo     bool
	prefetchedIdentity model.AgentIdentity
}

func NewPasswordAuth(pwd string) *PasswordAuth {
	return &PasswordAuth{
		BaseAuth: newBaseAuth(AUTH_TYPE_PASSWORD, AUTH_V1, AUTH_V2),
		pwd:      pwd,
	}
}

func (auth *PasswordAuth) SetLifetime(lifetime time.Duration) {
	auth.letftime = lifetime
	auth.method = nil
}

func (auth *PasswordAuth) GetLifetime() time.Duration {
	if auth.letftime == 0 {
		return 60 * time.Second
	}
	return auth.letftime
}

func (auth *PasswordAuth) SetAgentIdentity(identity model.AgentIdentity) {
	auth.prefetchedIdentity = identity
	auth.prefetchedInfo = true
}

func (auth *PasswordAuth) Auth(request request.Request, context *request.Context) error {
	method := auth.method
	if method != nil {
		return method.Auth(request, context)
	}

	switch auth.GetVersion() {
	case AUTH_V1:
		auth.method = newPasswordAuthV1(auth.pwd, auth.GetLifetime(), auth.prefetchedInfo, auth.prefetchedIdentity)
	case AUTH_V2:
		auth.method = newPasswordAuthV2(auth.pwd, auth.GetLifetime(), auth.prefetchedInfo, auth.prefetchedIdentity)
	default:
		return ErrNotSupportedAuthVersion
	}

	return auth.method.Auth(request, context)
}

type PasswordAuthMethod struct {
	pwd                string
	pk                 string
	identityCheck      bool
	letftime           time.Duration
	prefetchedIdentity model.AgentIdentity
}

func (auth *PasswordAuthMethod) Reset() {
	auth.pk = ""
	auth.identityCheck = false
	auth.prefetchedIdentity = model.UNIDENTIFIED
}

func (auth *PasswordAuthMethod) checkIdentity(req request.Request) error {
	if !auth.identityCheck {
		identity := auth.prefetchedIdentity
		if identity == model.UNIDENTIFIED {
			var err error
			identity, err = util.GetIdentityWithOptions(req.GetServer(), req.GetProtocol(), req.GetTLSConfig())
			if err != nil {
				return err
			}
		}
		if identity == model.SINGLE {
			auth.pwd = ""
			log.Warn("Identity is single, password is not needed.")
		}
		auth.identityCheck = true
		auth.prefetchedIdentity = model.UNIDENTIFIED
	}
	return nil
}

func newPasswordAuthMethod(pwd string, letftime time.Duration, prefetched bool, identity model.AgentIdentity) *PasswordAuthMethod {
	if !prefetched {
		identity = model.UNIDENTIFIED
	}
	return &PasswordAuthMethod{
		pwd:                pwd,
		letftime:           letftime,
		prefetchedIdentity: identity,
	}
}
