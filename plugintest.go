package rpcinterfaces

import (
	"errors"
	"fmt"
	"github.com/byzk-project-deploy/go-plugin"
	"github.com/byzk-project-deploy/grumble"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"github.com/hashicorp/go-hclog"
	"os"
	"time"
)

func loadPlugin(res *PluginServeCallbackResult) (PluginBaseInterface, PluginCmdWrapperInterface, *plugin.Client) {
	pluginMap := map[string]plugin.Plugin{
		PluginNameBase: &PluginBaseImpl{impl: res.BasePlugin},
		PluginNameCmd:  &PluginCmdImpl{impl: res.CmdPlugin},
	}

	res.HandshakeConfig = DefaultHandshakeConfig

	reattachConfigCh := make(chan *plugin.ReattachConfig, 1)
	go func() {
		plugin.Serve(&plugin.ServeConfig{
			HandshakeConfig: *res.HandshakeConfig,
			Plugins:         pluginMap,
			Test: &plugin.ServeTestConfig{
				ReattachConfigCh: reattachConfigCh,
			},
		})
	}()

	reattachConfig := <-reattachConfigCh

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: *res.HandshakeConfig,
		Plugins:         pluginMap,
		Reattach:        reattachConfig,
	})

	rpcClient, err := client.Client()
	if err != nil {
		panic(err.Error())
	}

	basePlugin, err := rpcClient.Dispense(PluginNameBase)
	if err != nil {
		panic(err)
	}

	cmdPlugin, err := rpcClient.Dispense(PluginNameCmd)
	if err != nil {
		panic(err)
	}

	return basePlugin.(PluginBaseInterface), cmdPlugin.(PluginCmdWrapperInterface), client

}

func pluginInfoPrint(info *PluginInfo) {

	pluginType := ""
	if info.Type.Is(PluginTypeCmd) {
		pluginType += "终端 "
	}

	if info.Type.Is(PluginTypeWeb) {
		pluginType += "WEB "
	}

	pluginInfoPrintTable := uitable.New()
	pluginInfoPrintTable.AddRow("名称:", info.Name).
		AddRow("作者:", info.Author).
		AddRow("类型:", pluginType).
		AddRow("描述:", info.ShortDesc).
		AddRow("详细说明:", info.Desc).
		AddRow("创建时间:", info.CreateTime.Format("2006-01-02 15:04:05"))

	fmt.Println(pluginInfoPrintTable)
	fmt.Println()
}

func TestWithInteractive(fn PluginServeCallback) {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	res := fn(logger, pluginTestRootCert)
	if res.BasePlugin == nil {
		logger.Error("缺失插件信息内容")
		os.Exit(1)
	}

	pluginInfo, err := res.BasePlugin.Info()
	if err != nil {
		logger.Error("获取插件中的插件信息失败: %s", err.Error())
		os.Exit(2)
	}

	if pluginInfo.Type.Is(PluginTypeCmd) {
		loadCmdPlugin(res)
	} else {
		loadWebPlugin()
	}

}

func loadWebPlugin() {

}

func loadCmdPlugin(res *PluginServeCallbackResult) {
	basePluginInterface, cmdPluginInterface, client := loadPlugin(res)
	defer client.Kill()

	var pluginCmdNameList []string

	currentApp := grumble.New(&grumble.Config{
		Name:                  "test",
		PromptColor:           color.New(color.FgGreen, color.Bold),
		HelpHeadlineColor:     color.New(color.FgGreen),
		HelpHeadlineUnderline: true,
		HelpSubCommands:       true,
		Stdin:                 os.Stdin,
	})

	currentApp.SetInterruptHandler(func(a *grumble.App, count int) {
	})
	currentApp.OnClose(func() error {
		return nil
	})
	currentApp.OnClosing(func() error {
		return nil
	})

	currentApp.AddCommand(&grumble.Command{
		Name: "exit",
		Help: "退出bypt工具",
		Run: func(c *grumble.Context) error {
			c.Stop()
			return nil
		},
	})

	currentApp.AddCommand(&grumble.Command{
		Name: "info",
		Help: "查看当前加载的插件信息",
		Run: func(c *grumble.Context) error {
			info, err := basePluginInterface.Info()
			if err != nil {
				return err
			}
			pluginInfoPrint(info)
			return nil
		},
	})

	currentApp.AddCommand(&grumble.Command{
		Name: "enable",
		Help: "启动当前测试插件",
		Run: func(c *grumble.Context) error {
			if err := basePluginInterface.Start(); err != nil {
				return err
			}

			pluginCmdList := cmdPluginInterface.Commands()
			pluginCmdNameList = pluginCmdLoad(currentApp, pluginCmdList, func(c *grumble.Context) error {
				return cmdPluginInterface.Call(c.Command.Name, c.Flags, c.Args)
			}, func(cmdName string) CompleterFn {
				return func(prefix string, args []string) []string {
					return cmdPluginInterface.Completer(cmdName, prefix, args)
				}
			})
			return nil
		},
	})

	currentApp.AddCommand(&grumble.Command{
		Name: "disable",
		Help: "禁用当前测试的插件",
		Run: func(c *grumble.Context) (err error) {
			defer func() {
				e := recover()
				if e != nil {
					err = errors.New("停止插件失败")
				}
			}()
			basePluginInterface.Stop()
			for i := range pluginCmdNameList {
				cmdName := pluginCmdNameList[i]
				currentApp.Commands().Remove(cmdName)
			}
			pluginCmdNameList = nil
			return
		},
	})

	currentApp.SetPrintASCIILogo(func(a *grumble.App) {
		_, _ = a.Println(` _`)
		_, _ = a.Println(`| |                 _`)
		_, _ = a.Println(`| |__  _   _ ____ _| |_`)
		_, _ = a.Println(`|  _ \| | | |  _ (_   _)`)
		_, _ = a.Println(`| |_) ) |_| | |_| || |_`)
		_, _ = a.Println(`|____/ \__  |  __/  \__)`)
		_, _ = a.Println(`      (____/|_|`)
		_, _ = a.Println()
		_, _ = a.Println("             版本: plugin.test")
		_, _ = a.Println("             作者: 无&痕")
		_, _ = a.Println()
		_, _ = a.Println("应用部署管理平台终端客户端一切只为便捷、高效与可靠的部署和管理应用^-^")
		_, _ = a.Println()
		_, _ = a.Println()
	})

	os.Args = os.Args[:1]
	if err := currentApp.Run(); err != nil {
		return
	}
	os.Exit(0)
}

