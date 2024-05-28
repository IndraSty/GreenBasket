package service

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type notificationService struct {
	repo     domain.NotificationRepository
	tmplRepo domain.TemplateRepository
	hub      *dto.Hub
}

func NewNotificationService(r domain.NotificationRepository, tr domain.TemplateRepository, hub *dto.Hub) domain.NotificationService {
	return &notificationService{
		repo:     r,
		tmplRepo: tr,
		hub:      hub,
	}
}

// FindByUser implements domain.NotificationService.
func (s *notificationService) FindByUser(ctx context.Context, userId string) ([]dto.NotificationRes, error) {
	notifications, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to find user notification :" + err.Error())
	}

	var result []dto.NotificationRes
	for _, v := range notifications {
		result = append(result, dto.NotificationRes{
			ID:         v.ID,
			Title:      v.Title,
			Body:       v.Body,
			Status:     v.Status,
			IsRead:     v.IsRead,
			Created_At: v.Created_At,
		})
	}

	if result == nil {
		result = make([]dto.NotificationRes, 0)
	}

	return result, nil
}

// Insert implements domain.NotificationService.
func (s *notificationService) Insert(ctx context.Context, email string, code string, data map[string]string) error {
	tmpl, err := s.tmplRepo.FindByCode(ctx, code)
	if err != nil {
		return errors.New("failed to find template notification :" + err.Error())
	}

	if tmpl == (&domain.Template{}) {
		return errors.New("template not found")
	}

	body := new(bytes.Buffer)
	tp := template.Must(template.New("notif").Parse(tmpl.Body))
	err = tp.Execute(body, data)
	if err != nil {
		return err
	}

	notification := domain.Notification{
		ID:         primitive.NewObjectID(),
		Email:      email,
		Title:      tmpl.Title,
		Body:       body.String(),
		Status:     1,
		IsRead:     false,
		Created_At: time.Now(),
	}

	err = s.repo.Insert(ctx, &notification)
	if err != nil {
		return errors.New("failed to insert notification :" + err.Error())
	}

	if channel, ok := s.hub.NotificationChannel[email]; ok {
		channel <- dto.NotificationRes{
			ID:         notification.ID,
			Title:      notification.Title,
			Body:       notification.Body,
			Status:     notification.Status,
			IsRead:     notification.IsRead,
			Created_At: notification.Created_At,
		}
	}

	return nil
}
