package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	ID    primitive.ObjectID `json:"id" bson:"_id"`
	Code  string             `json:"code"`
	Title string             `json:"title"`
	Body  string             `json:"body"`
}

type TemplateRepository interface {
	FindByCode(ctx context.Context, code string) (*Template, error)
}
