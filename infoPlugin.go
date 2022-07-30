package rpcinterfaces

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
	"time"
)

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

func (o OsOrArch) Is(os OsOrArch, arch OsOrArch) bool {
	return o&os == os && o&arch == arch
}

const (
	OsOrArchAll OsOrArch = 0
	// OsLinux Linux系统
	OsLinux OsOrArch = 1 << iota
	// OsDarwin Mac系统
	OsDarwin
	// ArchAmd64 amd64系统架构
	ArchAmd64
	// ArchArm arm系统架构
	ArchArm
	// ArchArm64 arm64系统架构
	ArchArm64
	// ArchMips64le Mips64le系统架构
	ArchMips64le
)

// PluginInfo 插件
type PluginInfo struct {
	// Author 作者名称
	Author string
	// Name 插件名称
	Name string
	// ShortDesc 短描述，一般大于30字将被裁剪
	ShortDesc string
	// Desc 插件描述（支持Markdown）
	Desc string
	// CreateTime 创建时间
	CreateTime time.Time
	// Type 插件类别
	Type PluginType
	// AllowOsAndArch 允许的系统或者架构
	AllowOsAndArch []OsOrArch
	// NotAllowOsAndArch 不被允许的平台或者架构
	NotAllowOsAndArch []OsOrArch
}

// PluginInfoInterface 插件信息接口
type PluginInfoInterface interface {
	// Info 插件信息
	Info() (*PluginInfo, error)
}

type PluginInfoImpl struct {
	impl PluginInfoInterface
}

func (p *PluginInfoImpl) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &pluginInfoRpcServer{impl: p.impl}, nil
}

func (p *PluginInfoImpl) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &pluginInfoRpc{client: client}, nil
}

type pluginInfoRpc struct {
	client *rpc.Client
}

func (p *pluginInfoRpc) Info() (resp *PluginInfo, err error) {
	return resp, withTimeout(5*time.Second, func() error {
		return p.client.Call("Plugin.Info", new(interface{}), &resp)
	})
}

type pluginInfoRpcServer struct {
	impl PluginInfoInterface
}

func (p pluginInfoRpcServer) Info(args interface{}, resp **PluginInfo) (err error) {
	*resp, err = p.impl.Info()
	return
}
