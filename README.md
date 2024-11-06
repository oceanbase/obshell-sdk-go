<p align="center">
  <a href="https://github.com/oceanbase/oceanbase/blob/master/LICENSE">
    <img alt="license" src="https://img.shields.io/badge/license-Apache--2.0-blue" />
  </a>
  <a href="https://en.oceanbase.com/docs/oceanbase-database">
    <img alt="English doc" src="https://img.shields.io/badge/docs-English-blue" />
  </a>
  <a href="https://www.oceanbase.com/docs/oceanbase-database-cn">
    <img alt="Chinese doc" src="https://img.shields.io/badge/文档-简体中文-blue" />
  </a>
</p>

English | [Chinese](README_CN.md)


**OBShell-SDK-GO** is an SDK provided by the[OceanBase Community](https://open.oceanbase.com/) to facilitate developers with quick access to OBShell services, allowing them to conveniently call OBShell interfaces using this SDK.

## Install
```shell
go get github.com/oceanbase/obshell-sdk-go@master
```

## Quick Start
Please ensure that OBShell is running when using it.
### Create a Client
You can choose to create a single-version client.
``` GO
package main

import (
	"github.com/oceanbase/obshell-sdk-go/services/v1"
)

func main() {
	client, err := v1.NewClientWithPassword("11.11.11.1", 2886, "***")
	if err != nil {
        // Handle error.
		return
	}
}
```
Or create a multi-version client set.
``` GO
package main

import (
	"github.com/oceanbase/obshell-sdk-go/services"
)

func main() {
	clientset, err := services.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
        // Handle error.
		return
	}
}
```
### Deploy Cluster
``` GO
package main

import (
	"github.com/oceanbase/obshell-sdk-go/services"
	"github.com/oceanbase/obshell-sdk-go/services/v1"
)

func main() {
	client, err := services.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
		return
	}
	joinReqeust1 := client.V1().NewJoinRequest("11.11.11.1", 2886, "zone1")
	if _, err := client.V1().JoinSyncWithRequest(joinReqeust1); err != nil {
		return
	}

	joinReqeust2 := client.V1().NewJoinRequest("11.11.11.2", 2886, "zone2")
	if _, err := client.V1().JoinSyncWithRequest(joinReqeust2); err != nil {
		return
	}

	joinReqeust3 := client.V1().NewJoinRequest("11.11.11.3", 2886, "zone3")
	if _, err := client.V1().JoinSyncWithRequest(joinReqeust3); err != nil {
		return
	}

	// Configure the cluster.
	configObclusterReq := client.V1().NewConfigObclusterRequest("obshell-sdk-test", 12358).SetRootPwd("****")
	if _, err := client.V1().ConfigObclusterSyncWithRequest(configObclusterReq); err != nil {
		return
	}

	// Configure the observers.
	configs := map[string]string{
		"datafile_size": "24G", "cpu_count": "16", "memory_limit": "16G", "system_memory": "8G", "log_disk_size": "24G",
	}
	configObserverReq := client.V1().NewConfigObserverRequest(configs, v1.SCOPE_GLOBAL)
	if _, err := client.V1().ConfigObserverSyncWithRequest(configObserverReq); err != nil {
		return
	}

	// Initialize the cluster.
	initReq := client.V1().NewInitRequest()
	if _, err := client.V1().InitSyncWithRequest(initReq); err != nil {
		return
	}
}
```
### Scale out
``` GO
package main

import (
	"github.com/oceanbase/obshell-sdk-go/services"
)

func main() {
	client, err := services.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
		return
	}

	// OBShell prior to 4.2.3.0 should use mysqlPort(rpcPort) instead of mysql_port(rpc_port).
	configs := map[string]string{
		"mysql_port": "2881", "rpc_port": "2882", "datafile_size": "24G", "cpu_count": "16", "memory_limit": "16G", "system_memory": "8G", "log_disk_size": "24G",
	}

	// Scale out a new server in zone3.
	scaleOutReq := client.V1().NewScaleOutRequest("11.11.11.3", 2886, "zone3", configs)
	if _, err := client.V1().ScaleOutSyncWithRequest(scaleOutReq); err != nil {
		return
	}
}
```
