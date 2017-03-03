package tcpClient

import (
	"fmt"
	"gt_msg"
	"net"
	"time"
	"utils"

	"code.google.com/p/protobuf/proto"

	"log"

	"github.com/robfig/cron"
)

type TcpManage struct {
	Conn        net.Conn
	NetBuf      *utils.NetBuffer
	Ip          string
	version     int
	errdescribe string
}

//创建一个tcp客户端
func CreateTcp() *TcpManage {
	this := &TcpManage{
		Ip:          "",
		version:     1,
		errdescribe: "",
		NetBuf:      utils.NewNetBuffer([]byte{69, 123, 132, 104, 67, 95, 33, 74, 120, 131, 61, 101, 55, 101, 69, 44}),
	}

	service := "192.168.0.213:16090"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Println("ResolveTCPAddr ", err.Error())
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println("DialTCP", err.Error())
	}
	this.Conn = conn
	return this
}

func (this *TcpManage) ReadTcp() {
	defer this.OnClose()
LoopConn:
	for {
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
				closeConn, resp := this.tcpNetworkRead(msg)
				if this.OnResp(closeConn, resp) {
					this.errdescribe = fmt.Sprintf("IP:[%s] socket Msgid [%d] return false ", this.Ip, utils.GetMsgCode(msg))
					break LoopConn
				}
			} else {
				break
			}
			time.Sleep(time.Millisecond * 1)
		}
	}
}
func (this *TcpManage) CronTask() {
	c := cron.New()
	spec := "*/3  * * * * *"
	c.AddFunc(spec, func() {
		resp := new(gt_msg.HHRequest)
		this.SendMsg(resp)
		log.Println(resp)
	})
	c.Start()
}

func (this *TcpManage) tcpNetworkRead(req proto.Message) (bool, proto.Message) {
	switch req.(type) {
	case *gt_msg.HHResponse:
		{
			return false, this.HH()
		}
	case *gt_msg.Order:
		{
			this.ReadOrder(req.(*gt_msg.Order))
			return false, nil
		}
	case *gt_msg.AppOrder:
		{
			return false, nil
		}
	default:
		{
			utils.SysLog.PutLineAsLog(fmt.Sprintf("Horn注册 Step err:UnReg %s", utils.ProtoToString(req)))
			return true, gt_msg.CommonErrorMsg("没有注册此消息", req)
		}
	}
	return true, nil
}
func (this *TcpManage) ReadOrder(req *gt_msg.Order) {
	fmt.Printf("\n收到订单:%+v \n", req)
}

func (this *TcpManage) OnResp(closeConn bool, resp proto.Message) bool {
	if resp != nil { //发响应回去
		if _, err := this.SendMsg(resp); err != nil {
			closeConn = true
		}
	}
	return closeConn
}

// SendMsg 发送msg给用户
func (this *TcpManage) SendMsg(resp proto.Message) (int, error) {
	if resp != nil {
		n, err := this.Conn.Write(this.toData(resp))
		if err != nil && utils.SysLog != nil {
			utils.SysLog.PutLineAsLog(fmt.Sprintf("SendMsg(%s) ERROR:%s", utils.ProtoToString(resp), err.Error()))
		} else {
			if utils.SysLog != nil {
				if utils.GetMsgCode(resp) != 2 {
					utils.SysLog.PutLineAsLog(fmt.Sprintf("SendMsg:OK %s CodeID=%d", resp.String(), utils.GetMsgCode(resp)))
				}
			}
		}
		return n, err
	}
	return 0, nil
}

//pkt转成加密的二进制数据
func (this *TcpManage) toData(pkt proto.Message) []byte {
	return utils.ToData(pkt, this.NetBuf.Key, this.version, 0)
}

// HH 心跳请求服务器
func (this *TcpManage) HH() proto.Message {
	resp := new(gt_msg.HHRequest)
	return resp
}

func (this *TcpManage) OnClose() {
	this.Conn.Close()
}
