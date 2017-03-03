实现了tcp发送消息的客户端和服务端
==============

## gt_msg 是protobuf 生成的go的文件
- protobuf 是消息传送方式

## tcpServe是tcp服务端
tcpServe是搭建了一个http服务器和tcp服务器。</br>
把通过http的url传来的参数，再通过tcp服务传给tcp的客户端。</br>
tcp搭建起来很简单。可以参考Go网络编程。但是搭建好一个tcp服务器不简单。</br>
其中这个tcp实现了：
- 通过protobuf来传送消息
- 有SetReadDeadline，就是单位时间内没有tcp的连接没有读写，就断开。
- 搭建一个tcp的步骤中，发消息的时候，还需要加密和解密。

## tcpClient是tcp客户端
tcpClient和tcpServe逻辑一样，读和写。
- 这其中的代码可以参考书:Go网络编程
   
## utils是第三方包
- 实现日志的记录
- 链表

tcp就是以上这些啦。另外就是用语言来实现上面这些。</br>
如果大家学习了go的基础，了解了go的数据结构，条件语句。以及go特有的slice,</br>
interface ,结构体和方法。就可以看懂了这些代码了。:-)</br>

另外大家觉得有帮助，可以请我喝杯cafe哦～
![image](cash.jpg)
