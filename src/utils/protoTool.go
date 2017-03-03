package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

const (
	MaxMsgLen  = 65536
	maxMsgLen  = 4096 * 4
	MaxMsgCode = 10240 //WxrAdd_con 9000
	magicLen   = 4
	MinVer     = 1
	MaxVer     = 1
	modName    = "[MSG]"
	IsWindows  = false
)

var (
	ServerWinNo int64 = 1
	msgCodes          = map[reflect.Type]uint16{} //用来保存每种类型对应的编号
	verMagics         = map[int]string{}
	newMsg      func(n uint16) proto.Message
)

type NetBuffer struct {
	Key    []byte
	buf    [MaxMsgLen * 2]byte
	bufLen int
}

func InitProtoTool(f func(n uint16) proto.Message) {
	if f == nil {
		return
	}
	var i uint16
	newMsg = f
	for i = 0; i < MaxMsgCode; i++ {
		if msg := newMsg(i); msg != nil {
			msgCodes[reflect.TypeOf(msg)] = i
		}
	}

	//初始化versionedMagics map
	vers := []rune{} //   rune数组是为了转换成unicode编码
	for i := '1'; i <= '9'; i++ {
		vers = append(vers, i)
	}
	for i := MinVer; i <= MaxVer; i++ {
		verMagics[i] = fmt.Sprintf("GTV%c", vers[i-1])
	}
}

func GetMsgCode(msg proto.Message) int {
	msgCode, ok := msgCodes[reflect.TypeOf(msg)]
	if !ok {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg unregistered: %s", modName, reflect.TypeOf(msg)))
		return -1
	}
	return int(msgCode)
}

//msg -> []byte  封包 format:     GTV[VERSION][uint64 cmdnumber][uint16 PktCode][uint16 loadLen]data
func ToData(msg proto.Message, key []byte, ver int, cmdNo int64) []byte {
	if ver < MinVer || ver > MaxVer {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg version %d err not in [%d, %d]", modName, ver, MinVer, MaxVer))
		return nil
	}
	//proto.message -> []byte
	data, err := proto.Marshal(msg)
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg marshal %+v error: %s", modName, ProtoToString(msg), err.Error()))
		return nil
	}
	codeData(data, key)
	//msgCode
	msgCode, ok := msgCodes[reflect.TypeOf(msg)]
	if !ok {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg unregistered: %s", modName, reflect.TypeOf(msg)))
		return nil
	}
	buf := new(bytes.Buffer)
	if n, err := buf.WriteString(verMagics[ver]); err != nil || n != len(verMagics[ver]) {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg write magic %s error: %d, %v", modName, verMagics[ver], n, err))
		return nil
	}

	var Currnttime int64

	if IsWindows {
		// 额外增加包序号
		if cmdNo == 0 { //发送给客户端的
			ServerWinNo++
			Currnttime = ServerWinNo
		} else {
			Currnttime = cmdNo //转发的保留客户端序号
		}
	} else {
		// 额外增加包序号
		if cmdNo == 0 { //发送给客户端的
			Currnttime = time.Now().UnixNano()
		} else {
			Currnttime = cmdNo //转发的保留客户端序号
		}
	}

	//SysLog.PutLineAsLog(fmt.Sprintf(" Currnttime = %d", Currnttime))

	if err = binary.Write(buf, binary.BigEndian, Currnttime); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf(" msg binary write ServercmdNo error"))
		return nil
	}

	//fmt.Println("fs timer", Currnttime)

	//int -> []byte
	if err = binary.Write(buf, binary.BigEndian, msgCode); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg binary write pktCode error: %s", modName, err.Error()))
		return nil
	}

	dataLen := uint16(len(data))
	if err = binary.Write(buf, binary.BigEndian, dataLen); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg binary write dataLen error: %s", modName, err.Error()))
		return nil
	}

	if n, err := buf.Write(data); n != len(data) || err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s msg write data error: %d, %v", modName, n, err))
		return nil
	}
	return buf.Bytes()
}

