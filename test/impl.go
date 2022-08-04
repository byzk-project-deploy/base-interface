package main

import (
	"fmt"
	"github.com/byzk-project-deploy/base-interface"
	"github.com/byzk-project-deploy/grumble"
	"github.com/hashicorp/go-hclog"
	"strings"
	"time"
)

type testPluginBaseImpl struct {
	logger hclog.Logger
}

func (t *testPluginBaseImpl) Start() error {
	fmt.Println("测试插件已启用")
	return nil
}

func (t *testPluginBaseImpl) Stop() {
	fmt.Println("测试插件已停用")
}

func (t *testPluginBaseImpl) Info() (*rpcinterfaces.PluginInfo, error) {
	return &rpcinterfaces.PluginInfo{
		Name:       "test",
		Author:     "bypt",
		ShortDesc:  "bypt plugin functional test",
		Desc:       "test bypt cmd plugin and web plugin functional",
		CreateTime: time.Now(),
		Type:       rpcinterfaces.PluginTypeCmd,
	}, nil
}

type testPluginCmdImpl struct {
	logger hclog.Logger
}

func (t *testPluginCmdImpl) Commands() []*rpcinterfaces.CmdInfo {
	return []*rpcinterfaces.CmdInfo{
		{
			Name:     "echo",
			Help:     "打印字符",
			LongHelp: "将命令之后的所有内容进行打印",
			Usage:    "echo a b c",
			Args: []*rpcinterfaces.CmdArgInfo{
				{
					Name: "args",
					Type: rpcinterfaces.DataTypeString | rpcinterfaces.DataTypeList,
					Help: "要打印的内容",
					Min:  1,
				},
			},
			Completer: func(prefix string, args []string) []string {
				if len(args) > 0 {
					return nil
				}
				return []string{"test"}
			},
			Run: func(flags grumble.FlagMap, args grumble.ArgMap) error {
				list := args.StringList("args")
				fmt.Println(strings.Join(list, " "))
				return nil
			},
		},
	}
}
