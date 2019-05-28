package tool

import (
	"encoding/binary"
	"fmt"
	"net"
)

//字节转换成整形
func BytesToInt(b []byte) int {
	fmt.Println("mmm",b)
	return int(binary.BigEndian.Uint16(b))
	//return int(x)
}

func Unpack(conn net.Conn,buf []byte) ([]byte, error) {
	//  read head
	//if 0 != this.readTimeoutSec {
	//	this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.readTimeoutSec) * time.Second))
	//}
	headBuf := buf[:6]
	n, err := conn.Read(headBuf)
	fmt.Println("nnnn",n,err)
	if err != nil {
		return nil, err
	}
	//  check length
	packetLength := BytesToInt(headBuf[4:6])
	fmt.Println("packetLength",packetLength,headBuf)
	//if packetLength > this.maxReadBufferLength ||
	if	0 == packetLength {
		fmt.Println("长度0")
	}

	//  read body
	//if 0 != this.readTimeoutSec {
	//	this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.readTimeoutSec) * time.Second))
	//}

	bodyLength := packetLength
	bufbody := make([]byte,bodyLength)
	_, err = conn.Read(bufbody[:bodyLength])
	if err != nil {
		return nil, err
	}

	//  ok
	msg := make([]byte, bodyLength + 6)
	copy(msg, buf[:6])
	copy(msg[6:], bufbody[:bodyLength])
	return msg, nil
}
func CheckError(err error){
	fmt.Println(err)
}