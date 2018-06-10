package netextends

import (
	"os"
	"os/signal"
	"net"
	"fmt"
	"github.com/grt1st/netgo/utils"
	"log"
	"io"
)

func RemotePortForward(addr string, port int, rhost string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}
	conn1, err := net.Dial("tcp", rhost)
	if err != nil {
		log.Fatal(err)
	}

	transformToEach(conn, conn1)
}

func ServerAndServer(addr string, port1,port2 int,) {
	stopChan := make(chan os.Signal) // 接收系统中断信号
	signal.Notify(stopChan, os.Interrupt)

	if addr == "" {
		addr = "localhost"
	}
	listener1, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port1))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	listener2, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port2))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	listener1 = utils.LimitListener(listener1, 1)
	listener2 = utils.LimitListener(listener2, 1)

	go func() {
		<-stopChan
		if err = listener1.Close(); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		if err = listener2.Close(); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	for {
		conn, err := listener1.Accept()
		if err != nil {
			fmt.Println("acc")
			os.Exit(1)
		}
		conn1, err := listener2.Accept()
		if err != nil {
			fmt.Println("acc")
			os.Exit(1)
		}
		go transformToEach(conn, conn1)
	}
}

func transformToEach(conn, conn1 net.Conn) {
	transformToConn(conn, conn1)
	transformToConn(conn1, conn)
	defer conn.Close()
	defer conn1.Close()
}

func transformToConn(l1, l2 net.Conn) {
	go func() {
		io.Copy(l2, l1)
		os.Exit(0)
	}()
	//go utils.Transform(os.Stdout, conn)
	utils.Transform(l1, l2)
}

