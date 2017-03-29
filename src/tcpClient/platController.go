package tcpClient

import (
	"fmt"
	"utils"
)

var MaxRegisterCnt = 100000

type HandleControl struct {
	tcpExample   *utils.SocketManage
	conns        *utils.SafeQueue
	orderDate    *utils.SafeQueue
	appOrderDate *utils.SafeQueue
}

func TcpMain(IP, Host string) {
	this := NewHandControl(IP)
	if this != nil {
		this.serve()
	}
}

func NewHandControl(IP string) *HandleControl {
	this := &HandleControl{
		conns:        utils.MakeNewSafeQueue(MaxRegisterCnt),
		orderDate:    utils.MakeNewSafeQueue(MaxRegisterCnt),
		appOrderDate: utils.MakeNewSafeQueue(MaxRegisterCnt),
	}
	this.tcpExample = utils.NewTCPCreate(IP, this)
	return this
}

func (this *HandleControl) serve() {
	this.tcpExample.TcpStart()
}

/*------------------以下3个函数是接口必须要实现--------------------------------------------------*/

//1：初始化入口类一些基本动作 （本函数不能处理耗时的东西 如有需求 go出去）
func (this *HandleControl) OnCallBackInit() {

}

//2：具体实现单个连接逻辑
func (this *HandleControl) OnTcpHCCallBackInit(tcphc *utils.TCPHandleConnection) {
	//建立单个conn
	MakeNewHanleExample(tcphc, this)
}

//3：IPC消息处理
func (this *HandleControl) OnMqResponse(MsgSender, MsgReceiver int32, msgHead string, msgData []byte) {

}

/*------------------以上3个函数是接口必须要实现--------------------------------------------------*/
func (this *HandleControl) Hornconnected(conn *HandleExample) {
	//保存单个tcp 连接
	this.conns.Push(conn)
	utils.SysLog.PutLineAsLog(fmt.Sprintf("存入conn %+v，size: %d，", conn, this.conns.Size()))
}

func (this *HandleControl) HorndisConnected(conn *HandleExample) {
	this.conns.Remove(conn)
	utils.SysLog.PutLineAsLog(fmt.Sprintf("移除 conn %+v,", conn))
}
