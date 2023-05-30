package main

import (
	"fmt"

	"go-chat/service"

	adapter "github.com/MohammadBnei/go-html-adapter/adapterHTML"

	"github.com/gin-gonic/gin"
)

var roomManager service.Manager

func main() {
	roomManager = service.GetRoomManager()
	adapter := adapter.NewGinHTMLAdapter(roomManager)
	router := gin.Default()
	router.SetHTMLTemplate(adapter.Template)

	router.GET("/room/:roomid", adapter.GetRoom)
	router.POST("/room/:roomid", adapter.PostRoom)
	router.DELETE("/room/:roomid", adapter.DeleteRoom)
	router.GET("/stream/:roomid", adapter.Stream)

	router.Run(fmt.Sprintf(":%v", 8080))
}
