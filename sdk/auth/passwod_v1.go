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
	"encoding/json"
	"time"

	"github.com/oceanbase/obshell-sdk-go/internal/util"
	"github.com/oceanbase/obshell-sdk-go/sdk/request"
)

type PasswordAuthV1 struct {
	*PasswordAuthMethod
}

func newPasswordAuthV1(pwd string, letftime time.Duration) *PasswordAuthV1 {
	return &PasswordAuthV1{
		PasswordAuthMethod: newPasswordAuthMethod(pwd, letftime),
	}
}

func (auth *PasswordAuthV1) Auth(req request.Request, context *request.Context) error {
	if !req.Authentication() {
		return nil
	}

	var err error
	if err = auth.checkIdentity(req); err != nil {
		return err
	}

	if auth.pk == "" {
		auth.pk, err = util.GetPublicKey(req.GetServer())
		if err != nil {
			return err
		}
	}

	authMap := map[string]interface{}{
		"password": auth.pwd,
		"ts":       time.Now().Unix() + int64(auth.letftime),
	}
	authJSON, err := json.Marshal(authMap)
	if err != nil {
		return err
	}
	encryptedPwd, err := RSAEncrypt(authJSON, auth.pk)
	if err != nil {
		return err
	}
	context.SetHeader("X-OCS-Auth", encryptedPwd)

	return nil
}
