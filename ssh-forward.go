package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	//本地监听地址
	laddr = "localhost:5900"
	//远程服务地址
	raddr   = "localhost:5900"
	sshaddr = "192.168.1.8:22"
)

/**
基于ssh连接的端口转发
https://tools.ietf.org/html/rfc4254#page-16 TCP/IP Port Forwarding

通过ssh隧道将远程服务器(192.168.1.8)上的5900端口映射到本地5900端口
*/
func main() {

	sshConfig := &ssh.ClientConfig{
		User: "pi",
		Auth: []ssh.AuthMethod{
			ssh.Password("123456"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		ClientVersion:   "",
		Timeout:         10 * time.Second,
	}

	//监听本地映射端口
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		//客户端连接
		go portForward(conn, sshConfig)
	}
}

/**
conn: 客户端到本地映射端口的连接
sshConfig: ssh配置
*/
func portForward(conn net.Conn, sshConfig *ssh.ClientConfig) {
	defer conn.Close()

	//建立与SSH服务器的连接
	sshClient, err := ssh.Dial("tcp", sshaddr, sshConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer sshClient.Close()

	//建立ssh到后端服务的连接
	forwardConn, err := sshClient.Dial("tcp", raddr)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("ssh端口映射隧道建立成功")

	defer forwardConn.Close()

	var wait2close = sync.WaitGroup{}
	wait2close.Add(1)

	go func() {
		n, err := io.Copy(forwardConn, conn)
		if err != nil {
			log.Fatalln(err.Error())
			wait2close.Done()
		}
		log.Printf("入流量共%s", formatFlowSize(n))
	}()

	go func() {
		n, err := io.Copy(conn, forwardConn)
		if err != nil {
			log.Fatalln(err.Error())
			wait2close.Done()
		}
		log.Printf("出流量共%s", formatFlowSize(n))
	}()

	wait2close.Wait()
	log.Println("ssh端口映射隧道关闭")
}

// 字节的单位转换 保留两位小数
func formatFlowSize(s int64) (size string) {
	if s < 1024 {
		return fmt.Sprintf("%.2fB", float64(s)/float64(1))
	} else if s < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(1024))
	} else if s < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(1024*1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(1024*1024*1024*1024))
	} else { //if s < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(s)/float64(1024*1024*1024*1024*1024))
	}
}
