package netextends

import (
	"fmt"
	"io"
	"os"
	"net"
	"log"
	"os/exec"
	"github.com/grt1st/netgo/utils"
)


func ConnectNormalMode(conn net.Conn) {
	go func() {
		io.Copy(os.Stdout, conn)
		os.Exit(0)
	}()
	//go utils.Transform(os.Stdout, conn)
	utils.Transform(conn, os.Stdin)
}


func ConnectExecMode(conn net.Conn, exeCmd string) {
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

	go func() {
		utils.Transform(stdin, conn)
		os.Exit(0)
	}()
	go utils.Transform(conn, stdout)

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
}


func ConnectHtmlMode(conn net.Conn) {
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	io.Copy(os.Stdout, conn)
}

