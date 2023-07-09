package main

import (
	"log"

	"github.com/riri95500/go-chat/config"
	"github.com/riri95500/go-chat/model"
	"github.com/riri95500/go-chat/service"

	"github.com/riri95500/go-chat/config"
	"github.com/riri95500/go-chat/service"
)

var roomManager service.Manager

func main() {
	conf := config.InitConfig()
	db, err := config.InitDB(conf)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.User{})

	// roomManager = service.GetRoomManager()
	// adapter := adapter.NewGinHTMLAdapter(roomManager)
	// router := gin.Default()
	// router.SetHTMLTemplate(adapter.Template)

	// router.GET("/room/:roomid", adapter.GetRoom)
	// router.POST("/room/:roomid", adapter.PostRoom)
	// router.DELETE("/room/:roomid", adapter.DeleteRoom)
	// router.GET("/stream/:roomid", adapter.Stream)

	// router.Run(fmt.Sprintf(":%v", 8080))
}
