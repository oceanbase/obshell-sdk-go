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

package sdk

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/oceanbase/obshell-sdk-go/sdk/auth"
	"github.com/oceanbase/obshell-sdk-go/sdk/option"
)

func WithPasswordAuth(pwd string) *auth.PasswordAuthOption {
	return auth.WithPasswordAuth(pwd)
}

// WithHttps configures HTTPS behavior via parameters.
//
// - insecureSkipVerify: when true, skips verification of the server certificate.
// - certPool: when non-nil, verifies the server certificate against this CA pool.
// - clientCert: when non-nil, sends the client certificate (mTLS).
func WithHttps(insecureSkipVerify bool, certPool *x509.CertPool, clientCert *tls.Certificate) *option.ProtocolOption {
	if clientCert != nil {
		if insecureSkipVerify && certPool == nil {
			return option.NewProtocolOptionWithClientCertInsecure("https", *clientCert)
		}
		return option.NewProtocolOptionWithClientCert("https", certPool, *clientCert)
	}
	if certPool != nil {
		return option.NewProtocolOptionWithCertPool("https", certPool)
	}
	return option.NewProtocolOption("https", insecureSkipVerify)
}
