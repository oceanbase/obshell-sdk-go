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
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/oceanbase/obshell-sdk-go/sdk/request"
	"github.com/oceanbase/obshell-sdk-go/util"
)

type PasswordAuthV2 struct {
	*PasswordAuthMethod
}

func newPasswordAuthV2(pwd string) *PasswordAuthV2 {
	return &PasswordAuthV2{
		newPasswordAuthMethod(pwd),
	}
}

func (auth *PasswordAuthV2) Auth(req request.Request) error {
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

	originalBody := req.OriginalBody()
	if originalBody == nil {
		req.SetOriginalBody(req.Body())
		originalBody = req.OriginalBody()
	}
	encryptedBody, header, err := auth.BuildBodyAndHeader(originalBody, auth.pwd, req.GetUri())
	if err != nil {
		return err
	}
	req.SetHeaderByKey("X-OCS-Header", header["X-OCS-Header"])
	req.SetBody(encryptedBody)
	return nil
}

func (auth *PasswordAuthV2) BuildBodyAndHeader(param interface{}, pwd string, uri string) (encryptedBody interface{}, header map[string]string, err error) {
	encryptedBody, Key, Iv, err := EncryptBodyWithAes(param)
	if err != nil {
		return nil, nil, err
	}
	header, err = auth.BuildHeader(pwd, uri, Key, Iv)
	if err != nil {
		return nil, nil, err
	}
	return encryptedBody, header, nil
}

type HttpHeader struct {
	Auth string
	Ts   string
	Uri  string
	Keys []byte
}

func (auth *PasswordAuthV2) BuildHeader(pwd, uri string, keys ...[]byte) (map[string]string, error) {
	headers := make(map[string]string)

	var aesKeys []byte
	if len(keys) != 2 {
		aesKeys = nil
	} else {
		aesKeys = append(keys[0], keys[1]...)
	}
	header := HttpHeader{
		Auth: pwd,
		Ts:   fmt.Sprintf("%d", time.Now().Add(100*time.Second).Unix()),
		Uri:  uri,
		Keys: aesKeys,
	}
	mAuth, err := json.Marshal(header)
	if err != nil {
		return nil, errors.Wrap(err, "marshal header failed")
	}
	encrypt, err := RSAEncrypt(mAuth, auth.pk)
	if err != nil {
		return nil, errors.Wrap(err, "rsa encrypt failed")
	}
	headers["X-OCS-Header"] = encrypt
	return headers, nil
}

func RSAEncrypt(raw []byte, pk string) (string, error) {
	pkix, err := base64.StdEncoding.DecodeString(pk)
	if err != nil {
		return "", err
	}
	pub, err := x509.ParsePKCS1PublicKey(pkix)
	if err != nil {
		return "", errors.Wrap(err, "parse public key failed")
	}
	if len(raw) == 0 {
		b, err := rsa.EncryptPKCS1v15(rand.Reader, pub, raw)
		return base64.StdEncoding.EncodeToString(b), errors.Wrap(err, "encrypt failed")
	}
	// 分段加密
	blockSize := 512/8 - 11
	numBlocks := (len(raw) + blockSize - 1) / blockSize
	ciphertext := make([]byte, 0)
	for i := 0; i < numBlocks; i++ {
		start := i * blockSize
		end := (i + 1) * blockSize
		if end > len(raw) {
			end = len(raw)
		}
		encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, pub, raw[start:end])
		if err != nil {
			return "", errors.Wrap(err, "rsa encrypt failed")
		}
		ciphertext = append(ciphertext, encrypted...)
	}
	return base64.StdEncoding.EncodeToString(ciphertext), err
}

func EncryptBodyWithAes(body interface{}) (encryptedBody interface{}, key []byte, iv []byte, err error) {
	if body == nil {
		return
	}
	var mBody []byte
	if _, ok := body.([]byte); !ok {
		mBody, err = json.Marshal(body)
		if err != nil {
			return
		}
	} else {
		mBody = body.([]byte)
	}
	key = make([]byte, 16)
	iv = make([]byte, 16) // equal to block_size，16 bytes
	_, err = rand.Read(key)
	if err != nil {
		return
	}
	_, err = rand.Read(iv)
	if err != nil {
		return
	}
	encryptedBody, err = AESEncrypt(mBody, key, iv)
	return
}

func AESEncrypt(raw []byte, key []byte, iv []byte) (string, error) {
	// 创建AES加密器
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	raw = pkcs5Padding(raw, block.BlockSize())
	// 创建AES的CBC模式加密器
	mode := cipher.NewCBCEncrypter(block, iv)

	// 加密数据
	ciphertext := make([]byte, len(raw))
	mode.CryptBlocks(ciphertext, raw)

	return base64.StdEncoding.EncodeToString(ciphertext), err
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
