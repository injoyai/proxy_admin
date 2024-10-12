package main

import "github.com/injoyai/proxy/core"

type Tunnel struct {
	ID     string `json:"id"`             //主键
	Name   string `json:"name"`           //名称
	Memo   string `json:"memo"`           //备注
	SN     string `json:"sn"`             //唯一标识
	Target string `json:"target"`         //目标地址
	Online bool   `json:"online" orm:"-"` //是否在线
}

func (this *Tunnel) Resp() {
	_, ok := Server.Clients.Get(this.SN)
	this.Online = ok
}

func (this *Tunnel) GetTarget() string {
	if len(this.Target) == 0 {
		return ":80"
	}
	return this.Target
}

func (this *Tunnel) Dialer() *core.Dial {
	return core.NewDialTCP(this.GetTarget())
}
