package rpcinterfaces

// PluginInfoRPCServer 部署应用的插件RPC服务
type PluginInfoRPCServer struct {
	Impl PluginInfoInterface
}

// Info 方法具体实现
func (d *PluginInfoRPCServer) Info(args interface{}, resp **PluginInfo) (err error) {
	*resp, err = d.Impl.Info()
	return err
}
