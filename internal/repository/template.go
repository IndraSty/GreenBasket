package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type templateRepository struct {
	Collection *mongo.Collection
}

func NewTemplateRepository(client *mongo.Client) domain.TemplateRepository {
	return &templateRepository{
		Collection: db.OpenCollection(client, "Templates"),
	}
}

// FindByCode implements domain.TemplateRepository.
func (repo *templateRepository) FindByCode(ctx context.Context, code string) (*domain.Template, error) {
	var template domain.Template
	filter := bson.M{"code": code}
	err := repo.Collection.FindOne(ctx, filter).Decode(&template)
	if err != nil {
		return nil, err
	}

	return &template, nil
}
