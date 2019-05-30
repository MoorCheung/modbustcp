package tool

import (
	"encoding/binary"
	"fmt"
	"github.com/imroc/biu"
	"net"
	"strconv"
	"time"
)

//字节转换成整形
func BytesToInt(b []byte) int {
	fmt.Println("mmm",b)
	return int(binary.BigEndian.Uint16(b))
	//return int(x)
}
type TcpConn struct {
	Conn net.Conn
	Address byte
	Online bool
	Time time.Time
	DeviceNum string
}
type DeviceInfo struct {
	DeviceName string
	Online bool
	DeviceNum string
	Time time.Time
}
var DeviceList = make(map[string]*DeviceInfo)

func init(){
	info := &DeviceInfo{}
	info.DeviceNum = "abcedf"
	info.Online = false
	info.DeviceName = "设备一"
	DeviceList[info.DeviceNum] = info
	info1 := &DeviceInfo{}
	info1.DeviceNum = "abcdefh"
	info1.Online = false
	info1.DeviceName = "设备二"
	DeviceList[info1.DeviceNum] = info1


}
func Unpack(C *TcpConn,buf []byte) ([]byte, error) {
	//  read head
	//if 0 != this.readTimeoutSec {
	//	this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.readTimeoutSec) * time.Second))
	//}
	headBuf := buf[:6]
	n, err :=C.Conn.Read(headBuf)
	fmt.Println("nnnn",n,err)
	if err != nil {
		return nil, err
	}
	fmt.Println("headbuf:",headBuf)
	if C.DeviceNum == "" {
		C.DeviceNum = string(headBuf)
		bufbody := make([]byte,2)
		_, err = C.Conn.Read(bufbody[:2])
		intbuf,_ := strconv.ParseInt(string(bufbody), 16, 10)
		fmt.Println(intbuf,byte(intbuf),uint16(intbuf))
		C.Address = byte(uint16(intbuf))
		fmt.Println(C.Address)
		return []byte("ok"),nil
	}else{
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
		_, err = C.Conn.Read(bufbody[:bodyLength])
		if err != nil {
			return nil, err
		}

		//  ok
		msg := make([]byte, bodyLength + 6)
		copy(msg, buf[:6])
		copy(msg[6:], bufbody[:bodyLength])
		return msg, nil
	}

}
func CheckError(err error){
	fmt.Println(err)
}
/**
	开关量1和0
 */
func KaiGuan(b []byte) []string {
	var arr = make([]string,len(b) *8)
	fmt.Println(b)
	for k,v := range b {
		//formatInt := strconv.FormatInt(int64(v), 2)
		formatInt := biu.ToBinaryString(uint8(v))
		//fmt.Println(formatInt)
		for kk,vv := range formatInt {
			//fmt.Println("第几列",k,string(vv),formatInt)
			arr[8*(k+1)-1-kk] = string(vv)
		}
	}
	//fmt.Println(arr)
	//for k,v := range arr {
	//	fmt.Println("k",k,v)
	//}
	return arr
}
func FuncDecode(b []byte)[]int{
	var arr = []int{}
	for k,_ :=range  b{
		if k % 2 == 1 {

			arr = append(arr,BytesToInt(b[k-1:k+1]))
		}
	}
	fmt.Println(arr)
	return arr
}