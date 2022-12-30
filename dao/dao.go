package dao

import (
	"go.backend/conf"
	"go.backend/initialization"
)

var dao *Dao

type Dao struct {
	c     *conf.AppConfig
	mysql *initialization.Mysql
}

func New(c *conf.AppConfig) *Dao {
	// 单例模式
	// todo 项目启动全局加载，下一版本考虑如何按需加载
	if dao == nil {
		dao = &Dao{
			c: c,
			// 初始化mysql
			//mysql: initialization.InitMysql(c.MySQLConfig),
		}
	}
	return dao
}

// Close release resource
func (d *Dao) Close() error {
	return nil
}
