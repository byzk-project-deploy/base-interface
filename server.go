package rpcinterfaces

// DeployApplicationPluginRPCServer 部署应用的插件RPC服务
type DeployApplicationPluginRPCServer struct {
	Impl DeployApplicationPluginInterface
}

// Info 方法具体实现
func (d *DeployApplicationPluginRPCServer) Info(args interface{}, resp **PluginInfo) (err error) {
	*resp, err = d.Impl.Info()
	return err
}
