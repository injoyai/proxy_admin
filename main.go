package main

import (
	_ "embed"
	"github.com/gogf/gf/net/ghttp"
	"github.com/injoyai/base/maps"
	"github.com/injoyai/goutil/frame/gf"
	"github.com/injoyai/goutil/frame/in"
	"github.com/injoyai/logs"
	"github.com/injoyai/minidb"
	"github.com/injoyai/proxy/core"
	"github.com/injoyai/proxy/tunnel"
	"net"
)

//go:embed index.html
var Index []byte

//go:embed offline.html
var Offline []byte

var (
	DB     = minidb.New("./database", "default", minidb.WithID("ID"))
	Server = tunnel.Server{
		Listen: core.NewListenTCP(7000), //给隧道客户端的端口
		OnRegister: func(tun *core.Tunnel, reg *core.RegisterReqExtend) error {
			tun.SetKey(reg.Key)
			reg.Listen = &core.Listen{}
			return nil
		},
	}
	Bridge       net.Listener
	BridgeTarget = maps.NewSafe()
	DefaultUID   = "default"
	TunnelCache  = maps.NewSafe()
)

func init() {
	core.DefaultLog.SetLevel(core.LevelInfo)
	logs.PrintErr(DB.Sync(new(Tunnel)))
	go Server.Run()
	data := []*Tunnel(nil)
	err := DB.Find(&data)
	logs.PanicErr(err)
	for _, v := range data {
		TunnelCache.Set(v.SN, v)
	}
}

func main() {

	s := gf.New().SetPort(8080)

	s.GET("/", func(r *ghttp.Request) {
		in.Html(200, Index)
	})

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
			if err == nil {
				TunnelCache.Set(req.SN, req)
			}
			in.Err(err)
		})

		g.DELETE("/tunnel", func(r *ghttp.Request) {
			id := r.GetString("id")
			err := DB.Where("ID=?", id).Delete(new(Tunnel))
			if err == nil {
				TunnelCache.Range(func(key, value interface{}) bool {
					if value.(*Tunnel).ID == id {
						TunnelCache.Del(key)
						return false
					}
					return true
				})
			}
			in.Err(err)
		})

		g.POST("/bridge", func(r *ghttp.Request) {
			sn := r.GetString("sn")
			BridgeTarget.Set(DefaultUID, sn)
			t := Server.GetTunnel(sn)
			if t == nil {
				in.Err("隧道不在线")
			}
			in.Succ(nil)
		})

	})

	s.Run()
}
