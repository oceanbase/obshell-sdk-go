package main

import (
	"fmt"

	"github.com/oceanbase/obshell-sdk-go/util"
)

func main() {
	ips := []string{
		"10.10.21.90",
		// "10.10.21.91",
		// "10.10.21.92",
	}

	workDir := "/data/ob" // OBServer 的工作目录，不需要提前创建，初始化 OBServer 时会自动创建
	nodeConfigs := make([]util.NodeConfig, 0)
	for _, ip := range ips {
		nodeConfigs = append(nodeConfigs, util.NewNodeConfig(ip, workDir, 2886))
	}
	errors, warns := util.CheckNodes(nodeConfigs...)
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
	}
	if len(warns) > 0 {
		for _, warn := range warns {
			fmt.Println(warn)
		}
	}
}