func pluginCmdLoad(app *grumble.App, cmdList []*CmdInfo, runFn func(c *grumble.Context) error, completer func(cmdName string) CompleterFn) (res []string) {
	res = make([]string, 0, len(cmdList))
	for i := range cmdList {
		cmdInfo := cmdList[i]
		if cmdInfo.Name == "" || cmdInfo.Help == "" {
			panic("命令名称或帮助信息不能为空")
		}

		res = append(res, cmdInfo.Name)
		cmd := &grumble.Command{
			Name:      cmdInfo.Name,
			Help:      cmdInfo.Help,
			LongHelp:  cmdInfo.LongHelp,
			Usage:     cmdInfo.Usage,
			Run:       runFn,
			Completer: completer(cmdInfo.Name),
		}

		if len(cmdInfo.Flags) > 0 {
			cmd.Flags = func(f *grumble.Flags) {
				for k := range cmdInfo.Flags {
					v := cmdInfo.Flags[k]
					if v.Type.Is(DataTypeString) {
						defaultValue, _ := v.DefaultValue.(string)
						if v.Short != "" {
							f.String(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.StringL(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeBool) {
						defaultValue, _ := v.DefaultValue.(bool)
						if v.Short != "" {
							f.Bool(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.BoolL(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeInt) {
						defaultValue, _ := v.DefaultValue.(int)
						if v.Short != "" {
							f.Int(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.IntL(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeInt64) {
						defaultValue, _ := v.DefaultValue.(int64)
						if v.Short != "" {
							f.Int64(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.Int64L(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeUint) {
						defaultValue, _ := v.DefaultValue.(uint)
						if v.Short != "" {
							f.Uint(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.UintL(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeUint64) {
						defaultValue, _ := v.DefaultValue.(uint64)
						if v.Short != "" {
							f.Uint64(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.Uint64L(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeFloat64) {
						defaultValue, _ := v.DefaultValue.(float64)
						if v.Short != "" {
							f.Float64(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.Float64L(v.Long, defaultValue, v.Help)
						}
					} else if v.Type.Is(DataTypeDuration) {
						defaultValue, _ := v.DefaultValue.(time.Duration)
						if v.Short != "" {
							f.Duration(v.Short, v.Long, defaultValue, v.Help)
						} else {
							f.DurationL(v.Long, defaultValue, v.Help)
						}
					}
				}
			}
		}

		if len(cmdInfo.Args) > 0 {
			cmd.Args = func(a *grumble.Args) {
				for i := range cmdInfo.Args {
					argInfo := cmdInfo.Args[i]

					opts := make([]grumble.ArgOption, 0)
					if argInfo.Min != 0 {
						opts = append(opts, grumble.Min(argInfo.Min))
					}

					if argInfo.Max != 0 {
						opts = append(opts, grumble.Max(argInfo.Max))
					}

					if argInfo.Default != "" {
						opts = append(opts, grumble.Default(argInfo.Default))
					}

					if argInfo.Type.Is(DataTypeString) {
						if argInfo.Type.Is(DataTypeList) {
							a.StringList(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.String(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeBool) {
						if argInfo.Type.Is(DataTypeList) {
							a.BoolList(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Bool(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeInt) {
						if argInfo.Type.Is(DataTypeList) {
							a.IntList(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Int(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeInt64) {
						if argInfo.Type.Is(DataTypeList) {
							a.Int64List(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Int64(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeUint) {
						if argInfo.Type.Is(DataTypeList) {
							a.UintList(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Uint(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeUint64) {
						if argInfo.Type.Is(DataTypeList) {
							a.Uint64List(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Uint64(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeFloat64) {
						if argInfo.Type.Is(DataTypeList) {
							a.Float64List(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Float64(argInfo.Name, argInfo.Help, opts...)
						}
					} else if argInfo.Type.Is(DataTypeDuration) {
						if argInfo.Type.Is(DataTypeList) {
							a.DurationList(argInfo.Name, argInfo.Help, opts...)
						} else {
							a.Duration(argInfo.Name, argInfo.Help, opts...)
						}
					}

				}
			}

		}

		app.AddCommand(cmd)
	}
	return
}
