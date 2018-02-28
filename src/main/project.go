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

func getAllProjectBaobeis(uid string) (list []Project, err error) {
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
		if uid != "" {
			for j := 0; j < len(p.Baobeis); j++ {
				b := &p.Baobeis[j]
				var c int
				if err = conn.Model(&UserBaobei{}).
					Joins("inner join users on user_baobeis.user_id = users.id").
					Where("users.uid = ? and user_baobeis.baobei_id = ?", uid, b.ID).Count(&c).Error; err != nil {
					gLogger.ELog("find favor list error %s", err.Error())
					return
				}
				if c > 0 {
					b.InFavor = true
				}
			}
		}
	}
	return list, nil
}

func projectBaobeiListHandler(context *gin.Context) {
	uid, _ := context.GetQuery("uid")
	list, err := getAllProjectBaobeis(uid)
	renderJSON(list, err, context)
}
