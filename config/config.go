package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
)

type ToolConfig struct {
	InputFile string `json:"InputFile" default:"device.fil"`
	DBConfig  struct {
		Host     string `json:"Host" default:"127.0.0.1"`
		Port     string `json:"Port" default:"3306"`
		Username string `json:"Username" default:"root"`
		Password string `json:"Password" default:"123456"`
		Database string `json:"Database" default:"fpbdsdb_child"`
	}
	LockInfo struct {
		StartID        int    `json:"StartID" default:"100"`
		StationId      int    `json:"StationId" default:"921"`
		ProductionDate string `json:"ProductionDate" default:"2015-01-20 00:00:00"`
		RegistDate     string `json:"RegistDate" default:"2015-01-20 00:00:00"`
		Register       string `json:"Register" default:"管理员"`
		CommissionDate string `json:"CommissionDate" default:"2015-08-13 00:00:00"`
		State          string `json:"Stateost" default:"正常"`
		ValidDate      string `json:"ValidDate" default:"2060-01-01 00:00:00"`
		Note           string `json:"Note" default:"-tool-批量导入"`
	}
}

type ConfigInterface interface {
	// 恢复默认值
	ReDefault(config *interface{}) error
	// 读取配置文件
	ReadFile() error
	// 写入配置文件
	WriteFile() error
}

// 配置
var Config ToolConfig

// 配置文件名
var configFile string = "config.json"

// 默认值标签名
var tagDefaultName string = "default"

// 恢复默认值
func (Config *ToolConfig) ReDefault(config interface{}) error {
	t, v := reflect.TypeOf(config), reflect.ValueOf(config)
	tE, vE := t.Elem(), v.Elem()
	if t.Kind() != reflect.Ptr || v.Kind() != reflect.Ptr {
		return errors.New("参数必须是指针")
	}
	if tE.Kind() != reflect.Struct || vE.Kind() != reflect.Struct {
		return errors.New("参数元素必须是结构体")
	}
	for i := 0; i < tE.NumField(); i++ {
		tEf, vEf := tE.Field(i), vE.Field(i)
		tagV := tEf.Tag.Get(tagDefaultName)
		switch vEf.Kind() {
		case reflect.Struct:
			if err := Config.ReDefault(vEf.Addr().Interface()); err != nil {
				return err
			}
		case reflect.String:
			if len(tagV) > 0 {
				vEf.SetString(tagV)
			}
		case reflect.Int:
			if len(tagV) > 0 {
				number, err := strconv.ParseInt(tagV, 0, 64)
				if err == nil {
					vEf.SetInt(number)
				} else {
					return err
				}
			}
		default:
			return errors.New("不能处理类型 " + vEf.Kind().String())
		}
	}
	return nil
}

// 读取配置文件
func (Config *ToolConfig) ReadFile() error {
	configFile, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err2 := json.Unmarshal(configFile, Config); err2 != nil {
		return err
	}
	return nil
}

// 写入配置文件
func (Config *ToolConfig) WriteFile() error {
	jsonStr, err := json.MarshalIndent(Config, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, jsonStr, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func init() {
	if err := Config.ReadFile(); err == nil {
		fmt.Println("读取配置文件成功")
	} else {
		fmt.Println("读取配置文件失败", err)
		if err := Config.ReDefault(&Config); err == nil {
			fmt.Println("恢复默认参数成功")
			if err := Config.WriteFile(); err == nil {
				fmt.Println("生成配置文件成功")
			} else {
				fmt.Println("生成配置文件失败", err)
			}
		} else {
			fmt.Println("恢复默认参数失败", err)
		}
	}
}
