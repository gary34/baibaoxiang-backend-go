package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//StartServer 启动服务器
func StartServer(port string) {
	r := gin.Default()
	api := r.Group("/api")
	api.GET("/items.json", projectBaobeiListHandler)
	api.GET("/favoies.json", getUserFavoriesHandler)
	r.Run(fmt.Sprintf(":%s", port))
}

func renderJSON(obj interface{}, err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, obj)
	}
}
