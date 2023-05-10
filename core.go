package main

import (
	"os"
	"strings"
	. "tool/config"
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

func Readfile() string {
	file, err := os.Open(JConfig.InputFile)
	if err != nil {
		logging.Error(err.Error())
		return ""
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		logging.Error(err.Error())
		return ""
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		logging.Error(err.Error())
		return ""
	}

	logging.Info("bytes read: ", bytesread)
	// fmt.Println("bytestream to string: ", string(buffer))
	return string(buffer)
}

func SplitLine(lines string) []*fpbd_assets {

	deviceList := make([]*fpbd_assets, 0)
	arrLines := strings.Split(lines, "\n")
	for i, line := range arrLines {

		if i == 0 {
			continue
		}
		device := new(fpbd_assets)
		block := strings.Split(line, ":")
		if len(block) == 8 {
			device.AssetsId = JConfig.LockInfo.StartID + i
			device.StationId = JConfig.LockInfo.StationId
			device.RelateDevice = block[0]
			device.AssetsType = block[1]
			if device.AssetsType == "机械锁" || device.AssetsType == "位置锁" || device.AssetsType == "其他" {
				device.AssetsModel = "FSWG-4K"
			} else if device.AssetsType == "电气锁" {
				device.AssetsModel = "FSZL-2A"
			}

			device.AssetsNum = block[4]
			device.AssetsName = block[6]

			device.ProductionDate = JConfig.LockInfo.ProductionDate
			device.RegistDate = JConfig.LockInfo.RegistDate
			device.Register = JConfig.LockInfo.Register
			device.CommissionDate = JConfig.LockInfo.CommissionDate
			device.State = JConfig.LockInfo.State
			device.ValidDate = JConfig.LockInfo.ValidDate
			device.Note = JConfig.LockInfo.Note
			deviceList = append(deviceList, device)
		} else {
			logging.Error("line err", line)
		}
	}
	logging.Info("split line end")
	return deviceList
}
