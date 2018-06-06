package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"net"
	"io"
)

const versionNumber = "1.0.0#20180606"

func main() {
	version := flag.Bool("version", false, "Show program's version number and exit")
	listen := flag.Bool("l", false, "Listening on the server")
	addr := flag.String("a", "", "Address to use")
	port := flag.Int("p", 0, "Port to use ")
	help := flag.Bool("h", false, "Show usage")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n  %s [Options]\n\nOptions\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if *version {
		fmt.Println(versionNumber)
		return
	}
	if *help || *port == 0 || (*addr == "" && *listen == false) {
		flag.Usage()
		return
	}



	if *listen {
		listenS(*addr, *port)
	}else {
		connectS(*addr, *port)
	}

}

func listenS(addr string, port int) {
	if addr == "" {
		addr = "localhost"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func connectS(addr string, port int) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go tran(os.Stdout, conn)
	tran(conn, os.Stdin)
}

func handleConn(conn net.Conn) {
	go tran(os.Stdout, conn)
	tran(conn, os.Stdin)
}

func tran(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}