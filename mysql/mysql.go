package mysql

import (
	"context"
	"database/sql"
	"fmt"
	. "tool/config"
	logger "tool/logger"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	DB          *sql.DB
	IsConnected bool
	Config      ToolConfig
}

// Connect 初始化mysql
func (c *Connection) Connect() (conn *sql.Conn, err error) {
	//连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", c.Config.DBConfig.Username, c.Config.DBConfig.Password, c.Config.DBConfig.Host, c.Config.DBConfig.Port, c.Config.DBConfig.Database, "charset=utf8&multiStatements=true")
	logger.Debug("mysql dsn：", dsn)
	//Open只会验证dsb的格式是否正确,不会验证是否连接成功,同理,密码是否正确也不知道
	c.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	conn, err = c.GetConnect()
	if err != nil {
		logger.Error(err.Error())
		//panic(err)
		return nil, err
	}
	return
}

func (c *Connection) GetConnect() (conn *sql.Conn, err error) {
	conn, err = c.DB.Conn(context.Background())
	logger.Info("Connect mysql server success [%s:%s]", c.Config.DBConfig.Host, c.Config.DBConfig.Port)
	return
}

func (c *Connection) CheckConnect() bool {
	// 此时尝试连接数据库,会判断用户,密码,ip地址,端口是否正确
	err := c.DB.Ping()
	if err != nil {
		logger.Error(err)
		return false
	}
	c.IsConnected = true
	return true
}

func (c *Connection) Close() (err error) {
	err = c.DB.Close()
	c.IsConnected = false
	if err != nil {
		return err
	}
	logger.Info("Close connect success", c.Config.DBConfig.Host, ":", c.Config.DBConfig.Port)
	return nil
}
