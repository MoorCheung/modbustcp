//通讯协议处理，主要处理封包和解包的过程
package tool

import (
	"bytes"
	"encoding/binary"
	"fmt"
)
const (
	//ConstHeader        =   "0000"
	ConstHeaderLength          = 4
	ConstSaveDataLength        = 2
)

var ConstHeader = []byte{0x00,0x00,0x00,0x00}
var ConstRegHeader = []byte{0x00,0x00,0x10,0x10}
//注册包
func PacketReg(message []byte) []byte {
	return append(append([]byte(ConstRegHeader), IntToBytes(len(message))...), message...)
}
//封包
func Packet(message []byte) []byte {
	return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//解包
func Unpack(buffer []byte, readerChannel chan []byte) []byte {
	fmt.Println("解包",buffer)
	length := len(buffer)

	var i int
	for i = 0; i < length; i = i + 1 {
		if length < i+ConstHeaderLength+ConstSaveDataLength {
			break
		}
		if string(buffer[i:i+ConstHeaderLength]) == string(ConstHeader) {
			messageLength := BytesToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstSaveDataLength])
			fmt.Println("messageLength",messageLength)
			if length < i+ConstHeaderLength+ConstSaveDataLength+messageLength {
				break
			}
			data := buffer[i+ConstHeaderLength+ConstSaveDataLength : i+ConstHeaderLength+ConstSaveDataLength+messageLength]
			readerChannel <- data

			i += ConstHeaderLength + ConstSaveDataLength + messageLength - 1
		}
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := uint16(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint16
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}