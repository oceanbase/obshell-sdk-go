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

package option

import (
	"crypto/tls"
	"crypto/x509"
)

type OptionType int

const (
	AUTH_OPT OptionType = iota + 1
	PROTOCOL_OPT
)

type Optioner interface {
	Type() OptionType
	Name() string
	Value() interface{}
}

type BaseOption struct {
	name       string
	optionType OptionType
	value      interface{}
}

func NewBaseOption(name string, optionType OptionType, value interface{}) *BaseOption {
	return &BaseOption{
		name:       name,
		optionType: optionType,
		value:      value,
	}
}

func (o *BaseOption) Type() OptionType {
	return o.optionType
}

func (o *BaseOption) Name() string {
	return o.name
}

func (o *BaseOption) Value() interface{} {
	return o.value
}

// ProtocolOption carries the HTTP protocol and TLS settings.
type ProtocolOption struct {
	protocol           string
	insecureSkipVerify bool
	certPool           *x509.CertPool
	clientCert         *tls.Certificate
}

func NewProtocolOption(protocol string, insecureSkipVerify bool) *ProtocolOption {
	return &ProtocolOption{
		protocol:           protocol,
		insecureSkipVerify: insecureSkipVerify,
	}
}

func NewProtocolOptionWithCertPool(protocol string, certPool *x509.CertPool) *ProtocolOption {
	return &ProtocolOption{
		protocol: protocol,
		certPool: certPool,
	}
}

func NewProtocolOptionWithClientCert(protocol string, certPool *x509.CertPool, clientCert tls.Certificate) *ProtocolOption {
	return &ProtocolOption{
		protocol:   protocol,
		certPool:   certPool,
		clientCert: &clientCert,
	}
}

func (o *ProtocolOption) Type() OptionType   { return PROTOCOL_OPT }
func (o *ProtocolOption) Name() string       { return "protocol" }
func (o *ProtocolOption) Value() interface{} { return o }

func (o *ProtocolOption) GetProtocol() string          { return o.protocol }
func (o *ProtocolOption) GetInsecureSkipVerify() bool  { return o.insecureSkipVerify }
func (o *ProtocolOption) GetCertPool() *x509.CertPool  { return o.certPool }
func (o *ProtocolOption) GetClientCert() *tls.Certificate { return o.clientCert }

// TLSConfig builds a *tls.Config from the option's TLS settings.
// Returns nil when using plain HTTP (no TLS fields set).
func (o *ProtocolOption) TLSConfig() *tls.Config {
	if !o.insecureSkipVerify && o.certPool == nil && o.clientCert == nil {
		return nil
	}
	cfg := &tls.Config{}
	if o.insecureSkipVerify {
		cfg.InsecureSkipVerify = true // #nosec G402
	}
	if o.certPool != nil {
		cfg.RootCAs = o.certPool
	}
	if o.clientCert != nil {
		cfg.Certificates = []tls.Certificate{*o.clientCert}
	}
	return cfg
}
