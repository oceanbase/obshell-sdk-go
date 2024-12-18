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

	"github.com/oceanbase/obshell-sdk-go/sdk/option"
)

type AuthOption struct {
	*option.BaseOption
}

func newAuthOption(value Auther) AuthOption {
	return AuthOption{
		BaseOption: option.NewBaseOption("auth", option.AUTH_OPT, value),
	}
}

type PasswordAuthOption struct {
	AuthOption
	password *PasswordAuth
}

func (auth *PasswordAuthOption) SetLifetime(lifetime time.Duration) *PasswordAuthOption {
	auth.password.SetLifetime(lifetime)
	return auth
}

func (auth *PasswordAuthOption) GetLifetime() time.Duration {
	return auth.password.GetLifetime()
}

func newPasswordAuthOption(value *PasswordAuth) *PasswordAuthOption {
	return &PasswordAuthOption{
		AuthOption: newAuthOption(value),
		password:   value,
	}
}

func WithPasswordAuth(pwd string) *PasswordAuthOption {
	return newPasswordAuthOption(NewPasswordAuth(pwd))
}
