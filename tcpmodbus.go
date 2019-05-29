package main

import (
	"ServerModbus/tool"
	"encoding/json"
	"fmt"
	"github.com/go-log/log"
	"io/ioutil"
	"net"
	"net/http"
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
type ReqData struct {
	Id int32
	Value int32
}
var tcpWrite = make(chan []byte,10)
type RespData struct {
	Code int32
	Data []*tool.DeviceInfo
}
func RouterList(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	//results := ModbusWrite(2,reqData)
	//tcpWrite <- results

	resp.Header().Set("Content-Type","application/json")
	infos := tool.DeviceList
	var lists []*tool.DeviceInfo
	for _,v := range infos {
		lists = append(lists,v)
	}
	data := &RespData{}
	data.Code = 200
	data.Data = lists
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
func Router(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//    r.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	results := ModbusWrite(2,reqData)
	tcpWrite <- results

	resp.Header().Set("Content-Type","application/json")
	type Data struct {
		Code int32
		Content string `json:"Content"`
		Name string `json:"name"`
		Email string `json:"email"`
	}
	data := Data{}
	data.Code = 200
	data.Content = "hello world"
	data.Name = "zhangsan"
	data.Email = "123@qq.com"
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
//tcp modbus 写入06 or 01
func ModbusWrite(SlaveId int,data *ReqData)[]byte{
	handler := &tool.TcpHandler{}
	handler.SlaveId = 2

	server := tool.NewServer(handler)
	fmt.Println(data,uint16(data.Id))
	var results []byte
	if data.Id > 17 {
		address := data.Id -18
		results, _ = server.ServerWriteSingleRegister(uint16(address),uint16(data.Value))
	}else{
		var err error
		var value uint16
		if data.Value == 1 {
			value = 0xFF00
		}else{
			value = 0x0000
		}
		results, err = server.ServerWriteSingleCoil(uint16(data.Id),value)
		fmt.Println(err)
	}
	//results, _ := server.ServerWriteSingleRegister(uint16(address),uint16(data.Value))
	return results
}

func  main(){
	http.HandleFunc("/test", Router)
	http.HandleFunc("/list", RouterList)
	go http.ListenAndServe("0.0.0.0:8088", nil)
	//初始化tcp
	netListen, err := net.Listen("tcp", ":9001")
	tool.CheckError(err)
	defer netListen.Close()
	log.Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		log.Log(conn.RemoteAddr().String(), " tcp connect success")
		tcpConn := &tool.TcpConn{}
		tcpConn.Conn = conn
		tcpConn.Time = time.Now()
		//协程读
		go TcphandleConnection(tcpConn)
		//协程写
		go TcpWrite(tcpConn)
		//go TesttcpWrite()
	}
}
func TcpWrite(C *tool.TcpConn){
	//通过通道进行写
	for {
		select {
		case results := <- tcpWrite :
			fmt.Println("发送数据",results)
			fmt.Sprintf("%x",results)
			//results,err := server.ServerReadHoldingRegisters(2,1)
			C.Conn.Write(results)
		}
	}
}
func TesttcpWrite(){
	//handler := &tool.TcpHandler{}
	//handler.SlaveId = 2
	//server := tool.NewServer(handler)
	//results, err := server.ServerReadHoldingRegisters(0, 8)
	//bytes, e := server.ServerWriteSingleRegister(2, 10)
	//bytes, e := server.ServerWriteMultipleRegisters(0,8,[]byte{0,1,0,2,0,80,0,100,0,60,0,1,0,2,0,3})
	//bytes, e := server.ServerWriteMultipleCoils(0,18,[]byte{0xff,0xff,3})
	//bytes, e := server.ServerReadCoils(0,18)
	//if err != nil {
	//	fmt.Println(err,e)
	//}else{
	//	for {
	//		fmt.Println("每5秒读一次")
	//		time.Sleep(5*time.Second)
	//		//tcpWrite <- results
	//		//tcpWrite <- bytes
	//	}
	//}
}

func TcphandleConnection(C *tool.TcpConn){
	for   {
		buf := make([]byte,2048)
		msg,err := tool.Unpack(C,buf)
		if err != nil {
			tool.DeviceList[C.DeviceNum].Online = false
			C.Conn.Close()
			break
		}else{
			if string(msg) == "ok" {
				C.Online = true
				tool.DeviceList[C.DeviceNum].Online = true
				fmt.Println("reg ok",C.Address,C.DeviceNum)
			}else{
				handler := tool.TcpHandler{}
				handler.SlaveId = C.Address
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
}
