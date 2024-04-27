package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// 监听连接
func StartProxy(connPool *ConnPool, addr string) {
	n := "unix"
	if strings.Contains(addr, ":") {
		n = "tcp"
	} else {
		// remove addr file
		os.Remove(addr)
	}

	l, err := net.Listen(n, addr)
	if err != nil {
		connPool.log.Fatal("listen error:", err)
		return
	}
	defer l.Close()

	for {
		local, err := l.Accept()
		if err != nil {
			connPool.log.Warn("accept error:", err)
			continue
		}

		go HandlerData(connPool, local)
	}

}

// 数据交换方法
func HandlerData(connPool *ConnPool, local net.Conn) {

	connPool.log.Debug("remote addr:", local.RemoteAddr(), connPool.Stats())

	conn, err := connPool.Get()
	if err != nil {
		local.Close()
		fmt.Println("get conn error:", err)
		return
	}

	forceClose := conn.SwapData(local)
	local.Close()
	connPool.Put(conn, forceClose)

}
