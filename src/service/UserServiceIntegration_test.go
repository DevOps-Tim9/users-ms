package service

import (
	"fmt"
	"os"
	"testing"
	"user-ms/src/auth0"
	"user-ms/src/dto"
	"user-ms/src/model"
	"user-ms/src/repository"
	"user-ms/src/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserServiceIntegrationTestSuite struct {
	suite.Suite
	service UserService
	db      *gorm.DB
	users   []model.User
}

func (suite *UserServiceIntegrationTestSuite) SetupSuite() {
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

	db.AutoMigrate(model.User{})

	userRepository := repository.UserRepository{Database: db}

	auth0Client := auth0.NewAuth0Client(os.Getenv("AUTH0_DOMAIN"), os.Getenv("AUTH0_CLIENT_ID"), os.Getenv("AUTH0_CLIENT_SECRET"), os.Getenv("AUTH0_AUDIENCE"))

	suite.db = db

	suite.service = UserService{
		UserRepo:    &userRepository,
		Auth0Client: auth0Client,
		Logger:      utils.Logger(),
	}

	gender := model.Female
	suite.users = []model.User{
		{
			ID:          1,
			FirstName:   "ime",
			LastName:    "prezime",
			Email:       "test@test.com",
			PhoneNumber: "0612456373",
			Gender:      &gender,
			Username:    "username",
			Auth0ID:     "auth0|62cac9f9117230969f05f366",
			Password:    "123",
		},
	}

	tx := suite.db.Begin()

	tx.Create(&suite.users[0])

	tx.Commit()
}

func TestUserServiceIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceIntegrationTestSuite))
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetByEmail_UserDoesNotExist() {
	email := "no@test.com"

	user, err := suite.service.GetByEmail(email)

	assert.Nil(suite.T(), user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetByEmail_UserExists() {
	email := "test@test.com"

	user, err := suite.service.GetByEmail(email)

	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), email, user.Email)
	assert.Nil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetByID_UserDoesNotExist() {
	id := 2000000

	user, err := suite.service.GetByID(id)

	assert.Nil(suite.T(), user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetByID_UserExists() {
	id := 1

	user, err := suite.service.GetByID(id)

	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), id, user.ID)
	assert.Nil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Register_PasswordLessThan8CharactersLong() {
	userDto := dto.RegistrationRequestDTO{}
	userDto.Password = "pass"

	user, err := suite.service.Register(&userDto)

	assert.Equal(suite.T(), -1, user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Register_RequiredFieldMissing() {
	userDto := dto.RegistrationRequestDTO{
		Username:  "username",
		FirstName: "test",
		LastName:  "test",
		Password:  "test",
	}

	user, err := suite.service.Register(&userDto)

	assert.Equal(suite.T(), -1, user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Register_ExistingEmail() {
	gender := model.Female
	userDto := dto.RegistrationRequestDTO{
		Username:    "usernametest",
		FirstName:   "test",
		LastName:    "test",
		Password:    "test",
		Email:       "test@test.com",
		PhoneNumber: "0123456",
		DateOfBirth: 1235679,
		Gender:      &gender,
	}

	user, err := suite.service.Register(&userDto)

	assert.Equal(suite.T(), -1, user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Register_Pass() {
	gender := model.Female
	userDto := dto.RegistrationRequestDTO{
		Username:    "test-username",
		FirstName:   "test",
		LastName:    "test",
		Password:    "testtest123",
		Email:       "testemailemail@test.com",
		PhoneNumber: "0123456",
		DateOfBirth: 1235679,
		Gender:      &gender,
	}

	user, err := suite.service.Register(&userDto)

	assert.NotNil(suite.T(), user)
	assert.Nil(suite.T(), err)

	suite.service.UserRepo.DeleteUser(user)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Update_UserDoesNotExist() {
	userDto := dto.UserUpdateDTO{}
	userDto.ID = 100

	user, err := suite.service.Update(&userDto)

	assert.Nil(suite.T(), user)
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_Update_UserExists() {
	gender := model.Female
	userDto := dto.UserUpdateDTO{
		ID:          1,
		FirstName:   "novo ime",
		LastName:    "novo prezime",
		Email:       "test@test.com",
		PhoneNumber: "0612456373",
		Gender:      &gender,
		Username:    "username",
	}

	user, err := suite.service.Update(&userDto)

	assert.NotNil(suite.T(), user)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), userDto.FirstName, user.FirstName)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_BlockUser_UserDoesNotExist() {

	err := suite.service.BlockUser(10000000, "auth0|62cac9f9117230969f05f366")

	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetBlockedUsers_EmptyList() {

	list := suite.service.GetBlockedUsers("auth0|62cac9f9117230969f05f366")

	assert.NotNil(suite.T(), list)
	assert.Equal(suite.T(), 0, len(list))
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_SetNotifications_Pass() {
	notifications := dto.NotificationsUpdateDTO{
		MessageNotifications: false,
		FollowNotifications:  true,
		LikeNotifications:    false,
		CommentNotifications: false,
	}

	err := suite.service.SetNotifications(&notifications, "auth0|62cac9f9117230969f05f366")

	n := suite.service.GetNotifications("auth0|62cac9f9117230969f05f366")

	assert.NotNil(suite.T(), notifications)
	assert.Equal(suite.T(), true, n.FollowNotifications)
	assert.Equal(suite.T(), false, n.MessageNotifications)
	assert.Equal(suite.T(), false, n.LikeNotifications)
	assert.Equal(suite.T(), false, n.CommentNotifications)

	assert.Nil(suite.T(), err)
}

func (suite *UserServiceIntegrationTestSuite) TestIntegrationUserService_GetNotifications_Pass() {

	notifications := suite.service.GetNotifications("auth0|62cac9f9117230969f05f366")

	assert.NotNil(suite.T(), notifications)
	assert.Equal(suite.T(), false, notifications.FollowNotifications)
	assert.Equal(suite.T(), false, notifications.MessageNotifications)
	assert.Equal(suite.T(), false, notifications.LikeNotifications)
	assert.Equal(suite.T(), false, notifications.CommentNotifications)
}
