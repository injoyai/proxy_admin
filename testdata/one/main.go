package main

import (
	"github.com/injoyai/logs"
	"github.com/injoyai/proxy/core"
	"github.com/injoyai/proxy/tunnel"
	"time"
)

func main() {
	sn := "001"
	t := tunnel.Client{
		Dialer: core.NewDialTCP(":7000"),
		Register: &core.RegisterReq{
			Key: sn,
		},
	}
	for {
		logs.PrintErr(t.Run())
		<-time.After(time.Second * 2)
	}
}
