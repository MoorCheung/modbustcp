package tool

import "C"
import (
	"encoding/binary"
	"fmt"
	"github.com/Luxurioust/excelize"
	_ "github.com/go-sql-driver/mysql"
	"github.com/imroc/biu"
	"github.com/jinzhu/gorm"
	"log"
	"net"
	"strconv"
	"time"
)

//字节转换成整形
func BytesToInta(b []byte) int {
	fmt.Println("mmm", b)
	return int(binary.BigEndian.Uint16(b))
	//return int(x)
}

type TcpConn struct {
	Conn      net.Conn
	//Address   byte
	Online    bool
	Time      time.Time
	DeviceNum string
	WriteChan chan []byte
	ReadChan chan []byte
	WriteContent []byte
	ReadContent []byte
}
type DeviceInfo struct {
	DeviceName string
	Online     bool
	DeviceNum  string
	Time       time.Time
}
func init() {
}
//func Unpackl(buf []byte) ([]byte, error) {
//	//00 00 06 abcdef
//	headBuf := buf[:6]
//	//n, err := C.Conn.Read(headBuf)
//	fmt.Println("nnnn", n, err)
//	if err != nil {
//		return nil, err
//	}
//	fmt.Println("headbuf:", headBuf)
//	//if C.DeviceNum == "" {
//	//	C.DeviceNum = string(headBuf)
//	//	dtu := &DTU{}
//	//	dtu.Id = "abcdef"
//	//	data := dtu.Find()
//	//
//	//	DTUS[dtu.Id] = data
//	//	fmt.Println("data:", data)
//	//	bufbody := make([]byte, 2)
//	//	_, err = C.Conn.Read(bufbody[:2])
//	//	intbuf, _ := strconv.ParseInt(string(bufbody), 16, 10)
//	//	plc := &PLC{}
//	//	plc.PLCAddress = string(bufbody)
//	//	plcs := plc.Select()
//	//	PLCS := map[uint16]*PLC{}
//	//	for _, v := range plcs {
//	//		i, _ := strconv.Atoi(v.PLCAddress)
//	//		PLCS[uint16(i)] = v
//	//	}
//	//	data.PLCS = PLCS
//	//	//fmt.Println(intbuf,byte(intbuf),uint16(intbuf))
//	//	//C.Address = byte(uint16(intbuf))
//	//	//fmt.Println(C.Address)
//	//	return []byte("ok"), nil
//	//} else {
//		//  check length
//		packetLength := BytesToInt(headBuf[4:6])
//		fmt.Println("packetLength", packetLength, headBuf)
//		//if packetLength > this.maxReadBufferLength ||
//		if 0 == packetLength {
//			fmt.Println("长度0")
//		}
//		//  read body
//		//if 0 != this.readTimeoutSec {
//		//	this.conn.SetReadDeadline(time.Now().Add(time.Duration(this.readTimeoutSec) * time.Second))
//		//}
//		bodyLength := packetLength
//		bufbody := make([]byte, bodyLength)
//		_, err = C.Conn.Read(bufbody[:bodyLength])
//		if err != nil {
//			return nil, err
//		}
//
//		//  ok
//		msg := make([]byte, bodyLength+6)
//		copy(msg, buf[:6])
//		copy(msg[6:], bufbody[:bodyLength])
//		return msg, nil
//	}
//
//}
func CheckError(err error) {
	fmt.Println(err)
}

/**
开关量1和0
*/
func KaiGuan(b []byte) []string {
	var arr = make([]string, len(b)*8)
	//fmt.Println(b)
	for k, v := range b {
		//formatInt := strconv.FormatInt(int64(v), 2)
		formatInt := biu.ToBinaryString(uint8(v))
		//fmt.Println(formatInt)
		for kk, vv := range formatInt {
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
func FuncDecode(b []byte) []int {
	var arr = []int{}
	for k, _ := range b {
		if k%2 == 1 {
			arr = append(arr, BytesToInt(b[k-1:k+1]))
		}
	}
	fmt.Println(arr)
	return arr
}
func ReadExcel(filename string, Sheet string) map[string]*Point {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println("错误", err)
		//os.Exit(1)
	}
	// Get all the rows in the Sheet1.
	var Sheet1 string
	if Sheet != "" {
		Sheet1 = Sheet
	} else {
		Sheet1 = "Sheet1"
	}
	rows := xlsx.GetRows(Sheet1)
	fmt.Println("rows", len(rows))
	var arr = make(map[string]*Point)

	for _, row := range rows[1:] {
		if len(row) > 0 {
			p := &Point{}
			for kkk, Cell := range row {
				if kkk == 0 {
					p.Name = Cell
				} else if kkk == 2 {
					if Cell[:2] == "读写" {
						p.Type = "RW"
					} else if Cell[:2] == "只读" {
						p.Type = "R"
					} else {
						p.Type = "RW"
					}
				} else if kkk == 3 {
				} else if kkk == 4 {
					sint, _ := strconv.Atoi(Cell)
					p.ModbusAddress = uint16(sint)
				} else if kkk == 5 {
					sint, _ := strconv.Atoi(Cell)
					p.FunRead = uint16(sint)
				} else if kkk == 6 {
					sint, _ := strconv.Atoi(Cell)
					p.FunWrite = uint16(sint)
				} else if kkk == 7 {
					p.Operation = Cell
				} else if kkk == 8 {
					p.FunType = Cell
				}
			}
			str := strconv.Itoa(int(p.FunRead)) + ":" + strconv.Itoa(int(p.ModbusAddress))
			arr[str] = p
			//sql := fmt.Sprintf(`insert into points(plc_id,name,type,address,funread,funwrite,operation,funtype,value,sort) values("%s")`,p.Name,p.)
			//GormDb.Model(p).Create(p)
		}
	}
	return arr
}

var GormDb *gorm.DB
var GormErr error

func init() {
	var configMap = map[string]string{
		"user":     "root",
		"password": "mmDongkaikjcx13579",
		"host":     "139.129.119.229",
		"port":     "3306",
		"database": "gongchang",
	}
	mysqlurl := configMap["user"] + ":" + configMap["password"] + "@tcp(" + configMap["host"] + ":" + configMap["port"] + ")/" + configMap["database"] + "?charset=utf8"
	//gorm_model
	GormDb, GormErr = gorm.Open("mysql", mysqlurl)
	if GormErr != nil {
		log.Fatal("gorm_model", GormErr)
	}
	// 全局禁用表名复数
	GormDb.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响

}
