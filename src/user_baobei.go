package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserBaobei struct {
	CommonModle
	UserId   uint
	BaobeiId uint
}

func addUserFavoriesHandler(context *gin.Context) {
	// context.Params.
	uid := context.PostForm("uid")
	gLogger.ILog("uid is %s", uid)
	baobei := context.PostForm("baobei_id")

	iBaobei, err := strconv.Atoi(baobei)
	if err != nil {
		renderJSON(nil, err, context)
		return
	}
	err = addUserBaobei(uid, uint(iBaobei))
	renderJSON(nil, err, context)
}

func addUserBaobei(uid string, baobei uint) error {
	u, err := getUserByUID(uid)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	err = GetDB().Create(&UserBaobei{UserId: u.ID, BaobeiId: baobei}).Error
	if err != nil {
		gLogger.ELog("create user favori error %s", err.Error)
	}
	return err
}
func rmUserFavoriesHandler(context *gin.Context) {
	// context.Params.
	uid := context.Query("uid")
	// gLogger.ILog("uid is %s", uid)
	baobei := context.Query("baobei_id")

	iBaobei, err := strconv.Atoi(baobei)
	if err != nil {
		renderJSON(nil, err, context)
		return
	}
	err = rmUserBaobei(uid, uint(iBaobei))
	renderJSON(nil, err, context)
}

func rmUserBaobei(uid string, baobei uint) error {
	u, err := getUserByUID(uid)
	if err != nil {
		return err
	}
	if u == nil {
		return nil
	}
	err = GetDB().Unscoped().Where(&UserBaobei{UserId: u.ID, BaobeiId: baobei}).Delete(&UserBaobei{}).Error
	if err != nil {
		gLogger.ELog("create user favori error %s", err.Error)
	}
	return err
}
