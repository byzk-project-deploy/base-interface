package rpcinterfaces

import (
	"github.com/byzk-project-deploy/go-plugin"
	"net/rpc"
	"sync"
	"time"
)

type PluginType byte

var (
	blockingMap = make(map[string]chan struct{})
	blockerLock = &sync.Mutex{}
)

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

type PluginBaseImpl struct {
	impl PluginBaseInterface
}

func (p *PluginBaseImpl) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &pluginBaseRpcServer{impl: p.impl}, nil
}

func (p *PluginBaseImpl) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &pluginBaseRpc{client: client}, nil
}

type pluginBaseRpc struct {
	client *rpc.Client
}

func (p *pluginBaseRpc) Start() error {
	return p.client.Call("Plugin.Start", new(any), new(any))
}

func (p *pluginBaseRpc) Stop() {
	_ = p.client.Call("Plugin.Stop", new(any), new(any))
}

func (p *pluginBaseRpc) Info() (resp *PluginInfo, err error) {
	return resp, withTimeout(5*time.Second, func() error {
		return p.client.Call("Plugin.Info", new(any), &resp)
	})
}

func (p *pluginBaseRpc) Blocking(id string) error {
	return p.client.Call("Plugin.Blocking", id, new(any))
}

func (p *pluginBaseRpc) Revoke(id string) error {
	return p.client.Call("Plugin.Revoke", id, new(any))
}

type pluginBaseRpcServer struct {
	impl PluginBaseInterface
}

func (p pluginBaseRpcServer) Info(args any, resp **PluginInfo) (err error) {
	*resp, err = p.impl.Info()
	return
}

func (p pluginBaseRpcServer) Start(args any, resp *any) error {
	return p.impl.Start()
}

func (p pluginBaseRpcServer) Stop(args any, resp *any) error {
	p.impl.Stop()
	return nil
}

func (p pluginBaseRpcServer) Blocking(args string, resp *any) error {
	blockerLock.Lock()
	ch, ok := blockingMap[args]
	if !ok {
		ch = make(chan struct{}, 1)
		blockingMap[args] = ch
	}
	blockerLock.Unlock()
	<-ch
	return nil
}

func (p pluginBaseRpcServer) Revoke(args string, resp *any) error {
	blockerLock.Lock()
	defer blockerLock.Unlock()
	if ch, ok := blockingMap[args]; !ok {
		return nil
	} else {
		close(ch)
		delete(blockingMap, args)
	}
	return nil
}
