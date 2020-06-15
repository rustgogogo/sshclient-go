package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

//https://github.com/pkg/sftp
//
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

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer sftpClient.Close()

	//获取当前目录
	cwd, err := sftpClient.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("当前目录：", cwd)

	//显示文件/目录详情
	fi, err := sftpClient.Lstat(cwd)
	log.Println(fi)

	{
		//上传文件(将本地file.dat文件通过sftp传到远程服务器)
		remoteFileName := fmt.Sprintf("upload_by_sftp_%d.dat", rand.Int())
		remoteFile, err := sftpClient.Create(sftp.Join(cwd, remoteFileName))
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer remoteFile.Close()

		localFileName := "file.dat"
		//打开本地文件file.dat
		localFile, err := os.Open(localFileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer localFile.Close()

		//本地文件流拷贝到上传文件流
		n, err := io.Copy(remoteFile, localFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		//获取本地文件大小
		localFileInfo, err := os.Stat(localFileName)
		if err != nil {
			log.Fatalln(err.Error())
		}

		log.Printf("文件上传成功[%s->%s]本地文件大小：%s，上传文件大小：%s", localFileName, remoteFileName, formatFileSize(localFileInfo.Size()), formatFileSize(n))

		//计算文件MD5
		//windows计算：certutil -hashfile .\file.dat MD5
		//linux计算：md5sum upload_by_sftp_5577006791947779410.dat
	}
	{
		//下载文件
		//将远程服务器的/bin/bash文件下载到本地
		remoteFileName := "/bin/bash"
		remoteFile, err := sftpClient.Open(remoteFileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer remoteFile.Close()

		localFileName := "local-bash"
		localFile, err := os.Create(localFileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer localFile.Close()
		n, err := io.Copy(localFile, remoteFile)
		if err != nil {
			log.Fatalln(err.Error())
		}

		//获取远程文件大小
		remoteFileInfo, err := sftpClient.Stat(remoteFileName)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("文件下载成功[%s->%s]远程文件大小：%s，下载文件大小：%s", remoteFileName, localFileName, formatFileSize(remoteFileInfo.Size()), formatFileSize(n))
	}
}

// 字节的单位转换 保留两位小数
func formatFileSize(s int64) (size string) {
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
