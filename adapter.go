package adapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/riri95500/go-chat/service"
)

type RestAdapter struct {
	roomManager    service.Manager
	userRepository service.UserRepository
}

func NewRestAdapter(roomManager service.Manager, userRepository service.UserRepository) *RestAdapter {
	return &RestAdapter{
		roomManager:    roomManager,
		userRepository: userRepository,
	}
}

func (a *RestAdapter) Stream(c *gin.Context) {
	roomID := c.Param("roomid")

	// Ouvrir un flux SSE (Server-Sent Events)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	messageChan := a.roomManager.OpenListener(roomID)
	defer a.roomManager.CloseListener(roomID, messageChan)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	for {
		select {
		case message := <-messageChan:
			c.SSEvent("message", message)
			flusher.Flush()
		case <-c.Writer.CloseNotify():
			// Le client a fermé la connexion SSE
			return
		}
	}
}

func (a *RestAdapter) Submit(c *gin.Context) {
	var messageData struct {
		UserID string `json:"userId"`
		RoomID string `json:"roomId"`
		Text   string `json:"text"`
	}

	if err := c.ShouldBindJSON(&messageData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	a.roomManager.Submit(messageData.UserID, messageData.RoomID, messageData.Text)

	c.JSON(http.StatusOK, gin.H{"message": "Message submitted"})
}

func (a *RestAdapter) GetUsers(c *gin.Context) {
	users, err := a.userRepository.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (a *RestAdapter) AddUser(c *gin.Context) {
	var userData struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	userID := generateUserID()
	user := User{
		ID:       userID,
		Username: userData.Username,
	}

	if err := a.userRepository.AddUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"userID": userID})
}

func (a *RestAdapter) RemoveUser(c *gin.Context) {
	userID := c.Param("userid")

	if err := a.userRepository.RemoveUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed"})
}
