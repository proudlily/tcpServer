实现了tcp发送消息的客户端和服务端
==============

## gt_msg 
protobuf 生成的go的文件

## tcpServe
是tcp服务端

tcpServe是搭建了tcp服务器。</br>
tcp搭建起来很简单。可以参考Go网络编程。但是搭建好一个tcp服务器不简单。</br>
其中这个tcp实现了：
- 通过protobuf来传送消息
- 有SetReadDeadline，就是单位时间内没有tcp的连接没有读写，就断开。
- 搭建一个tcp的步骤中，发消息的时候，还需要加密和解密。

## tcpClient
是tcp客户端

tcpClient和tcpServe逻辑一样，读和写。
- 这其中的代码可以参考书:Go网络编程
   
## utils
是第三方包

- 实现日志的记录
- 链表

tcp就是以上这些啦。另外就是用语言来实现上面这些。</br>

## 讲解
### 1、 这里用到的数据结构除了golang原生的,还构造了链表数据结构

`utils/safeQueue.go`
```golang
type SafeQueue struct {
	list   *list.List
	maxNum int //最大容纳量
	lock   *sync.RWMutex   //加锁
}
相当于放进去的是interface,取出来的时候需要判断interface的类型
```
另外在utils包里面还有一种类似map的数据结构 `key:value`

`utils/safeStrMap.go`

```golang
type SafeStrMap struct {
	m    map[string]interface{}
	lock *sync.RWMutex //加锁
}
```
两种的用法区别是，当你需要只查询key的时候，来查数据，就用map.这样子代码比较方便理解。

-------
### 2、  讲解的是interface{}

根据定义: 
`interface类型定义了一组方法，如果某个对象实现了某个接口的所有方法，则此对象就实现了此接口`

这里有用到的是在
`utils/tcpServer/`

```golang
type TcpHandleConnectionEvent interface {
	//建立连接
	OnEventTCPNetworkBind(Remoteip string)
	//读取数据
	OnEventTCPNetworkRead(req proto.Message) (bool, proto.Message)
	//关闭网络连接
	OnEventTCPNetworkShut(err string)
}
```
相当于回调，调用有关实现了`TcpHandleConnectionEvent `接口的对象。

可以参考 `src/tcpServer/platConn.go`。

--------

### 3、 另外需要讲的是go的并发,在 `/src/tcpServer/tcpServer.go` TcpStart()函数。

一个client连接上来, 就分派一个goroutine给它。




另外大家觉得有帮助，可以请我喝杯cafe哦～

![image](cash.jpg)
