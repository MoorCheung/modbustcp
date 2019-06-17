package main

import (
	"ServerModbus/tool"
	"encoding/json"
	"fmt"
	"github.com/go-log/log"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)


type ReqData struct {
	Id int32
	Value int32
	PlcAddress string
}
var tcpWrite = make(chan []byte,10)
type RespData struct {
	Code int32
	Data []*tool.DeviceInfo
	KaiGuan []string
	FuncValue []int
	MapData *tool.DTU
	PLcData map[uint16]*tool.PLC
	DtuData *tool.DTU
	Points []*tool.Point
}
var Kaiguan []string
var FuncValue []int
func RouterList(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	//results := ModbusWrite(2,reqData)
	//tcpWrite <- results

	resp.Header().Set("Content-Type","application/json")
	var lists []*tool.DeviceInfo

	data := &RespData{}
	data.Code = 200
	data.Data = lists
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
func RouterResult(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	//results := ModbusWrite(2,reqData)
	//tcpWrite <- results
	resp.Header().Set("Content-Type","application/json")
	data := &RespData{}
	data.Code = 200
	data.KaiGuan = Kaiguan
	data.FuncValue = FuncValue
	//data.MapData = funmap2
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
//获取我自己dtu 下对应的设备
func RouterMydtu(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	//results := ModbusWrite(2,reqData)
	//tcpWrite <- results
	resp.Header().Set("Content-Type","application/json")
	data := &RespData{}
	data.Code = 200
	//通过设备id 获取dtu
	dtu := tool.DTUS["abcdef"]
	var plcs = []*tool.PLC{}
	for _,v := range dtu.PLCS {
		plcs = append(plcs,v)
	}
	data.PLcData = dtu.PLCS
	data.DtuData = dtu
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
//我的plc设备
func RouterMyplc(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	//results := ModbusWrite(2,reqData)
	//tcpWrite <- results
	resp.Header().Set("Content-Type","application/json")
	data := &RespData{}
	data.Code = 200
	//通过设备id 获取dtu
	dtu := tool.DTUS["abcdef"]
	//var plcs = []*tool.PLC{}
	//for _,v := range dtu.PLCS {
	//	plcs = append(plcs,v)
	//}
	i,_ := strconv.Atoi(reqData.PlcAddress)
	plc := dtu.PLCS[uint16(i)]
	var Points  []*tool.Point

	tool.GormDb.Model(Points).Order("sort asc").Find(&Points)
	var r = make(map[uint16]int)
	for _,v := range plc.Points {
		r[v.ModbusAddress] = v.Value
	}
	for _,v := range Points {
		//fmt.Println(v.ModbusAddress,r,r[v.ModbusAddress])
		v.Value = r[v.ModbusAddress]
	}
	data.Points = Points
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
//写入数据
func RouterWrite(resp http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(body)
	//defer  req.Body.Close()
	reqData := &ReqData{}
	json.Unmarshal(body,reqData)
	fmt.Println(reqData)
	results := ModbusWrite(2,reqData)
	//tool.DTUS["abcdef"].PLCS["2"].
	tcpWrite <- results
	resp.Header().Set("Content-Type","application/json")
	data := &RespData{}
	data.Code = 200
	byte,_ := json.Marshal(data)
	resp.Write(byte)
}
func Router(resp http.ResponseWriter, req *http.Request) {
	var b = []byte{196,255,1}
	Kaiguan = tool.KaiGuan(b)
	var c = []byte{00,01,00,03,00,64,00,100,00,32,00,20,00,10,00,40}
	FuncValue = tool.FuncDecode(c)
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
	//fmt.Println(data,uint16(data.Id))
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
	tool.InitPlc()
	http.HandleFunc("/test", Router)
	http.HandleFunc("/list", RouterList)
	http.HandleFunc("/result", RouterResult)
	http.HandleFunc("/mydtu", RouterMydtu)
	http.HandleFunc("/myplc", RouterMyplc)
	http.HandleFunc("/write", RouterWrite)
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

	}
}
func TcpWrite(C *tool.TcpConn){
	//通过通道进行写
	for {
		select {
		case results := <- tcpWrite :
			fmt.Println("发送数据",results)
			//results,err := server.ServerReadHoldingRegisters(2,1)
			C.Conn.Write(results)
		}
	}
}
func TickRead(C *tool.TcpConn){
	handler := &tool.TcpHandler{}
	if C.DeviceNum != "" {
		//dtu上线获取dtu下所有plc 循环对plc发送命令
		plcs := tool.DTUS[C.DeviceNum].PLCS
		for _,plc := range plcs {
			i,_ := strconv.Atoi(plc.PLCAddress)
			fmt.Println("plc设备地址",plc.PLCAddress)
			handler.SlaveId = byte(i)
			server := tool.NewServer(handler)
			//获取开关量
			results, err := server.ServerReadCoils(0, 18)
			results1, err := server.ServerReadHoldingRegisters(0, 8)
			fmt.Println(err)
			for {
				tcpWrite <- results
				time.Sleep(1*time.Second)
				tcpWrite <- results1
				time.Sleep(1*time.Second)
			}
		}

	}

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
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan []byte, 16)
	go tcpreader(C ,readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := C.Conn.Read(buffer)
		fmt.Println("粘包:",buffer[:n],string(buffer[:n]))
		if err != nil {
			log.Logf(C.Conn.RemoteAddr().String(), " connection error: ", err)
			return
		}

		tmpBuffer = tool.Unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
		fmt.Println("tmpBuffer",tmpBuffer)
	}
}
//读入消息并处理
func tcpreader(C *tool.TcpConn,readerChannel chan []byte) {
	for {
		select {
		case data := <-readerChannel:
			fmt.Println("data",data,string(data))
			if C.DeviceNum == "" {
				C.DeviceNum = string(data)[:6]
			}else{
				fmt.Println("data",string(data))

			}

		}
	}
}