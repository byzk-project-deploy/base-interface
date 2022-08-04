package rpcinterfaces

import (
	"encoding/json"
	"fmt"
	"github.com/byzk-project-deploy/go-plugin"
	"github.com/byzk-project-deploy/grumble"
	"net/rpc"
)

var cmdMap map[string]*CmdInfo

type PluginCmdImpl struct {
	impl PluginCmdInterface
}

func (p *PluginCmdImpl) Server(broker *plugin.MuxBroker) (interface{}, error) {
	return &pluginCmdRpcServer{impl: p.impl}, nil
}

func (p *PluginCmdImpl) Client(broker *plugin.MuxBroker, client *rpc.Client) (interface{}, error) {
	return &pluginCmdRpc{client: client}, nil
}

type cmdCallRpcArgs struct {
	CmdName string
	Flags   grumble.FlagMap
	Args    grumble.ArgMap
}

type cmdCompleterRpcArgs struct {
	CmdName string
	Prefix  string
	Args    []string
}

type pluginCmdRpcServer struct {
	impl PluginCmdInterface
}

func (p pluginCmdRpcServer) Commands(args interface{}, resp *[]byte) error {
	cmds := p.impl.Commands()
	cmdLen := len(cmds)
	cmdMap = make(map[string]*CmdInfo, cmdLen)

	if cmdLen == 0 {
		return nil
	}

	for i := range cmds {
		cmd := cmds[i]
		if cmd.Name != "" {
			cmdMap[cmd.Name] = cmd
		}
	}
	*resp, _ = json.Marshal(cmds)
	return nil
}

func (p pluginCmdRpcServer) Call(args []byte, resp *interface{}) error {
	var rpcArgs *cmdCallRpcArgs
	if err := json.Unmarshal(args, &rpcArgs); err != nil {
		return nil
	}

	cmdInfo, ok := cmdMap[rpcArgs.CmdName]
	if !ok {
		return fmt.Errorf("命令[%s]不存在", rpcArgs.CmdName)
	}

	if cmdInfo.Run == nil {
		return nil
	}

	return cmdInfo.Run(rpcArgs.Flags, rpcArgs.Args)
}

func (p pluginCmdRpcServer) Completer(args []byte, resp *[]string) error {
	var rpcArgs *cmdCompleterRpcArgs
	if err := json.Unmarshal(args, &rpcArgs); err != nil {
		return nil
	}

	cmdInfo, ok := cmdMap[rpcArgs.CmdName]
	if !ok {
		return nil
	}

	if cmdInfo.Completer == nil {
		return nil
	}

	*resp = cmdInfo.Completer(rpcArgs.Prefix, rpcArgs.Args)
	return nil
}

type pluginCmdRpc struct {
	client *rpc.Client
}

func (p *pluginCmdRpc) Call(cmdName string, flags grumble.FlagMap, args grumble.ArgMap) error {
	marshal, _ := json.Marshal(&cmdCallRpcArgs{
		CmdName: cmdName,
		Flags:   flags,
		Args:    args,
	})
	return p.client.Call("Plugin.Call", marshal, new(interface{}))
}

func (p *pluginCmdRpc) Commands() (res []*CmdInfo) {
	var resBytes []byte
	if err := p.client.Call("Plugin.Commands", new(interface{}), &resBytes); err != nil {
		panic(err)
	}

	_ = json.Unmarshal(resBytes, &res)
	return
}

func (p *pluginCmdRpc) Completer(cmdName, prefix string, args []string) (res []string) {
	marshal, _ := json.Marshal(&cmdCompleterRpcArgs{
		CmdName: cmdName,
		Prefix:  prefix,
		Args:    args,
	})
	_ = p.client.Call("Plugin.Completer", marshal, &res)
	return
}
