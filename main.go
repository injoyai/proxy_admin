package main

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/goutil/frame/gf"
	"github.com/injoyai/goutil/frame/in"
	"github.com/injoyai/logs"
	"github.com/injoyai/minidb"
	"github.com/injoyai/proxy/core"
	"github.com/injoyai/proxy/tunnel"
	"net"
)

var (
	DB     = minidb.New("./database", "default", minidb.WithID("ID"))
	Server = tunnel.Server{
		Listen:      core.NewListenTCP(7000), //给隧道客户端的端口
		OnRegister:  nil,
		OnProxy:     nil,
		OnConnected: nil,
		OnClosed:    nil,
	}
	Bridge net.Listener
)

func init() {
	logs.PrintErr(DB.Sync(new(Tunnel)))
	go Server.Run()
}

func main() {

	s := gf.New().SetPort(8080)

	s.Group("/api", func(g *ghttp.RouterGroup) {

		g.GET("/tunnel/list", func(r *ghttp.Request) {
			data := []*Tunnel(nil)
			co, err := DB.FindAndCount(&data)
			in.CheckErr(err)
			for _, v := range data {
				v.Resp()
			}
			in.Succ(data, co)
		})

		g.POST("/tunnel", func(r *ghttp.Request) {
			req := &Tunnel{}
			err := r.Parse(req)
			in.CheckErr(err)
			if len(req.SN) == 0 {
				in.Err("sn不能为空")
			}
			if len(req.Name) == 0 {
				in.Err("名称不能为空")
			}
			if req.ID == "" {
				err = DB.Insert(req)
				in.Err(err)
			}
			err = DB.Where("ID=?", req.ID).Update(req)
			in.Err(err)
		})

		g.DELETE("/tunnel", func(r *ghttp.Request) {
			id := r.GetString("id")
			err := DB.Where("ID=?", id).Delete(new(Tunnel))
			in.Err(err)
		})

	})

	s.Run()
}
