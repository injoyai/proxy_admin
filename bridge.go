package main

import (
	"bufio"
	"github.com/injoyai/logs"
	"github.com/injoyai/proxy/core"
	"net"
	"net/http"
)

func init() {
	var err error
	Bridge, err = net.Listen("tcp", ":8090")
	logs.PanicErr(err)

	go func() {
		for {
			conn, err := Bridge.Accept()
			if err != nil {
				return
			}
			go func() {
				defer conn.Close()
				err := handlerEasy(conn)
				logs.PrintErr(err)
			}()
		}
	}()
}

func handlerEasy(conn net.Conn) error {
	buf := bufio.NewReader(conn)

	r, err := http.ReadRequest(buf)
	if err != nil {
		return err
	}

	co, err := r.Cookie("sn")
	if err != nil {
		return err
	}

	sn := co.Value

	c, ok := Server.Clients.Get(sn)
	if !ok {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		return nil
	}

	t := c.(*core.Tunnel)

	return t.DialAndSwap(
		conn.RemoteAddr().String(),
		core.NewDialTCP(":80"),
		conn,
	)

}
