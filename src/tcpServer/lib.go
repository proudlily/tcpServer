package platForm

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"time"
	"utils"
)

//返回到页面
type Regs struct {
	Err     int    `json:"err"`
	Err_msg string `json:"err_msg"`
}

//验证时间
func check_time(urlTime string) bool {
	Nowtime := time.Now().Unix()
	urlTimeInt, err := strconv.ParseInt(urlTime, 10, 64)
	if err != nil {
		utils.SysLog.PutLineAsLog(fmt.Sprintf("转换时间失败:url.Time %s:err%v", urlTime, err.Error()))
		return false
	}
	if urlTimeInt+1800 > Nowtime {
		return true
	}
	return false
}

//验证时间
func TimeCheck(urlTime string, w http.ResponseWriter) bool {
	resTime := check_time(urlTime)
	if resTime != true {
		resNews := &Regs{
			Err:     1,
			Err_msg: "时间戳失效",
		}
		utils.SysLog.PutLineAsLog(fmt.Sprintf("时间戳失效"))
		ResToHtml(resNews, w)
		return false
	}
	return true
}

//返回到页面的json格式
func ResToHtml(res *Regs, w http.ResponseWriter) {
	reply, err := json.Marshal(res)
	if err != nil {
		fmt.Println("格式化json失败", err)
		return
	}
	fmt.Fprintf(w, "%s", string(reply))
}

//验证sign
func CheckSign(funcName, sign string, urlFiled interface{}, w http.ResponseWriter) bool {
	res := signControl(funcName, urlFiled)
	log.Printf("CheckUser sign is %s,md5 is %s", sign, res)
	if res != sign {
		resNews := &Regs{
			Err:     1,
			Err_msg: "sign与本地不匹配",
		}
		ResToHtml(resNews, w)
		return false
	}
	return true
}

//盐
const key string = "manage%D^%U8"

//加密
func signControl(ctName string, urlFiled interface{}) string {
	fileds := ctName
	value := reflect.ValueOf(urlFiled)
	for i := 0; i < value.NumField(); i++ {
		filed := fmt.Sprintf("%v", value.Field(i))
		fileds = fileds + filed
	}
	fileds += key
	data := []byte(fileds)
	s := fmt.Sprintf("%x", md5.Sum(data))
	return s
}
