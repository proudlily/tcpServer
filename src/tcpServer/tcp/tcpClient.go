package main

import (
	"gt_msg"
	"tcpClient1"
	"utils"
)

//模拟一个客户端
func main() {
	tcpSocker := tcpClient1.CreateTcp()
	//定时发消息
	tcpSocker.CronTask()
	//接收消息
	tcpSocker.ReadTcp()

	//开始等待信号量
	///utils.SysLog.PutLineAsLog("hronServer wait-for-signal")
	//s := utils.WaitForSignal() //阻塞，等待对应的信号量 syscall.SIGHUP, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT
	//utils.SysLog.PutLineAsLog(fmt.Sprintf("hronServer signal got: %v", s))
}
func init() {
	utils.InitProtoTool(gt_msg.NewMsg)
}
