package utils

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	XdomainReq = "<policy-file-request/>\x00"
	XdomainRes = "<?xml version=\"1.0\"?><cross-domain-policy><site-control permitted-cross-domain-policies=\"all\"/><allow-access-from domain=\"*\" to-ports=\"*\" secure=\"false\"/></cross-domain-policy>\x00\""
)

var version string

func SendUdpPacket(addr string, data []byte) (int, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return 0, err
	}

	udp, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return 0, err
	}
	defer udp.Close()

	return udp.Write(data)
}

func WaitForSignal() os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)

	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT) //等到上述信号才停止函数
	s := <-signalChan                                                                                            //返回该信号
	signal.Stop(signalChan)                                                                                      //这些信号再不发给signalChan管道
	return s
}

type ServerCtlCallback interface {
	BeforeDown()
	GetLoad() int
}

type ServerControl struct {
	Name         string
	ServerConn   net.Listener
	addr         string
	internalAddr string
	cmdChan      chan bool
	connections  []net.Conn
	callback     ServerCtlCallback
	lock         *sync.Mutex
}

func NewServerControl(name, addr, internalAddr string, cmdChan chan bool, listener net.Listener, callback ServerCtlCallback) *ServerControl {
	SysLog.PutLineAsLog(fmt.Sprintf("%s init: %s/%s", name, addr, internalAddr))
	return &ServerControl{
		Name:         name,
		ServerConn:   listener,
		addr:         addr,
		internalAddr: internalAddr,
		cmdChan:      cmdChan,
		callback:     callback,
		lock:         new(sync.Mutex),
	}
}

func (control *ServerControl) ConnAdded(c net.Conn) {
	control.lock.Lock()
	defer control.lock.Unlock()
	control.connections = append(control.connections, c)
}

func (control *ServerControl) ConnRemoved(c net.Conn) {
	control.lock.Lock()
	defer control.lock.Unlock()

	for i, c2 := range control.connections {
		if c2 == c {
			//这种从队列删除的办法真是有趣
			control.connections = append(control.connections[:i], control.connections[i+1:]...)
			break
		}
	}
}

func (control *ServerControl) heartbeat() {
	for _ = range time.Tick(time.Second * 5) { //5 s 一次心跳
		load := len(control.connections)
		if control.callback != nil {
			l := control.callback.GetLoad()
			if l >= 0 {
				load = l
			}
		}
		SysLog.PutLineAsLog(fmt.Sprintf("%s load: %d %s/%s", control.Name, load, control.addr, control.internalAddr))
	}
}

//启动以后阻塞，等待退出事件，进行收尾处理，serverControl专用
func (control *ServerControl) Run() {
	defer func() {
		if err := recover(); err != nil {
			SysLog.PutLineAsLog(fmt.Sprintf("Panic error (control *ServerControl) Run() %d", err))
		}
	}()
	go control.heartbeat()

	<-control.cmdChan //阻塞于此
	SysLog.PutLineAsLog("contrlo Closeing......")
	control.ServerConn.Close()
	for _, c := range control.connections {
		c.Close()
	}
	//	if control.callback != nil {
	//		control.callback.BeforeDown()
	//	}
	SysLog.PutLineAsLog(fmt.Sprintf("%s exit", control.Name))
	control.cmdChan <- false
}

func GetVersion() string {
	return version
}
