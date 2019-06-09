package tool

import (
	"fmt"
	"strconv"
)

//plc设备
type PLC struct {
	PlCName string //plc设备名称
	PLCAddress string //plc设备地址
	Online bool //plc是否在线
	Points map[string]*Point
}
type Point struct {
	Name string `gorm:"column:name"` //名称
	Type string `gorm:"column:type"`//类型
	ModbusAddress uint16 `gorm:"column:address"`//MODBUS寄存区地址
	FunRead uint16 `gorm:"column:funread"`//功能码(读）
	FunWrite uint16 `gorm:"column:funwrite"`//功能码(写）
	Operation string `gorm:"column:operation"`//位操作 16位无符号
	FunType string `gorm:"column:funtype"`//功能类型需要在屏上显示
	Value int `gorm:"column:value"`//当前值
	Sort int `gorm:"column:sort"`//排序
}
//DTU设备
type DTU struct {
	Id string `gorm:"column:dtu_id"`  //设备id
	Pwd string `gorm:"column:dtu_pwd"` //密码
	Name string `gorm:"column:dtu_name"` //名称
	Mobile string `gorm:"column:mobile"`  //手机号
	PLCS map[uint16]*PLC //pcl数组
	Online bool `gorm:"column:online"` //在线状态
	Time int64 `gorm:"column:time"` //心跳时间 在线时间
	Status int64 `gorm:"column:status"` //心跳时间 在线时间
}
func (d *DTU)Find()*DTU{
	GormDb.Model(d).Find(d)
	return d
}
func (p *PLC)Find()*PLC{
	GormDb.Model(p).Find(p)
	return p
}
func (p *PLC)Select()[]*PLC{
	plcs := []*PLC{}
	GormDb.Model(p).Select(plcs)
	return plcs
}
var DTUS = make(map[string]*DTU)
//初始化plc数据
func InitPlc(){
	PlcPoints := ReadExcel("/Users/dongkai/go/src/ServerModbus/tool/modbus.xlsx", "Sheet1")
	plc := &PLC{}
	plc.PlCName = "plc1"
	plc.PLCAddress = "02"
	plc.Points = PlcPoints
	dtu := &DTU{}
	dtu.Id = "abcdef"
	plcs := make(map[uint16]*PLC)
	i,_ := strconv.Atoi(plc.PLCAddress)
	plcs[uint16(i)] = plc
	dtu.PLCS = plcs
	DTUS[dtu.Id] = dtu
	fmt.Println("Dtus:",DTUS,plcs)
}
