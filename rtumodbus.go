package main

import (
	"ServerModbus/tool"
	"fmt"
	"github.com/go-log/log"
	"net"
	"time"
)
var funmap1 = map[int]string{
	0:"屏启动",
	1:"屏停止",
	2:"屏报警消音",
	3:"运行标识",
	4:"点火完成",
	5:"温度上限警告",
	6:"点火报警1",
	7:"点火报警2",
	8:"火灭报警",
	9:"风机反馈",
	10:"风机反馈",
	11:"火焰反馈",
	12:"变频故障",
	13:"点火反馈",
	14:"报警器",
	15:"开阀",
	16:"变频启停",
	17:"点火线圈",

}
var funmap3 = map[int]string{
	0:"吹扫延时设定",
	1:"点火延时设定",
	2:"设定温度",
	3:"温度上限",
	4:"实际温度",
	5:"调节阀开度输出",
	6:"风机吹扫延时",
	7:"点火延时",
}

var rtuWrite = make(chan []byte,10)
func  main(){
	//初始化tcp
	netListen, err := net.Listen("tcp", ":9002")
	tool.CheckError(err)
	defer netListen.Close()
	log.Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log.Log(conn.RemoteAddr().String(), " tcp connect success")
		//协程读
		go RtuhandleConnection(conn)
		//协程写
		go RtuWrite(conn)
		go TestrtuWrite()
	}
}
func RtuWrite(conn net.Conn){
	//通过通道进行写

	for {
		select {
		case results := <- rtuWrite :
			//results,err := server.ServerReadHoldingRegisters(2,1)
			conn.Write(results)
		}
	}
}
func TestrtuWrite(){
	handler := &tool.RtuHandler{}
	handler.SlaveId = 2
	server := tool.NewServer(handler)
	results, err := server.ServerReadHoldingRegisters(0, 8)
	bytes, e := server.ServerWriteSingleRegister(2, 100)
	if err != nil {
		fmt.Println(err,e)
	}else{
		for {
			fmt.Println("每5秒读一次")
			time.Sleep(5*time.Second)
			rtuWrite <- results
			time.Sleep(1*time.Second)
			rtuWrite <- bytes
		}
	}
}

func RtuhandleConnection(conn net.Conn){
	for   {
		buf := make([]byte,2048)
		msg,err := tool.Unpack(conn,buf)
		if err != nil {
			conn.Close()
			break
		}else{
			handler := tool.TcpHandler{}
			handler.SlaveId = 2
			pdu,err := handler.Decode(msg)
			fmt.Println("pdu",pdu,err)
			switch pdu.FunctionCode {
			case 3:
				fmt.Println("接收功能码03",time.Now(),pdu.Data)
				for k,v := range pdu.Data{
					fmt.Println(k,v)
					name := funmap3[k]
					fmt.Println(name,v)
				}
			case 1:
				fmt.Println("01功能码")
				fmt.Println("接收功能码03",time.Now(),pdu.Data)
				for k,v := range pdu.Data{
					fmt.Println(k,v)
					name := funmap1[k]
					fmt.Println(name,v)
				}
			}

		}

	}
}
