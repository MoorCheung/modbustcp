package main

import (
	"ServerModbus/tool"
	"fmt"
)

func main(){
	handler := &tool.RtuHandler{}
	handler.SlaveId = 2
	rtuserver := tool.NewServer(handler)
	results, err := rtuserver.ServerReadHoldingRegisters(2, 1)
	fmt.Println(results,err)
	tcpHandler := &tool.TcpHandler{}
	tcpserver := tool.NewServer(tcpHandler)
	bytes, e := tcpserver.ServerReadHoldingRegisters(2, 1)
	fmt.Println(bytes,e)
	//handler := &inter.ServerHandler{}
	//handler.SaveId = "1"
	////handler.Encode()
	////handler.Decode()
	//client := inter.NewClient(handler)
	//client.Add()
	//clientHandler := &inter.ClientHandler{}
	//clientHandler.SaveId = "2"
	////clientHandler.Encode()
	////clientHandler.Decode()
	//newClient := inter.NewClient(clientHandler)
	//newClient.Add()
}
