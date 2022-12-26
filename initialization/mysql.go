package initialization

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormopentracing "gorm.io/plugin/opentracing"
	"hertz_demo/conf"
	"log"
)

type Mysql struct {
	*gorm.DB
}

func InitMysql(c *conf.MySQLConfig) *Mysql {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		c.User, c.Password, c.Host,
		c.Port, c.DB)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Println(err)
		return nil
	}
	if err = db.Use(gormopentracing.New()); err != nil {
		log.Println(err)
		return nil
	}
	return &Mysql{
		db,
	}
}
