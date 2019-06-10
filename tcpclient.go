package main

import "C"
import (
	"ServerModbus/tool"
	"fmt"
	"net"
	"net/http"
)

func main(){
	fmt.Println("tcp client")
	go TcpConn()
	http.ListenAndServe("0.0.0.0:8082", nil)
	
}
func TcpConn(){
	tcpaddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:9001")
	if err != nil {
		fmt.Println(err)
	}else{
		tcpconn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			fmt.Println(err)
		}else{
			str := "abcdef02"
			tcpconn.Write([]byte(str))
			go TcpRead(tcpconn)
			go Tick()
		}
	}

}
var TcpReadChan = make(chan []byte,1000)
func TcpRead(c *net.TCPConn){
	for {
		buf := make([]byte,2048)
		headBuf := buf[:6]
		n, err := c.Read(headBuf)
		if err != nil {
			fmt.Println(err,n)
		}
		fmt.Println("n",n)
		packetLength := tool.BytesToInt(headBuf[4:6])
		fmt.Println("packetLength",packetLength,headBuf)
		//if packetLength > this.maxReadBufferLength ||
		if	0 == packetLength {
			fmt.Println("长度0")
		}
		bodyLength := packetLength
		bufbody := make([]byte,bodyLength)
		_, err = c.Read(bufbody[:bodyLength])
		if err != nil {

		}
		//  ok
		msg := make([]byte, bodyLength + 6)
		copy(msg, buf[:6])
		copy(msg[6:], bufbody[:bodyLength])
		fmt.Println(msg)
		TcpReadChan <- msg
	}
}

func Tick(){
	for {
		select {
		case bytes,ok := <- TcpReadChan:
			if ok {
				fmt.Println("接收消息",string(bytes))

			}

		}
	}
}