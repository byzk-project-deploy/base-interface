package rpcinterfaces

import (
	"errors"
	"net/rpc"
	"time"
)

var ErrTimeout = errors.New("业务接口访问超时")

func errHandle(err error) error {
	if _, ok := err.(rpc.ServerError); ok {
		return err
	}
	panic(err)
}

func withTimeout(duration time.Duration, fn func() error) error {
	ch := make(chan error, 1)
	go func() {
		ch <- fn()
	}()

	after := time.After(duration)
	select {
	case <-after:
		return ErrTimeout
	case err := <-ch:
		return err
	}
}
