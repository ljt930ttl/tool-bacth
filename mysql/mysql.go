package mysql

import (
	"context"
	"database/sql"
	"fmt"
	. "tool/config"
	"tool/logger"

	_ "github.com/go-sql-driver/mysql"
)

type Connection struct {
	DB          *sql.DB
	IsConnected bool
}

var logging = logger.GetStaticLogger()

// Connect 初始化mysql
func (c *Connection) Connect() (conn *sql.Conn, err error) {
	//连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=%s",
		JConfig.DBConfig.Username,
		JConfig.DBConfig.Password,
		JConfig.DBConfig.Host,
		JConfig.DBConfig.Port,
		JConfig.DBConfig.Database, "utf8", "true")
	logging.Debug("mysql dsn：", dsn)
	//Open只会验证dsb的格式是否正确,不会验证是否连接成功,同理,密码是否正确也不知道
	c.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	conn, err = c.GetConnect()
	if err != nil {
		logging.Error(err.Error())
		//panic(err)
		return nil, err
	}
	return
}

func (c *Connection) GetConnect() (conn *sql.Conn, err error) {
	conn, err = c.DB.Conn(context.Background())
	logging.Info("Connect mysql server success [%s:%s]", JConfig.DBConfig.Host, JConfig.DBConfig.Port)
	return
}

func (c *Connection) CheckConnect() bool {
	// 此时尝试连接数据库,会判断用户,密码,ip地址,端口是否正确
	err := c.DB.Ping()
	if err != nil {
		logging.Error(err.Error())
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
	logging.Info("Close connect success [%s:%s]", JConfig.DBConfig.Host, JConfig.DBConfig.Port)
	return nil
}
