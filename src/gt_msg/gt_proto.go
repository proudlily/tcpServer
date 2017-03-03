package gt_msg

import (
	"fmt"
	"utils"

	"github.com/golang/protobuf/proto"
)

var MsgMinDelayMSec [utils.MaxMsgCode]int //毫秒单位

func CommonErrorMsg(info string, pkt proto.Message) *CommonError {
	msg := new(CommonError)
	msg.Code = int32(utils.GetMsgCode(pkt))
	if info != "" {
		msg.SzDescribeString = info
	} else {
		msg.SzDescribeString = utils.ProtoToString(pkt)
	}
	utils.SysLog.PutLineAsLog(fmt.Sprintf("CommonErrorMsg:%s", utils.ProtoToString(pkt)))
	return msg
}
func NewMsg(msgCode uint16) proto.Message {
	MsgMinDelayMSec[msgCode] = 1000
	//控制消息注册
	if msgCode >= 0 {
		return Control_Msg(msgCode)
	}
	return nil
}

func Control_Msg(msgCode uint16) proto.Message {
	switch msgCode {
	case 0:
		return new(CommonError)
	case 1:
		return new(HHRequest)
	case 2:
		return new(HHResponse)
	case 3:
		return new(Order)
	case 4:
		return new(AppOrder)
	default:
		return nil
	}
}
