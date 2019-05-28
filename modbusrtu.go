package main

import (
	"ServerModbus/tool"
	"fmt"
)

func main(){
	handler := tool.RTUClientHandler{}
	handler.SlaveId = 2
	pdu := tool.ProtocolDataUnit{
		FunctionCode: 3,
		Data:         []byte{0x00, 0x01},
	}
	fmt.Println(pdu)
	adu, err := handler.Encode(&pdu)
	fmt.Println(adu,err)
}