package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"net"
	"io"
	"bufio"
	"sync"
	"os/signal"
	"golang.org/x/net/netutil"
	"os/exec"
)

const versionNumber = "1.0.0#20180606"

var wg = sync.WaitGroup{}

func main() {
	version := flag.Bool("version", false, "Show program's version number and exit")
	listen := flag.Bool("l", false, "Listening on the server")
	addr := flag.String("a", "", "Address to use")
	port := flag.Int("p", 0, "Port to use ")
	htmlFlag := flag.Bool("html", false, "Send html request of GET")
	exeCmd := flag.String("e", "", "")
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
		listenS(*addr, *port, *exeCmd)
	}else {
		connectS(*addr, *port, *htmlFlag, *exeCmd)
	}

}

func listenS(addr string, port int, exeCmd string) {

	defer wg.Done()
	stopChan := make(chan os.Signal) // 接收系统中断信号
	signal.Notify(stopChan, os.Interrupt)

	if addr == "" {
		addr = "localhost"
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listener.Close()
	listener = netutil.LimitListener(listener, 1)

	allClients := make(map[net.Conn]int)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	inMessages := make(chan string)
	outMessages := make(chan string)

	go func() {
		<-stopChan
		for c := range allClients {
			delete(allClients, c)
			c.Close()
		}
		if err = listener.Close(); err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("acc")
				os.Exit(1)
			}
			wg.Add(1)
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
						if err != io.EOF {
							log.Print("err")
							log.Print(err)
						}else if err == io.EOF {
							//  关闭连接
							log.Print("eof")
							os.Exit(0)
						}
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
						log.Print("close")
						log.Print(err)
						deadConnections <- c
					}
				}(c)
			}
		case conn := <-deadConnections:
			delete(allClients, conn)
			conn.Close()
		}
	}

	wg.Wait()

}

func connectS(addr string, port int, htmlFlag bool, exeCmd string) {

	stopChan := make(chan os.Signal) // 接收系统中断信号
	signal.Notify(stopChan, os.Interrupt)

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		log.Fatal(err)
	}
	//defer conn.Close()
	go func() {
		<-stopChan
		if err = conn.Close(); err != nil {
			log.Print(err)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	if htmlFlag {
		fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
		io.Copy(os.Stdout, conn)
		return
	}else if exeCmd != ""{
		cmd := exec.Command(exeCmd)

		//创建获取命令输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
			return
		}
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
			return
		}

		go tran(conn, stdout)
		go tran(stdin, conn)

		//执行命令
		if err := cmd.Start(); err != nil {
			log.Println("Error:The command is err,", err)
			return
		}

		//wait 方法会一直阻塞到其所属的命令完全运行结束为止
		if err := cmd.Wait(); err != nil {
			fmt.Println("wait:", err.Error())
			return
		}

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

