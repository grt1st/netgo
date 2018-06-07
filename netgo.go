package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"net"
	"io"
	"bufio"
)

const versionNumber = "1.0.0#20180606"

func main() {
	version := flag.Bool("version", false, "Show program's version number and exit")
	listen := flag.Bool("l", false, "Listening on the server")
	addr := flag.String("a", "", "Address to use")
	port := flag.Int("p", 0, "Port to use ")
	html := flag.Bool("html", false, "Send html request of GET")
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
		connectS(*addr, *port, *html)
	}

}

func listenS(addr string, port int) {
	if addr == "" {
		addr = "localhost"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	allClients := make(map[net.Conn]int)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	inMessages := make(chan string)
	outMessages := make(chan string)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			message, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			outMessages <- message
		}
	}()

	for {
		select {
		case conn := <-newConnections:
			allClients[conn] = 1
			go func(conn net.Conn) {
				reader := bufio.NewReader(conn)
				for {
					message, err := reader.ReadString('\n')
					if err != nil {
						break
					}
					inMessages <- message
				}
				deadConnections <- conn
			}(conn)
		case message := <-inMessages:
			fmt.Print(message)
		case message := <-outMessages:
			for c := range allClients{
				go func(c net.Conn) {
					_, err := c.Write([]byte(message))
					if err != nil {
						deadConnections <- c
					}
				}(c)
			}
		case conn := <-deadConnections:
			delete(allClients, conn)
			conn.Close()
		}
	}


}

func connectS(addr string, port int, html bool) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	if html {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		io.Copy(os.Stdout, conn)
		return
	}
	go tran(os.Stdout, conn)
	tran(conn, os.Stdin)
}

func tran(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}