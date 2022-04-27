package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Scheme   string
	Host     string
	Username string
	Password string
}

func NewMongoDB(ctx context.Context, config Config) (*mongo.Client, error) {
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(
			fmt.Sprintf("%s://%s:%s@%s", config.Scheme, config.Username, config.Password, config.Host),
		),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		<-ctx.Done()
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}
