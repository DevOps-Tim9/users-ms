package service

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	_ "user-ms/src/dto"
	"user-ms/src/model"
	"user-ms/src/rabbitmq"
	"user-ms/src/repository"
	"user-ms/src/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FollowingServiceIntegrationTestSuite struct {
	suite.Suite
	service   FollowingService
	db        *gorm.DB
	followers []model.Follower
	users     []model.User
}

func (suite *FollowingServiceIntegrationTestSuite) SetupSuite() {
	host := os.Getenv("DATABASE_DOMAIN")
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)
	db, _ := gorm.Open("postgres", connectionString)

	db.AutoMigrate(model.FollowingRequest{})
	db.AutoMigrate(model.Follower{})
	db.AutoMigrate(model.User{})

	followerRepository := repository.FollowerRepository{Database: db}
	followingRequestRepository := repository.FollowingRequestRepository{Database: db}
	userRepository := repository.UserRepository{Database: db}

	suite.db = db

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	suite.service = FollowingService{
		FollowerRepository:         &followerRepository,
		FollowingRequestRepository: &followingRequestRepository,
		UserRepository:             &userRepository,
		RabbitMQChannel:            channel,
		Logger:                     utils.Logger(),
	}

	gender := model.Female
	suite.users = []model.User{
		{
			ID:          1234,
			FirstName:   "ime",
			LastName:    "prezime",
			Email:       "test@test.com",
			PhoneNumber: "0612456373",
			Gender:      &gender,
			Username:    "username",
			Auth0ID:     "auth0|62cac9f9117230969f05f366",
			Password:    "123",
		},
		{
			ID:          2222,
			FirstName:   "ime2",
			LastName:    "prezim2",
			Email:       "tes2t@test.com",
			PhoneNumber: "0612456373",
			Gender:      &gender,
			Username:    "username2",
			Auth0ID:     "auth0|62cac9f9117230969f05f366",
			Password:    "123",
		},
	}

	suite.followers = []model.Follower{
		{
			ID:          1,
			FollowerId:  1234,
			FollowingId: 2222,
		},
		{
			ID:          1,
			FollowerId:  2222,
			FollowingId: 1234,
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.users[0])
	tx.Create(&suite.users[1])
	tx.Commit()

	tx.Create(&suite.followers[0])
	tx.Create(&suite.followers[1])

	tx.Commit()
}

func TestFollowingServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(FollowingServiceIntegrationTestSuite))
}

func (suite *FollowingServiceIntegrationTestSuite) TestIntegrationGetFollowers() {
	followers, err := suite.service.GetFollowers(2222)

	assert.Equal(suite.T(), 1, len(followers))
	assert.NotNil(suite.T(), err)
}
