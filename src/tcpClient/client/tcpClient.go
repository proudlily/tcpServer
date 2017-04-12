package main

import (
	"gt_msg"
	"tcpClient"
	"utils"
)

//模拟一个客户端
func main() {
	tcpSocker := tcpClient.CreateTcp()
	//定时发消息
	tcpSocker.CronTask()
	//接收消息
	tcpSocker.ReadTcp()
}
func init() {
	utils.InitProtoTool(gt_msg.NewMsg)
}
