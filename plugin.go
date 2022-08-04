package rpcinterfaces

import (
	"crypto/tls"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

const (
	PluginNameBase = "BYPT_PLUGIN_BASE"
	PluginNameCmd  = "BYPT_PLUGIN_CMD"
	PluginNameWeb  = "BYPT_PLUGIN_WEB"
)

var DefaultHandshakeConfig = &plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BYPT",
	MagicCookieValue: "BYPT_VERIFY",
}

// PluginServeCallbackResult 插件监听回调结果
type PluginServeCallbackResult struct {
	// BasePlugin 信息插件( 必传 )
	BasePlugin PluginBaseInterface
	// CmdPlugin 命令行插件( 当信息插件内的插件类型包含cmd时生效 )
	CmdPlugin PluginCmdInterface
	// WebPlugin Web插件( 当信息插件内的插件类型包含web时生效 )
	WebPlugin PluginCmdInterface
	// HandshakeConfig 握手协议配置
	HandshakeConfig *plugin.HandshakeConfig
	// TLSProvider tls认证
	TLSProvider func() (*tls.Config, error)
}

type PluginServeCallback func(logger hclog.Logger) *PluginServeCallbackResult

// PluginServe 插件监听
func PluginServe(fn PluginServeCallback) {

	//unixFile := os.Getenv("UNIX_FILE")

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	res := fn(logger)
	if res.BasePlugin == nil {
		logger.Error("缺失插件信息内容")
		os.Exit(1)
	}

	pluginInfo, err := res.BasePlugin.Info()
	if err != nil {
		logger.Error("获取插件中的插件信息失败: %s", err.Error())
		os.Exit(2)
	}

	if res.HandshakeConfig == nil {
		res.HandshakeConfig = DefaultHandshakeConfig
	}

	pluginMap := map[string]plugin.Plugin{
		PluginNameBase: &PluginBaseImpl{impl: res.BasePlugin},
	}

	if pluginInfo.Type.Is(PluginTypeCmd) {
		if res.CmdPlugin == nil {
			logger.Error("插件缺失接口内容")
			os.Exit(2)
		}
		pluginMap[PluginNameCmd] = &PluginCmdImpl{impl: res.CmdPlugin}
	}

	if pluginInfo.Type.Is(PluginTypeWeb) {
		//TODO 注册WEB插件接口
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: *res.HandshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     res.TLSProvider,
	})
}
