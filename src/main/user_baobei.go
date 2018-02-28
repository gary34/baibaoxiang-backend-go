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
	user, err := getUserFromContext(context)
	if err != nil {
		context.Abort()
		return
	}
	baobei := context.PostForm("baobei_id")

	iBaobei, err := strconv.Atoi(baobei)
	if err != nil {
		renderJSON(nil, err, context)
		return
	}
	err = addUserBaobei(user.ID, uint(iBaobei))
	renderJSON(nil, err, context)
}

func addUserBaobei(userId uint, baobei uint) error {

	err := GetDB().Create(&UserBaobei{UserId: userId, BaobeiId: baobei}).Error
	if err != nil {
		gLogger.ELog("create user favori error %s", err.Error)
	}
	return err
}

func rmUserFavoriesHandler(context *gin.Context) {
	user, err := getUserFromContext(context)
	if err != nil {
		context.Abort()
		return
	}
	baobei := context.Query("baobei_id")

	iBaobei, err := strconv.Atoi(baobei)
	if err != nil {
		renderJSON(nil, err, context)
		return
	}
	err = rmUserBaobei(user.ID, uint(iBaobei))
	renderJSON(nil, err, context)
}

func rmUserBaobei(userId uint, baobei uint) error {
	err := GetDB().Unscoped().Where(&UserBaobei{UserId: userId, BaobeiId: baobei}).Delete(&UserBaobei{}).Error
	if err != nil {
		gLogger.ELog("create user favori error %s", err.Error)
	}
	return err
}