//[]byte -> msg   解包        {head format：GTV[VERSION][uint64 cmdnumber][uint16 PktCode][uint16 loadLen]}
func ToMsg(data []byte, key []byte) (int, proto.Message, int, int64) {
	if len(data) < magicLen+12 {
		if len(data) > 0 { //avoid writing unnecessary logs here
			SysLog.PutLineAsLog(fmt.Sprintf("%s data less than header: %d", modName, len(data)))
		}
		return MaxVer, nil, 0, 0
	}

	s1 := string(data[:magicLen])
	ver := MinVer - 1
	for v, magic := range verMagics {
		if magic == s1 {
			ver = v
			break
		}
	}

	buf := bytes.NewReader(data[magicLen:])

	//读取客户端的序号
	var msgClinetNumber int64
	if err := binary.Read(buf, binary.BigEndian, &msgClinetNumber); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s error reading binary msgClinetNumber: %s", modName, err.Error()))
		return MinVer - 1, nil, 0, 0
	}
	//fmt.Println("msgClinetNumber==", msgClinetNumber)

	//读取pktCode两字节
	var msgCode uint16
	if err := binary.Read(buf, binary.BigEndian, &msgCode); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s error reading binary msgCode: %s", modName, err.Error()))
		return MinVer - 1, nil, 0, 0
	}

	//读取dataLen两字节
	var dataLen uint16
	if err := binary.Read(buf, binary.BigEndian, &dataLen); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s error reading binary dataLen: %s", modName, err.Error()))
		return MinVer - 1, nil, 0, 0
	}

	//超长
	if dataLen > maxMsgLen {
		SysLog.PutLineAsLog(fmt.Sprintf("%s dataLen too large: %d", modName, dataLen))
		return MinVer - 1, nil, 0, 0
	}
	//用于截取的整包长度
	used := magicLen + 4 + 8 + int(dataLen)
	if len(data) < used {
		SysLog.PutLineAsLog(fmt.Sprintf("%s data too small: %d/%d/%d", modName, len(data), used, dataLen))
		return ver, nil, 0, 0
	}

	msg := newMsg(msgCode)
	if msg == nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s unregistered pktCode: %d", modName, msgCode))
		return MinVer - 1, nil, used, 0
	}
	//解码
	raw := data[magicLen+4+8 : used]
	codeData(raw, key)

	//proto封包
	if err := proto.Unmarshal(raw, msg); err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("%s Unmarshal error: %s", modName, err.Error()))
		SysLog.PutLineAsLog(string(raw))
		SysLog.PutHexAsLog(raw, len(raw))
		return MinVer - 1, nil, used, 0
	}
	return ver, msg, used, msgClinetNumber
}

//用key对data异或计算
func codeData(data []byte, key []byte) {
	l := len(key)
	if l == 0 {
		return
	}
	for i, n := range data {
		data[i] = n ^ key[i%l]
	}
}

//String
func ProtoToString(msg proto.Message) string {
	if msg == nil {
		return ""
	}
	strs := strings.Split(reflect.TypeOf(msg).String(), ".")
	str := strs[len(strs)-1]
	return str + " {" + strings.TrimSpace(msg.String()) + "}"
}

//Clone
// func Clone(msg proto.Message) proto.Message {
// 	if msg == nil {
// 		return nil
// 	}
// 	data := ToData(msg, nil, MinVer)
// 	ver, newMsg, used := ToMsg(data, nil)
// 	if ver != MinVer || used != len(data) {
// 		return nil
// 	}
// 	return newMsg
// }

//
func NewNetBuffer(key []byte) *NetBuffer {
	return &NetBuffer{
		Key:    key,
		bufLen: 0,
	}
}

//
func (this *NetBuffer) GetAMsg() (int, proto.Message, int64) {
	ver, msg, used, cmdNo := ToMsg(this.buf[:this.bufLen], this.Key)
	if used > 0 && used <= this.bufLen {
		if this.bufLen > used {
			copy(this.buf[0:], this.buf[used:])
		}
		this.bufLen -= used
	}
	return ver, msg, cmdNo
}

//代理直接拷贝
func (this *NetBuffer) CompleteData() []byte {
	defer func() {
		this.bufLen = 0
	}()
	return this.buf[:this.bufLen]
}

//
func (this *NetBuffer) DelFromHead(buf []byte) bool {
	if len(buf) == 0 || len(buf) > this.bufLen {
		return false
	}
	if bytes.Compare(buf, this.buf[:len(buf)]) == 0 {
		copy(this.buf[0:], this.buf[len(buf):])
		this.bufLen -= len(buf)
		return true
	}
	return false
}

//
func (this *NetBuffer) Read(reader io.Reader) error {
	n, err := reader.Read(this.buf[this.bufLen:])
	if err == nil && n >= 0 {
		this.bufLen += n
		return nil
	}
	if err != nil {
		SysLog.PutLineAsLog(fmt.Sprintf("NetBuffer.Read ERR:%s", err.Error()))
	}
	return err
}
