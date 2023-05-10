package main

import (
	"flag"
	. "tool/config"
	"tool/logger"
	// "tool/mysql"
)

var logging = logger.GetStaticLogger()

func main() {
	var ReConfig bool
	logging.SetRollingDaily("./logs", "bacth.log")
	flag.BoolVar(&ReConfig, "r", false, "是否恢复默认配置")
	//解析命令行参数，写入注册的flag里面
	flag.Parse()
	if ReConfig {
		if err := JConfig.ReDefault(&JConfig); err == nil {
			logging.Info("恢复默认参数成功")
			if err := JConfig.WriteFile(); err == nil {
				logging.Info("生成配置文件成功")
			} else {
				logging.Error("生成配置文件失败", err)
			}
		}
		return
	}
	logging.SetRollingDaily("./logs", "bacth.log")
	line := Readfile()
	SplitLine(line)
	// lineList := splitLine(line)

	// connect := mysql.Connection{}
	// defer func() {
	// 	err := connect.Close()
	// 	if err != nil {
	// 		return
	// 	}
	// }()

	// conn, err := connect.Connect()
	// if err != nil {
	// 	return
	// }

	// var opera mysql.OperatorInterface = &mysql.ConnDB{
	// 	Conn:    conn,
	// 	CTX:     context.TODO(),
	// 	IsDebug: true,
	// }

	// for _, line := range lineList {
	// 	ret, inter := mysql.GetInsertSql(line)
	// 	opera.ExecControl(ret, inter...)
	// }

}
