package tests

import (
	"context"
	"log"
	"testing"

	"github.com/arthurshafikov/tg-blackjack/internal/config"
	"github.com/arthurshafikov/tg-blackjack/internal/logger"
	"github.com/arthurshafikov/tg-blackjack/internal/repository"
	"github.com/arthurshafikov/tg-blackjack/internal/repository/mongodb"
	"github.com/arthurshafikov/tg-blackjack/internal/services"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
)

var r *require.Assertions

type APITestSuite struct {
	suite.Suite

	db *mongo.Client

	repos    *repository.Repository
	services *services.Services
	config   *config.Config

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) SetupSuite() {
	r = s.Require()
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	s.config = &config.Config{
		Database: config.Database{
			Scheme:   "mongodb",
			Host:     "mongo",
			Username: "root",
			Password: "secret",
		},
	}

	mongo, err := mongodb.NewMongoDB(s.ctx, mongodb.Config{
		Scheme:   s.config.Database.Scheme,
		Host:     s.config.Database.Host,
		Username: s.config.Database.Username,
		Password: s.config.Database.Password,
	})
	if err != nil {
		log.Fatalln(err)
	}

	repository := repository.NewRepository(mongo)

	logger := logger.NewLogger()
	s.services = services.NewServices(services.Deps{
		Config:     s.config,
		Repository: repository,
		Logger:     logger,
	})
}

func (s *APITestSuite) TearDownTest() {
}

func (s *APITestSuite) TearDownSuite() {
	s.ctxCancel()
}
