package main

import (
	"github.com/gin-gonic/gin"
)

//Project 项目分组
type Project struct {
	CommonModle
	Name    string   `json:"name"`
	Baobeis []Baobei `json:"items"`
}

func getAllProjectBaobeis() (list []Project, err error) {
	conn := GetDB()
	if conn == nil {
		return
	}
	if err = conn.Find(&list).Error; err != nil {
		gLogger.ELog("find project list error %s", err.Error())
		return
	}
	for i := 0; i < len(list); i++ {
		p := &list[i]
		conn.Model(p).Related(&p.Baobeis)
	}
	return list, nil
}

func projectBaobeiListHandler(context *gin.Context) {
	list, err := getAllProjectBaobeis()
	renderJSON(list, err, context)
}
