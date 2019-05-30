package main

import (
	"ServerModbus/tool"
	"encoding/binary"
	"fmt"
)
func dataBlock(value ...uint16) []byte {
	data := make([]byte, 2*len(value))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	return data
}
func dataBlockSuffix(suffix []byte, value ...uint16) []byte {
	length := 2 * len(value)
	data := make([]byte, length+1+len(suffix))
	for i, v := range value {
		binary.BigEndian.PutUint16(data[i*2:], v)
	}
	data[length] = uint8(len(suffix))
	copy(data[length+1:], suffix)
	return data
}
func main(){
	var b = []byte{04,00,02,00,03,00,15}
	var arr = []int{}
	for k,_ :=range b[1:] {
		if k % 2 == 1 {
			arr = append(arr,tool.BytesToInt(b[1:][k-1:k+1]))
		}

	}
	fmt.Println(arr)
	//a输出结果:00011110


}
