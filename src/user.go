package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type User struct {
	CommonModle
	UID     string   `json:"uid" gorm:"column:uid"`
	Name    string   `json:"name"`
	Baobeis []Baobei `gorm:"many2many:user_baobeis;"`
}

func getUserFavoriesHandler(context *gin.Context) {
	uid, _ := context.GetQuery("uid")
	list, err := getUserFavories(uid)
	renderJSON(list, err, context)
}
func getUserFavories(uid string) (list []Baobei, err error) {

	u, err := getUserByUID(uid)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return make([]Baobei, 0), nil
	}
	conn := GetDB()
	if err = conn.Model(&u).Related(&u.Baobeis, "Baobeis").Error; err != nil {
		gLogger.ELog("query user baobeis error %s", err.Error)
		return
	}
	for i := 0; i < len(u.Baobeis); i++ {
		b := &u.Baobeis[i]
		b.InFavor = true
	}
	return u.Baobeis, err
}

func getUserByUID(uid string) (user *User, err error) {
	conn := GetDB()
	// var u User
	user = &User{UID: uid}
	if err = conn.Where("uid = ?", uid).First(user).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			gLogger.ELog("query user error %s,uid %s", err.Error(), uid)
		} else {
			err = nil
		}
		return
	}
	return
}
