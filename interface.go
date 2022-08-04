package rpcinterfaces

import "github.com/byzk-project-deploy/grumble"

type DataType uint

func (d DataType) Is(t DataType) bool {
	return d&t == t
}

const (
	// DataTypeString  字符串
	DataTypeString DataType = 1 << iota
	DataTypeBool
	DataTypeInt
	DataTypeInt64
	DataTypeUint
	DataTypeUint64
	DataTypeFloat64
	DataTypeDuration
	DataTypeList
)

type CmdFlagInfo struct {
	Type         DataType
	Long         string
	Short        string
	DefaultValue interface{}
	Help         string
}

type CmdArgInfo struct {
	Type    DataType
	Name    string
	Help    string
	Min     int
	Max     int
	Default string
}

type CmdRunFn func(flags grumble.FlagMap, args grumble.ArgMap) error
type CompleterFn func(prefix string, args []string) []string

type CmdInfo struct {
	// Name 名称
	Name string
	// Help 帮助
	Help string
	//LongHelp 长文本帮助
	LongHelp string
	// Usage 使用说明
	// Sample: start [OPTIONS] CONTAINER [CONTAINER...]
	Usage string
	// Flags 命令行选项
	Flags map[string]*CmdFlagInfo
	// Args 命令行参数
	Args []*CmdArgInfo
	// Run 运行命令
	Run CmdRunFn `json:"-"`
	// Completer 自动补全
	Completer CompleterFn `json:"-"`
}

type PluginCmdInterface interface {
	// Commands 获取终端命令
	Commands() []*CmdInfo
	// Completer 命令补全
	//Completer(cmdName, prefix string, args []string) []string
	// Call 命令调用
	//Call(cmdName string, flags grumble.FlagMap, args grumble.ArgMap) error
}

type PluginCmdWrapperInterface interface {
	PluginCmdInterface
	// Call 命令调用
	Call(cmdName string, flags grumble.FlagMap, args grumble.ArgMap) error
	// Completer 命令补全
	Completer(cmdName, prefix string, args []string) []string
}

type PluginWebInterface interface {
}

// PluginBaseInterface 插件信息接口
type PluginBaseInterface interface {
	// Info 插件信息
	Info() (*PluginInfo, error)
	// Start 启动插件
	Start() error
	// Stop 停止插件
	Stop()
}
