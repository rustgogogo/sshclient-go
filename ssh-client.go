package main

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"time"
)

/**
golang版本的SSH客户端
SSH协议RFC文档
https://tools.ietf.org/html/rfc4254

一个ssh连接可以打开多个会话session
linux tty和pty区别
开机后登录系统的终端称为tty
远程登录的终端称为pty
pts是pty的实现方式
w命令可以显示当前系统登录的终端列表
针对交互式会话的操作
1.请求伪终端 pty-req
2.X11转发 x11-req
3.X11通道 x11
4.环境变量 env
5.启动shell或命令 shell/exec/subsystem

默认不支持上下键和tab键，还不支持clear清屏指令
通过VT100终端支持tab和clear指令
VT100终端包括一些控制符，可以在终端中显示不同颜色，支持光标控制，清屏指令等
http://www.termsys.demon.co.uk/vtansi.htm
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
	//建立与SSH服务器的连接
	sshClient, err := ssh.Dial("tcp", "192.168.1.8:22", sshConfig)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer sshClient.Close()
	log.Println("sessionId: ", sshClient.SessionID())
	log.Println("user: ", sshClient.User())
	log.Println("ssh server version: ", string(sshClient.ServerVersion()))
	log.Println("ssh client version: ", string(sshClient.ClientVersion()))

	//打开交互式会话(A session is a remote execution of a program.)
	//https://tools.ietf.org/html/rfc4254#page-10
	session, err := sshClient.NewSession()
	if err != nil {
		log.Fatalln("Failed to create ssh session", err)
	}

	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     //打开回显
		ssh.TTY_OP_ISPEED: 14400, //输入速率 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //输出速率 14.4kbaud
		ssh.VSTATUS:       1,
	}

	//使用VT100终端来实现tab键提示，上下键查看历史命令，clear键清屏等操作
	//VT100 start
	//windows下不支持VT100
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer terminal.Restore(fd, oldState)
	//VT100 end

	termWidth, termHeight, err := terminal.GetSize(fd)

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	//打开伪终端
	//https://tools.ietf.org/html/rfc4254#page-11
	err = session.RequestPty("xterm", termHeight, termWidth, modes)
	if err != nil {
		log.Fatalln(err.Error())
	}

	//启动一个远程shell
	//https://tools.ietf.org/html/rfc4254#page-13
	err = session.Shell()
	if err != nil {
		log.Fatalln(err.Error())
	}

	//等待远程命令结束或远程shell退出
	err = session.Wait()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
