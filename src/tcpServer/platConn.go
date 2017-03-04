package platForm

import (
	"fmt"
	"gt_msg"
	"time"
	"utils"

	"github.com/golang/protobuf/proto"
)

type HandleExample struct {
	m_hc *utils.TCPHandleConnection
	Ctrl *HandleControl
}

func MakeNewHanleExample(tcphc *utils.TCPHandleConnection, handlecontrol *HandleControl) *HandleExample {
	this := &HandleExample{
		m_hc: tcphc,
		Ctrl: handlecontrol,
	}
	this.m_hc.HandleConnEvent = this
	return this
}

//-------------------------------------------------
//读取数据
func (this *HandleExample) OnEventTCPNetworkRead(req proto.Message) (bool, proto.Message) {
	switch req.(type) {
	case *gt_msg.HHRequest:
		{
			return false, this.HH()
		}
	default:
		{
			utils.SysLog.PutLineAsLog(fmt.Sprintf("Horn注册 Step err:UnReg %s", utils.ProtoToString(req)))
			return true, gt_msg.CommonErrorMsg("没有注册此消息", req)
		}
	}
	return true, nil
}

//关闭网络连接
func (this *HandleExample) OnEventTCPNetworkShut(err string) {
	fmt.Println(fmt.Sprintf("onClose--%s", err))
	this.Ctrl.HorndisConnected(this)
}

//建立连接
func (this *HandleExample) OnEventTCPNetworkBind(Remoteip string) {
	//在conn的管理器保存这个HandleExample
	this.Ctrl.Hornconnected(this)
	fmt.Println(fmt.Sprintf("%s connection", Remoteip))
}

//---------------------------------------------

// HH 心跳返回服务器时间戳
func (this *HandleExample) HH() proto.Message {
	resp := new(gt_msg.HHResponse)
	resp.ServerTimeNow = time.Now().Unix()
	return resp
}

//发送信息给conn
func (this *HandleExample) SendMsg(resp proto.Message) {
	this.m_hc.SendMsg(resp)
}
