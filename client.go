package rpcinterfaces

import "net/rpc"

// DeployApplicationPluginRPC 部署应用的插件RPC客户端
type DeployApplicationPluginRPC struct {
	client *rpc.Client
}
