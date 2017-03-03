package platForm

import (
	"fmt"
	"net/http"
	"sync"
	"utils"
)

//苹果订单传的url的参数
type appOrderUrl struct {
	openUserID string //渠道id
	shareID    string //充值类型id
	orderID    string //订单ID
	ipAddress  string //ip地址
	cardPrice  string //card的价格

	cardGold string //card的金币
	gameID   string
	insertID string //商品id
	time     string //时间戳
	sign     string //sign验证
}

// 解析url
func parseAppOrderUrl(w http.ResponseWriter, r *http.Request) *appOrderUrl {
	params := r.URL.Query()
	return &appOrderUrl{
		openUserID: params.Get("openUserID"),
		shareID:    params.Get("shareID"),
		orderID:    params.Get("orderID"),
		ipAddress:  params.Get("ipAddress"),
		cardPrice:  params.Get("cardPrice"),

		cardGold: params.Get("cardGold"),
		gameID:   params.Get("gameID"),
		insertID: params.Get("insertID"),
		time:     params.Get("time"),
		sign:     params.Get("sign"),
	}
}

//苹果订单处理的函数
func AppOrderController(w http.ResponseWriter, r *http.Request) {
	lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	//解析url
	url := parseAppOrderUrl(w, r)
	utils.SysLog.PutLineAsLog(fmt.Sprintf("苹果充值,url is %v", url))
	//验证sign
	var urlFiled interface{} = appOrderUrl{
		url.openUserID,
		url.shareID,
		url.orderID,
		url.ipAddress,
		url.cardPrice,

		url.cardGold,
		url.gameID,
		url.insertID,
		url.time, ""}
	funcName := "FinishAppOrder"
	if !CheckSign(funcName, url.sign, urlFiled, w) {
		return
	}
	//验证url
	if !checkAppUrl(url, w) {
		return
	}
	//成功返回到页面
	resNews := &Regs{
		Err:     0,
		Err_msg: "",
	}
	ResToHtml(resNews, w)
}

//验证url的参数
func checkAppUrl(url *appOrderUrl, w http.ResponseWriter) bool {
	if url.openUserID == "" ||
		url.shareID == "" ||
		url.orderID == "" ||
		url.ipAddress == "" ||
		url.cardPrice == "" ||
		url.gameID == "" ||
		url.insertID == "" {
		resNews := &Regs{
			Err:     1,
			Err_msg: "参数输入不完整",
		}
		utils.SysLog.PutLineAsLog(fmt.Sprintf("%s", "参数输入不完整"))
		ResToHtml(resNews, w)
		return false
	}
	return true
}
