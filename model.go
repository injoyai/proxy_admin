package main

type Tunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Memo   string `json:"memo"`
	SN     string `json:"sn"`
	Online bool   `json:"online" xorm:"-"`
}

func (this *Tunnel) Resp() {
	_, ok := Server.Clients.Get(this.SN)
	this.Online = ok
}
