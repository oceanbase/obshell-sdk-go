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

// WithHttps uses HTTPS with certificate verification enabled.
func WithHttps() *option.ProtocolOption {
	return option.NewProtocolOption("https", false)
}

// WithHttpsInsecure uses HTTPS and skips TLS certificate verification.
func WithHttpsInsecure() *option.ProtocolOption {
	return option.NewProtocolOption("https", true)
}

// WithHttpsCA uses HTTPS and verifies the server certificate against the provided CA pool.
func WithHttpsCA(certPool *x509.CertPool) *option.ProtocolOption {
	return option.NewProtocolOptionWithCertPool("https", certPool)
}

// WithHttpsClientCert uses HTTPS with a client certificate for mutual TLS (mTLS).
// certPool may be nil to use the system CA pool for server certificate verification.
func WithHttpsClientCert(clientCert tls.Certificate, certPool *x509.CertPool) *option.ProtocolOption {
	return option.NewProtocolOptionWithClientCert("https", certPool, clientCert)
}
