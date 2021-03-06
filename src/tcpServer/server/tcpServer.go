package main

import (
	"flag"
	"gt_msg"
	"tcpServer"

	"runtime"
	"utils"
)

var AppName string = "platForm"
var pcfgpath *utils.Config

var (
	IP string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//启动
	tcpServer.TcpMain(IP)
}

func init() {
	cfgfile := flag.String("cfg", "server.ini", "config file for platForm")
	flag.Parse()
	//读取配置
	NewDefenseConfig(*cfgfile)
	//初始化日志
	if InitDefenseSysLog() == false {
		utils.SysLog.PutLineAsLog("创建日志失败")
		return
	}
	utils.SysLog.PutLineAsLog("创建日志成功")
	//注册proto
	utils.InitProtoTool(gt_msg.NewMsg)
}

//初始化日志
func InitDefenseSysLog() bool {
	if utils.SysLog == nil {
		utils.SysLog = utils.MakeNewMyLog(AppName+"logs", AppName+"_sys.log", 50000000, 15)
	}
	if utils.SysLog == nil {
		return false
	} else {
		return true
	}
}

func NewDefenseConfig(configGlobalPath string) {
	pcfgpath = utils.SetConfig(configGlobalPath)
	IP = pcfgpath.GetValueString("TcpIp", "Ip")
}
