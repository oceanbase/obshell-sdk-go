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

package services

import (
	"github.com/oceanbase/obshell-sdk-go/sdk/option"
	v1 "github.com/oceanbase/obshell-sdk-go/services/v1"
	v2 "github.com/oceanbase/obshell-sdk-go/services/v2"
)

type Clientset struct {
	clientV1 *v1.Client
	clientV2 *v2.Client
}

// NewClient creates a new client with the given host and port.
//
// The client will use the given options to configure the client.
// You can use the sdk.WithPasswordAuth option to set the password for the client.
//
// AS: sdk.NewClient("127.0.0.1", 8080, sdk.WithPasswordAuth("password"))
func NewClient(host string, port int, options ...option.Optioner) (*Clientset, error) {
	clientV1, err := v1.NewClient(host, port, options...)
	if err != nil {
		return nil, err
	}
	return &Clientset{
		clientV1: clientV1,
	}, nil
}

func NewClientWithServer(host string, port int) (*Clientset, error) {
	clientV1, err := v1.NewClientWithServer(host, port)
	if err != nil {
		return nil, err
	}
	return &Clientset{
		clientV1: clientV1,
	}, nil
}

func NewClientWithPassword(host string, port int, password string) (*Clientset, error) {
	clientV1, err := v1.NewClientWithPassword(host, port, password)
	if err != nil {
		return nil, err
	}
	return &Clientset{
		clientV1: clientV1,
	}, nil
}

func (c *Clientset) V1() *v1.Client {
	return c.clientV1
}

func (c *Clientset) V2() *v2.Client {
	return c.clientV2
}
