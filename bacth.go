package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	. "tool/config"
	_ "tool/config"
	"tool/mysql"
)

type fpbd_assets struct {
	AssetsId       int
	StationId      int
	AssetsName     string
	AssetsModel    string
	AssetsNum      string
	AssetsType     string
	ProductionDate string
	RegistDate     string
	Register       string
	CommissionDate string
	State          string
	ValidDate      string
	RelateDevice   string
	Note           string
}

var conf = Config

func readfile() string {
	file, err := os.Open(conf.InputFile)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	fmt.Println("bytes read: ", bytesread)
	// fmt.Println("bytestream to string: ", string(buffer))
	return string(buffer)
}

func splitLine(lines string) []*fpbd_assets {

	deviceList := make([]*fpbd_assets, 0)
	arrLines := strings.Split(lines, "\n")
	for i, line := range arrLines {

		if i == 0 {
			continue
		}
		device := new(fpbd_assets)
		block := strings.Split(line, ":")
		if len(block) == 8 {
			device.AssetsId = conf.LockInfo.StartID + i
			device.StationId = conf.LockInfo.StationId
			device.RelateDevice = block[0]
			device.AssetsType = block[1]
			if device.AssetsType == "机械锁" || device.AssetsType == "位置锁" || device.AssetsType == "其他" {
				device.AssetsModel = "FSWG-4K"
			} else if device.AssetsType == "电气锁" {
				device.AssetsModel = "FSZL-2A"
			}

			device.AssetsNum = block[4]
			device.AssetsName = block[6]

			device.ProductionDate = conf.LockInfo.ProductionDate
			device.RegistDate = conf.LockInfo.RegistDate
			device.Register = conf.LockInfo.Register
			device.CommissionDate = conf.LockInfo.CommissionDate
			device.State = conf.LockInfo.State
			device.ValidDate = conf.LockInfo.ValidDate
			device.Note = conf.LockInfo.Note
			deviceList = append(deviceList, device)
		} else {
			fmt.Sprintln("line err", line)
		}
	}
	fmt.Print("end")
	return deviceList
}
func main() {
	var ReConfig bool
	flag.BoolVar(&ReConfig, "r", false, "是否恢复默认配置")
	//解析命令行参数，写入注册的flag里面
	flag.Parse()
	if ReConfig {
		if err := Config.ReDefault(&Config); err == nil {
			fmt.Println("恢复默认参数成功")
			if err := Config.WriteFile(); err == nil {
				fmt.Println("生成配置文件成功")
			} else {
				fmt.Println("生成配置文件失败", err)
			}
		}
		return
	}
	line := readfile()
	lineList := splitLine(line)

	connect := mysql.Connection{Config: conf}
	defer func(connect *mysql.Connection) {
		err := connect.Close()
		if err != nil {
			return
		}
	}(&connect)

	conn, err := connect.Connect()
	if err != nil {
		return
	}

	var opera mysql.Operator = &mysql.ConnDB{
		Conn:    conn,
		CTX:     context.TODO(),
		IsDebug: true,
	}

	for _, line := range lineList {
		ret, inter := mysql.GetInsertSql(line)
		opera.ExecControl(ret, inter...)
	}

}
