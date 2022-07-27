package rpcinterfaces

// PluginInfo 插件
type PluginInfo struct {
	// Name 插件名称
	Name string
	// Desc 插件描述
	Desc string
	// Type 插件类别
	Type string
}

// DeployApplicationPluginInterface 部署应用的插件接口
type DeployApplicationPluginInterface interface {
	// Info 插件信息
	Info() *PluginInfo
}
