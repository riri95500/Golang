package main

import (
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/riri95500/go-chat/adapter"
	"github.com/riri95500/go-chat/service"
)

func main() {
	// Création du gestionnaire de salles
	roomManager := service.GetRoomManager()

	// Création du gestionnaire d'utilisateurs
	userRepository := service.NewUserRepository()

	// Création de l'adaptateur REST
	restAdapter := adapter.NewRestAdapter(roomManager, userRepository)

	// Création du routeur Gin
	router := gin.Default()

	// Configuration du template HTML
	htmlTemplate, err := template.ParseFiles("template.html")
	if err != nil {
		fmt.Println("Failed to parse HTML template:", err)
		return
	}

	// Utilisation de l'adaptateur HTML
	adapterHTML := adapter.NewHTMLAdapter(htmlTemplate, roomManager, userRepository)

	// Routes
	router.GET("/room/:roomid", adapterHTML.GetRoom)
	router.POST("/room/:roomid", adapterHTML.PostRoom)
	router.DELETE("/room/:roomid", adapterHTML.DeleteRoom)
	router.GET("/stream/:roomid", adapterHTML.Stream)

	// Routes pour la gestion des utilisateurs
	router.GET("/users", restAdapter.GetUsers)
	router.POST("/users", restAdapter.AddUser)
	router.DELETE("/users/:userid", restAdapter.RemoveUser)

	// Lancement du serveur sur le port 8080
	router.Run(":8080")
}
