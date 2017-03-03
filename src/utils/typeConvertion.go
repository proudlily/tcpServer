package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	"code.google.com/p/mahonia"

	log "github.com/cihub/seelog"
)

func S2i(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Debugf("s2i(%s) err:%s", s, err.Error())
		return 0
	}
	return int(i)
}

func S2iD(s string, defaultv int) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		log.Debugf("s2i(%s) err:%s", s, err.Error())
		return defaultv
	}
	return int(i)
}

func S2i64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		//log.Debugf("s2i(%s) err:%s", s, err.Error())
		return 0
	}
	return i
}

func S2f64(s string) float64 {
	i, err := strconv.ParseFloat(s, 10)
	if err != nil {
		return 0
	}
	return i
}

func S2f32(s string) float32 {
	i, err := strconv.ParseFloat(s, 10)
	if err != nil {
		return 0
	}
	return float32(i)
}

func I2s(i int) string {
	return fmt.Sprintf("%d", i)
}

func I32s(i uint32) string {
	return fmt.Sprintf("%d", i)
}

func I2s64(i int64) string {
	return fmt.Sprintf("%d", i)
}

func Ints2s(i []int, sep string) string {
	str := ""
	for _, j := range i {
		if str != "" {
			str += sep
		}
		str += I2s(j)
	}
	return str
}

func S2ints(str, sep string) []int {
	if str == "" {
		return nil
	}
	//解析数据
	nums := []int{}
	strs := strings.Split(str, sep)
	for _, j := range strs {
		nums = append(nums, S2i(j))
	}
	return nums
}

func Ints(ints32 []int32) []int {
	result := make([]int, len(ints32))
	for i, n := range ints32 {
		result[i] = int(n)
	}
	return result
}

func Ints32(ints []int) []int32 {
	result := make([]int32, len(ints))
	for i, n := range ints {
		result[i] = int32(n)
	}
	return result
}

func Ints64(ints []int64) []int64 {
	result := make([]int64, len(ints))
	for i, n := range ints {
		result[i] = int64(n)
	}
	return result
}

//string -> []byte
func Str2Byte(str string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(str)
	return buf.Bytes()
}

func SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}
	if begin >= lth {
		begin = lth
	}
	end := begin + length
	if end > lth {
		end = lth
	}
	// 返回子串
	return string(rs[begin:end])
}

func IntsCount(ints []int, n int) int {
	count := 0
	for _, d := range ints {
		if d == n {
			count++
		}
	}
	return count
}

func IntsPos(ints []int, n int) int {
	for i, d := range ints {
		if d == n {
			return i
		}
	}
	return -1
}

func IntsSum(ints []int) int {
	sum := 0
	for _, d := range ints {
		sum += d
	}
	return sum
}

func StringPos(strs []string, dst string) int {
	for i, s := range strs {
		if s == dst {
			return i
		}
	}
	return -1
}

func Utf82GBK(text string) string {
	var enc mahonia.Encoder
	enc = mahonia.NewEncoder("gbk")
	if ret, ok := enc.ConvertStringOK(text); ok {
		return ret
	}
	return ""
}

func Byte2Base10(b []byte) (n uint64, err error) {
	base := uint64(10)
	n = 0
	for i := 0; i < len(b); i++ {
		var v byte
		d := b[i]
		switch {
		case '0' <= d && d <= '9':
			v = d - '0'
		default:
			n = 0
			err = errors.New("Base10err")
			return
		}
		n *= base
		n += uint64(v)
	}
	return n, err
}

func IntsDelInts(sints, dints []int) []int {
	dels := []int{}
	offindex := 0
	for _, i := range dints {
		for index, c := range sints {
			if c == i {
				if IntsPos(dels, index) == -1 {
					dels = append(dels, index)
				}
			}
		}
	}
	sort.Ints(dels)
	for _, i := range dels {
		index := i - offindex
		sints = append(sints[:index], sints[index+1:]...)
		offindex++
	}
	return append([]int{}, sints...)
}

func Trace(msg string) func() {
	start := time.Now()
	SysLog.PutLineAsLog(fmt.Sprintf("enter %s\n", msg))

	return func() {
		SysLog.PutLineAsLog(fmt.Sprintf("exit %s (%s)\n", msg, time.Since(start)))
	}
}

func UTF8to16(in string) []uint16 {

	var Out = []uint16{}
	uti_str := []rune(in)

	for _, j := range uti_str {
		Out = append(Out, uint16(j))
	}
	return Out
}

func UTF16to8(in []uint16) string {

	var out []byte
	var codepoint uint
	var charCount int
	var retString string
	count := len(in)

	for charCount = 0; charCount < count; charCount++ {

		if in[charCount] >= 0xd800 && in[charCount] <= 0xdbff {
			codepoint = uint((in[charCount]-0xd800)<<10) + 0x10000
		} else {
			if in[charCount] >= 0xdc00 && in[charCount] <= 0xdfff {
				codepoint |= uint(in[charCount] - 0xdc00)
			} else {
				codepoint = uint(in[charCount])
			}

			if codepoint <= 0x7f {
				out = append(out, byte(codepoint))
			} else if codepoint <= 0x7ff {
				out = append(out, byte(0xc0|((codepoint>>6)&0x1f)))
				out = append(out, byte(0x80|(codepoint&0x3f)))
			} else if codepoint <= 0xffff {
				out = append(out, byte(0xe0|((codepoint>>12)&0x0f)))
				out = append(out, byte(0x80|((codepoint>>6)&0x3f)))
				out = append(out, byte(0x80|(codepoint&0x3f)))
			} else {
				out = append(out, byte(0xf0|((codepoint>>18)&0x07)))
				out = append(out, byte(0x80|((codepoint>>12)&0x3f)))
				out = append(out, byte(0x80|((codepoint>>6)&0x3f)))
				out = append(out, byte(0x80|(codepoint&0x3f)))
			}
			codepoint = 0
		}
	}

	var index int = 0
	for i := 0; i < len(out); i++ {
		if out[i] != 0 {
			index++
		}
	}
	//fmt.Println("index", index)
	//fmt.Println("out...", out)
	retString = string(out[0:index])
	return retString
}

func Inet_ntoa(ipnr uint32) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3])
}
