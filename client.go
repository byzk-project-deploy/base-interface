package rpcinterfaces

import "net/rpc"

// DeployApplicationPluginRPC 部署应用的插件RPC客户端
type DeployApplicationPluginRPC struct {
	client *rpc.Client
}

func (d *DeployApplicationPluginRPC) Info() (*PluginInfo, error) {
	var resp *PluginInfo
	err := d.client.Call("Plugin.Info", new(interface{}), &resp)
	if err != nil {
		return nil, d.errHandle(err)
	}

	return resp, nil
}

func (d *DeployApplicationPluginRPC) errHandle(err error) error {
	if _, ok := err.(rpc.ServerError); ok {
		return err
	}
	panic(err)
}
