package rpcinterfaces

import "net/rpc"

// DeployApplicationPluginRPC 部署应用的插件RPC客户端
type DeployApplicationPluginRPC struct {
	client *rpc.Client
}

func (d *DeployApplicationPluginRPC) Info() *PluginInfo {
	var resp *PluginInfo
	err := d.client.Call("Plugin.Info", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}
