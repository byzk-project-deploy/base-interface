package rpcinterfaces

import "net/rpc"

// DeployApplicationPluginRPC 部署应用的插件RPC客户端
type DeployApplicationPluginRPC struct {
	client *rpc.Client
}

func (d *DeployApplicationPluginRPC) Info() *PluginInfo {
	var resp *PluginInfo
	err := d.client.Call("Plugin.Info", nil, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}
