package rpcinterfaces

import "net/rpc"

func errHandle(err error) error {
	if _, ok := err.(rpc.ServerError); ok {
		return err
	}
	panic(err)
}
