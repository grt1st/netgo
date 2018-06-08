package netextends

import (
	"net"
	"os/exec"
	"log"
	"os"
	"fmt"
	"github.com/grt1st/netgo/utils"
	"bufio"
	"io"
)


func ListenNormalMode(listener net.Listener) {
	allClients := make(map[net.Conn]int)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	inMessages := make(chan string)
	outMessages := make(chan string)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("acc")
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
				log.Print(err)
				os.Exit(1)
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
						os.Exit(1)
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
}


func ListenExecMode(listener net.Listener, exeCmd string) {
	newConnections := make(chan net.Conn)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("acc")
				os.Exit(1)
			}
			newConnections <- conn
		}
	}()

	for {
		select {
		case conn := <-newConnections:
			go func(conn net.Conn) {
				cmd := exec.Command(exeCmd)

				//创建获取命令输出管道
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
					os.Exit(1)
				}
				stdin, err := cmd.StdinPipe()
				if err != nil {
					log.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
					os.Exit(1)
				}

				go utils.Transform(conn, stdout)
				go utils.Transform(stdin, conn)

				//执行命令
				if err := cmd.Start(); err != nil {
					log.Println("Error:The command is err,", err)
					os.Exit(1)
				}

				//wait 方法会一直阻塞到其所属的命令完全运行结束为止
				if err := cmd.Wait(); err != nil {
					fmt.Println("wait:", err.Error())
					os.Exit(1)
				}

				os.Exit(0)
			}(conn)
		}
	}
}