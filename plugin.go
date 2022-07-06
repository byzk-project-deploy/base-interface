package rpcinterfaces

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

// DeployApplicationPlugin 部署应用插件
type DeployApplicationPlugin struct {
	// Impl 插件具体实现
	Impl DeployApplicationPluginInterface
}

func (d *DeployApplicationPlugin) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &DeployApplicationPluginRPCServer{Impl: d.Impl}, nil
}

func (d *DeployApplicationPlugin) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &DeployApplicationPluginRPC{client: client}, nil
}
