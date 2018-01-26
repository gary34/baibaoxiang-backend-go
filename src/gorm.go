package main

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var gDBConnect *gorm.DB

//GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return gDBConnect
}

type CommonModle struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

//InitDB 初始化数据库
func InitDB(host, name, user, password, port string) error {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, port, name)
	db, err := gorm.Open("mysql", url)
	if err != nil {
		return err
	}
	db.LogMode(true)
	gDBConnect = db

	gDBConnect.AutoMigrate(&Baobei{}, &Project{})
	return nil
}
