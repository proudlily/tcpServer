package utils

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

type TcpServerEvent interface {
	//IPC消息封装回调
	OnMqResponse(MsgSender, MsgReceiver int32, msgHead string, msgData []byte)
	//基类监听回调初始化自己数据
	OnCallBackInit()
	//每一个连接的初始化
	OnTcpHCCallBackInit(tcphc *TCPHandleConnection)
}

type TcpHandleConnectionEvent interface {
	//建立连接
	OnEventTCPNetworkBind(Remoteip string)
	//读取数据
	OnEventTCPNetworkRead(req proto.Message) (bool, proto.Message)
	//关闭网络连接
	OnEventTCPNetworkShut(err string)
}

type SocketManage struct {
	ServerConn net.Listener
	tcpEvent   TcpServerEvent
	ConnTimes  int64
}

type TCPHandleConnection struct {
	NetBuf          *NetBuffer //msg 解析
	HandleConnEvent TcpHandleConnectionEvent
	Conn            net.Conn
	Ip              string
	version         int
	errdescribe     string
}

func NewTCPCreate(ListenIp string, tcpserverevent TcpServerEvent) *SocketManage {
	this := &SocketManage{
		tcpEvent:  tcpserverevent,
		ConnTimes: 0,
	}

	this.tcpEvent.OnCallBackInit()

	//监听用户端口
	var err error
	this.ServerConn, err = net.Listen("tcp", ListenIp)
	if err != nil {
		if SysLog != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("NewHandleControl.net.Listen(tcp, %s) error: %s", ListenIp, err.Error()))
		}
		return nil
	} else {
		if SysLog != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("NewHandleControl listen sueccss = %s", ListenIp))
		}
	}
	return this
}

//启动serverHandle
func (this *SocketManage) TcpStart() {
	for {
		//没有err，即代表建立了一个TCP的连接
		conn, err := this.ServerConn.Accept()
		if err != nil && SysLog != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("error accepting: %s", err.Error()))
			break
		}

		this.ConnTimes++
		//创建一个gw的server协程
		go MakeNewHandleConnection(this, conn).serve()
		if this.ConnTimes > 999999999 {
			this.ConnTimes = 0
		}
	}
}

func MakeNewHandleConnection(socketman *SocketManage, conn net.Conn) *TCPHandleConnection {
	Tcphc := &TCPHandleConnection{
		Conn:        conn,
		Ip:          "",
		version:     1,
		errdescribe: "",
		NetBuf:      NewNetBuffer([]byte{69, 123, 132, 104, 67, 95, 33, 74, 120, 131, 61, 101, 55, 101, 69, 44}),
	}

	//回调初始化一些子连接数据
	socketman.tcpEvent.OnTcpHCCallBackInit(Tcphc)

	return Tcphc
}

func (this *TCPHandleConnection) OnClose() {
	this.Conn.Close()
	this.HandleConnEvent.OnEventTCPNetworkShut(this.errdescribe)
}

func (this *TCPHandleConnection) serve() {
	defer this.OnClose()
	this.Ip = strings.Split(this.Conn.RemoteAddr().String(), ":")[0]
	this.HandleConnEvent.OnEventTCPNetworkBind(this.Ip)
LoopConn:
	for {
		//设置超时
		this.Conn.SetReadDeadline(time.Now().Add(time.Second * 130000)) //13秒没有心跳就超时

		//接收数据
		if err := this.NetBuf.Read(this.Conn); err != nil {
			this.errdescribe = err.Error()
			break
		}
		//处理数据
		for {
			if version, msg, _ := this.NetBuf.GetAMsg(); msg != nil {
				this.version = version
				closeConn, resp := this.HandleConnEvent.OnEventTCPNetworkRead(msg)
				if this.OnResp(closeConn, resp) {
					this.errdescribe = fmt.Sprintf("IP:[%s] socket Msgid [%d] return false ", this.Ip, GetMsgCode(msg))
					break LoopConn
				}
			} else {
				break
			}
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func (this *TCPHandleConnection) OnResp(closeConn bool, resp proto.Message) bool {
	if resp != nil { //发响应回去
		if _, err := this.SendMsg(resp); err != nil {
			closeConn = true
		}
	}
	return closeConn
}

// SendMsg 发送msg给用户
func (this *TCPHandleConnection) SendMsg(resp proto.Message) (int, error) {
	if resp != nil {
		n, err := this.Conn.Write(this.toData(resp))
		if err != nil && SysLog != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("SendMsg(%s) ERROR:%s", ProtoToString(resp), err.Error()))
		} else {
			if SysLog != nil {
				if GetMsgCode(resp) != 2 {
					SysLog.PutLineAsLog(fmt.Sprintf("SendMsg:OK %s CodeID=%d", resp.String(), GetMsgCode(resp)))
				}
			}
		}
		return n, err
	}
	return 0, nil
}

//pkt转成加密的二进制数据
func (this *TCPHandleConnection) toData(pkt proto.Message) []byte {
	return ToData(pkt, this.NetBuf.Key, this.version, 0)
}
