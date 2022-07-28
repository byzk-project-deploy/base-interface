package rpcinterfaces

type PluginType byte

func (p PluginType) Is(pluginType PluginType) bool {
	return p&pluginType == pluginType
}

const (
	// PluginTypeCmd 命令行插件
	PluginTypeCmd PluginType = 1 << iota
	// PluginTypeWeb web插件
	PluginTypeWeb
)

type OsOrArch byte

const (
	OsOrArchAll OsOrArch = 0
	OsLinux     OsOrArch = 1 << iota
	OsDarwin
	ArchAmd64
	ArchArm
	ArchArm64
	ArchMips64le
)

// PluginInfo 插件
type PluginInfo struct {
	// Name 插件名称
	Name string
	// Desc 插件描述
	Desc string
	// Type 插件类别
	Type PluginType
	// AllowOsAndArch 允许的系统或者架构
	AllowOsAndArch []OsOrArch
	// NotAllowOsAndArch 不被允许的平台或者架构
	NotAllowOsAndArch []OsOrArch
}

// DeployApplicationPluginInterface 部署应用的插件接口
type DeployApplicationPluginInterface interface {
	// Info 插件信息
	Info() (*PluginInfo, error)
}
