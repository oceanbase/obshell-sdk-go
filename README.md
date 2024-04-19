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
go get github.com/oceanbase/obshell-sdk-go
```

## Quick Start
Please ensure that OBShell is running when using it.
### Creating a Client
You can choose to create a single-version client.
``` GO
package main

import (
	"github.com/obshell-sdk-go/services/v1"
)

func main() {
	client, err := v1.NewClientWithPassword("11.11.11.1", 2886, "***")
	if err != nil {
        // Handle exceptions
		panic(err)
	}
}
```
Or create a multi-version client set.
``` GO
package main

import (
	"github.com/obshell-sdk-go/services"
)

func main() {
	clientset, err := services.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
        // Handle exceptions
		panic(err)
	}
}
```
### Deploying Cluster
``` GO
package main

import (
    "fmt"
	"github.com/obshell-sdk-go/services/v1"
)

func main() {
	client, err := v1.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
		panic(err)
	}
	joinReqeust1 := client.NewJoinRequest("11.11.11.1", 2886, "zone1")
	if _, err := client.JoinSyncWithRequest(joinReqeust1); err != nil {
		panic(err)
	}

	joinReqeust2 := client.NewJoinRequest("11.111.11.2", 2886, "zone2")
	if _, err := client.JoinSyncWithRequest(joinReqeust2); err != nil {
		panic(err)
	}

	joinReqeust3 := client.NewJoinRequest("11.11.11.3", 2886, "zone3")
	if _, err := client.JoinSyncWithRequest(joinReqeust3); err != nil {
		panic(err)
	}

	// config server
	for _, ip := range []string{"11.11.11.1", "11.11.11.2", "11.11.11.3"} {
		// OBShell prior to 4.2.3.0 should use mysqlPort(rpcPort) instead of mysql_port(rpc_port).
		configs := map[string]string{
			"mysql_port": "2881", "rpc_port": "2882", "datafile_size": "24G", "cpu_count": "16", "memory_limit": "16G", "system_memory": "8G", "log_disk_size": "24G",
		}
		configObserverReq := client.NewConfigObserverRequest(configs, v1.SCOPE_SERVER, fmt.Sprintf("%s:2886", ip))
		if _, err := client.ConfigObserverSyncWithRequest(configObserverReq); err != nil {
			panic(err)
		}
	}

	// config obcluster
	configObclusterReq := client.NewConfigObclusterRequest("obshell-sdk-test", 12358).SetRootPwd("****")
	if _, err := client.ConfigObclusterSyncWithRequest(configObclusterReq); err != nil {
		panic(err)
	}

	// init
	initReq := client.NewInitRequest()
	if _, err := client.InitSyncWithRequest(initReq); err != nil {
		panic(err)
	}
}
```
### Scale out
``` GO
package main

import (
	"github.com/obshell-sdk-go/services/v1"
)

func main() {
	client, err := v1.NewClientWithPassword("11.11.11.1", 2886, "****")
	if err != nil {
		panic(err)
	}

	// OBShell prior to 4.2.3.0 should use mysqlPort(rpcPort) instead of mysql_port(rpc_port).
	configs := map[string]string{
		"mysql_port": "2881", "rpc_port": "2882", "datafile_size": "24G", "cpu_count": "16", "memory_limit": "16G", "system_memory": "8G", "log_disk_size": "24G",
	}

	// scale-out
	scaleOutReq := client.NewScaleOutRequest("11.11.11.3", 2886, "zone3", configs)
	if _, err := client.ScaleOutSyncWithRequest(scaleOutReq); err != nil {
		panic(err)
	}
}
```
