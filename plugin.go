package rpcinterfaces

import (
	"github.com/byzk-project-deploy/go-plugin"
	"github.com/hashicorp/go-hclog"
	"github.com/tjfoc/gmsm/gmtls"
	"github.com/tjfoc/gmsm/x509"
	"os"
)

const (
	PluginNameBase = "BYPT_PLUGIN_BASE"
	PluginNameCmd  = "BYPT_PLUGIN_CMD"
	PluginNameWeb  = "BYPT_PLUGIN_WEB"
)

const pluginTestRootCert = `-----BEGIN CERTIFICATE-----
MIIFwDCCA6igAwIBAgIIO86GiXsG/2AwDQYJKoZIhvcNAQELBQAwZTELMAkGA1UE
BhMCemgxEDAOBgNVBAgTB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxDTALBgNV
BAoTBGJ5emsxDTALBgNVBAsTBGJ5cHQxFDASBgNVBAMTC3BsdWdpbiB0ZXN0MCAX
DTIyMDgwNDA2MjEwMFoYDzI5OTkwODI4MDYyMTAwWjBlMQswCQYDVQQGEwJ6aDEQ
MA4GA1UECBMHQmVpSmluZzEQMA4GA1UEBxMHQmVpSmluZzENMAsGA1UEChMEYnl6
azENMAsGA1UECxMEYnlwdDEUMBIGA1UEAxMLcGx1Z2luIHRlc3QwggIiMA0GCSqG
SIb3DQEBAQUAA4ICDwAwggIKAoICAQC9wiIAmS+ioNcm+ryK0IfWZ0kdr+eURYJ2
VXb1DqcJMPHZHMAQG9DCPCMWfVShO537WNWGXNJI1IwHSZBaEs8rlHStwnUkfwJM
gnkHFIFRz8vBUV39o6Dmnkkxnklc4kLMgZX3ACYuhpl08UTopw72WkUfvUXvIZvw
f4MG9Gjt+n8qtkKKWk3nLbSBgxo1UPp084bT6dBcNw84BhSg/r6rA9r0OrygcMJz
jt0ecGZJDUffF47OLPhhwySZN3gwE9okD0Npu9fI3muyUYYhJZfLw8XUdIzpfmcr
taXMvqA8ioiOpVAVqRplvY05hIOuB0E9ZDM8PCMSisxkRvx5GVTK3E5xUPeafGBP
4ALZBBCzG1ny19KaJqF9uVkDYcILRpZKztKYn1TxfId1GP/igBV7pBnM2jB5NO4V
T+1hFdPyYRbVQTBVn+ecgjXjlEbsQlpNvef0+MvHibTNXvykF0FKqx5FWDVXh7e0
QqyeFjtrtpli3xP2EAF/cRP5DX9GtvmZWJ8nYrp1Ybnc9E8tz900EdUDhQdeURQY
jLZJ1bs2KAnSDkHkhlXFO3m5hnKARvnYakxa92uWMpOPD9WbykKNRNReNFyYQZYp
W0Bjequ2GN6BO3FO0MRS4ETuwDNp1DnAhclHizmYJq2s8a+2KNdB6AC59PN8KEBj
U/3hrpIWIQIDAQABo3IwcDAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTBe3gm
eQQY+YoTXrheyV43gB1KqzALBgNVHQ8EBAMCAQYwEQYJYIZIAYb4QgEBBAQDAgAH
MB4GCWCGSAGG+EIBDQQRFg94Y2EgY2VydGlmaWNhdGUwDQYJKoZIhvcNAQELBQAD
ggIBACAEq9pg+oXJ9Si7AqnUJwk6DOHZRbrKm8KFbn6/Wr8qTY7kQtZqmzg1TivB
yaHQUrd7orRKsuKVwd9zyLYDKi1/5npRedrJrXFTM0lI8xZ+R0m4wfhW1YnagNnh
JS+MuhCF2VHCp9gr8n0A6ROSYvefz+Qcj0HoemlDavdVCNgdVoYS2XkWd+VnjAhf
3XrudnUVxoEz0waBZkywOEcVUSq6yZvn2lHgg8EuZlS2e8D+jshgztVQskz+7Tkn
Fh5kQ7US2DIyf2ftzIgJ6xUzWiitiapk2lxi/AHJx+b1TXB90+relV4seB1p2AKO
vdB2fktphxoiOFQtJaybg5FzW1Sc9D+Fsxpxlol73RkF9x2DEGBftzfCXXnN44CJ
CqRjSl9JObf6Hzx1J/n7HpKeH1VNYhjUjAgvJBH3tvS3WUUiUj4+6VChQg0VNU8O
HMGNlCu7Qcvry0MSo7AmJVKa83xG7heDkfnIK4Z4pCNAW4ntvOxF6qQLdfaF9Wuu
huB+kKbhGRnkOjNDZ4vmCF4v132D98NiQQ2fGeyfil3FHWrz54NzchEU1jFUN0Nb
XkKJwFn6BtrU/XecBV93Uy/2p/iKxIgVPJQEmZkRHqkPSMWiEy8180RnaFGj9mzf
YIBJ89tjKkm5BvIwBfLxfmrzf3vHzuPURicU8h0WGzei5ueY
-----END CERTIFICATE-----`

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
	TLSProvider func() (*gmtls.Config, error)
	// CertPem 客户端证书PEM
	CertPem string
	// CertPrivateKeyPem 客户端证书私钥PEM
	CertPrivateKeyPem string
}

type PluginServeCallback func(logger hclog.Logger, rootCertPem string) *PluginServeCallbackResult

type emptyPluginCmdImpl struct {
}

func (e emptyPluginCmdImpl) Commands() []*CmdInfo {
	return nil
}

// PluginServe 插件监听
func PluginServe(fn PluginServeCallback) {
	//unixFile := os.Getenv("UNIX_FILE")

	rootCert := os.Getenv("ROOT_CERT_PEM")
	if rootCert == "" {
		rootCert = pluginTestRootCert
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	res := fn(logger, rootCert)
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
			pluginMap[PluginNameCmd] = &PluginCmdImpl{impl: &emptyPluginCmdImpl{}}
		} else {
			pluginMap[PluginNameCmd] = &PluginCmdImpl{impl: res.CmdPlugin}
		}
	}

	if pluginInfo.Type.Is(PluginTypeWeb) {
		//TODO 注册WEB插件接口
	}

	if res.TLSProvider == nil && res.CertPem != "" && res.CertPrivateKeyPem != "" {
		pair, err := gmtls.LoadX509KeyPair(res.CertPem, res.CertPrivateKeyPem)
		if err != nil {
			logger.Error("插件证书解析失败: %s", err.Error())
			os.Exit(3)
		}

		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM([]byte(rootCert))

		res.TLSProvider = func() (*gmtls.Config, error) {
			return &gmtls.Config{
				Certificates:       []gmtls.Certificate{pair},
				ClientAuth:         gmtls.RequireAndVerifyClientCert,
				ServerName:         pluginInfo.Name,
				InsecureSkipVerify: false,
				RootCAs:            pool,
			}, nil
		}
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: *res.HandshakeConfig,
		Plugins:         pluginMap,
		TLSProvider:     res.TLSProvider,
	})
}
