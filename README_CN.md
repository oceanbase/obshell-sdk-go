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

[英文版](README.md) | 中文版

**OBShell-SDK-GO** 是 [OceanBase 社区](https://open.oceanbase.com/) 为了方便开发者快速使用 OBShell 服务而提供的 SDK，开发者可以使用该 SDK 便捷地调用 OBShell 的接口。

## 安装
```shell
go get github.com/obshell/obshell-sdk-go
```

## 快速使用
使用时请确保OBShell处于运行状态
### 创建客户端
您可以选择创建单一版本客户端。
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
或者创建多版本客户端集。
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
### 部署OBShell集群
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
### 发起扩容
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
