package rpcinterfaces

import (
	"crypto/tls"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"net/rpc"
	"os"
)

const PluginName = "BYPT_PLUGIN"

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

type PluginServeCallback func(logger hclog.Logger) (plugin.Plugin, *plugin.HandshakeConfig, *tls.Config)

// PluginServe 插件监听
func PluginServe(fn PluginServeCallback) {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	p, config, tlsConfig := fn(logger)
	if p == nil || config == nil {
		logger.Error("缺失的参数配置")
		os.Exit(1)
	}

	pluginMap := map[string]plugin.Plugin{
		"PluginName": p,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: *config,
		Plugins:         pluginMap,
		TLSProvider: func() (*tls.Config, error) {
			return tlsConfig, nil
		},
	})
}