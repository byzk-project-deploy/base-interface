package main

import (
	rpcinterfaces "github.com/byzk-project-deploy/base-interface"
	"github.com/hashicorp/go-hclog"
)

func main() {
	rpcinterfaces.TestWithInteractive(func(logger hclog.Logger, rootCert string) *rpcinterfaces.PluginServeCallbackResult {
		return &rpcinterfaces.PluginServeCallbackResult{
			BasePlugin: &testPluginBaseImpl{
				logger: logger,
			},
			CmdPlugin: &testPluginCmdImpl{
				logger: logger,
			},
		}
	})
}
