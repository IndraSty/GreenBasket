package sse

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/gin-gonic/gin"
)

type NotificationSSE struct {
	hub      *dto.Hub
	userRepo domain.UserRepository
}

func NewNotificationSSE(hub *dto.Hub, ur domain.UserRepository) *NotificationSSE {
	return &NotificationSSE{
		hub:      hub,
		userRepo: ur,
	}
}

func (s NotificationSSE) StreamNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")

		email := ctx.MustGet("email").(string)
		user, err := s.userRepo.FindUserByEmail(ctx, email)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find user: " + err.Error()})
			return
		}
		s.hub.NotificationChannel[user.User_Id] = make(chan dto.NotificationRes)
		log.Println("berhasil membuat channel user ID :", user.User_Id)
		ctx.Stream(func(w io.Writer) bool {
			event := fmt.Sprintf("event: %s\n"+
				"data: \n\n", "initial")
			_, _ = fmt.Fprint(w, event)

			for notification := range s.hub.NotificationChannel[user.User_Id] {
				data, _ := json.Marshal(notification)

				event = fmt.Sprintf("event: %s\n"+
					"data: %s\n\n", "notification-updated", data)

				_, _ = fmt.Fprint(w, event)
				log.Println("ini data notif ", notification.Body)
			}

			return true // continue the stream
		})
	}
}
