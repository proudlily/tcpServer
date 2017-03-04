package platForm

import (
	"fmt"
	"gt_msg"
	"net/http"
	"sync"
	"utils"
)

//完成订单传的url
type finishOrderUrl struct {
	ordersID    string //订单ID
	orderAmount string //订单金额
	sign        string //sign 验证
}

var lock *sync.Mutex

//完成订单的函数处理
func (this *HandleControl) FinishOrder(w http.ResponseWriter, r *http.Request) {
	lock = new(sync.Mutex)
	lock.Lock()
	defer lock.Unlock()
	//解析url
	url := parseFinishOrder(r)
	utils.SysLog.PutLineAsLog(fmt.Sprintf("完成订单,url is %+v", url))

	//验证sign
	var urlFiled interface{} = finishOrderUrl{
		url.ordersID,
		url.orderAmount,
		""}
	funcName := "FinishOrder"
	if !CheckSign(funcName, url.sign, urlFiled, w) {
		return
	}
	//验证url
	if !checkOrderUrl(url, w) {
		return
	}
	//存数据
	orderMeta := new(gt_msg.Order)
	orderMeta.OrderAmount = url.orderAmount
	orderMeta.OrderID = url.ordersID
	this.orderDate.Push(orderMeta)
	//---
	//发数据
	//在线
	if this.conns.Size() != 0 {
		this.orderDate.EachItem(func(e interface{}) {
			if mess, ok := e.(*gt_msg.Order); ok && mess != nil {
				this.conns.EachItem(func(e interface{}) {
					if ret_m_hc, ok := e.(*HandleExample); ok && ret_m_hc != nil {
						ret_m_hc.SendMsg(mess)
					}
				})
			}
		})
		//将队列清零
		this.orderDate.Clear()
	}
	//成功返回到页面
	resNews := &Regs{
		Err:     0,
		Err_msg: "",
	}
	ResToHtml(resNews, w)
}

//解析url的参数
func parseFinishOrder(r *http.Request) *finishOrderUrl {
	params := r.URL.Query()
	return &finishOrderUrl{
		ordersID:    params.Get("order_id"),
		orderAmount: params.Get("order_amount"),
		sign:        params.Get("sign"),
	}
}

//验证订单的参数的合法性
func checkOrderUrl(url *finishOrderUrl, w http.ResponseWriter) bool {
	if url.ordersID == "" ||
		url.orderAmount == "" ||
		url.sign == "" {
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
