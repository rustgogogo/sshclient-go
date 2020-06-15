golang实现ssh客户端

**ssh模式**


1.shell模式

2.port forward 端口转发

3.sftp 基于ssh安全的文件传输协议 (SSH File Transfer Protocol)

ssh-client.go shell模式

ssh-forward.go 端口转发

ssh-sftp.go 文件传输

---

shell模式

导入golang.org/x/crypto/ssh包

VT100终端需要导入golang.org/x/crypto/ssh/terminal包

默认不支持上下键和tab键，还不支持clear清屏指令

通过VT100终端支持tab和clear指令

VT100终端包括一些控制符，可以在终端中显示不同颜色，支持光标控制，清屏指令等

http://www.termsys.demon.co.uk/vtansi.htm

---

port forward模式

导入golang.org/x/crypto/ssh包

通过ssh隧道来映射端口

---
sftp模式

导入https://github.com/pkg/sftp包

