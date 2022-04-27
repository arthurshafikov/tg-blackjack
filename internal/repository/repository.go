package repository

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{}
}
