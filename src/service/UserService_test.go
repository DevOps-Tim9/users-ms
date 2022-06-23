package service

import (
	"errors"
	"fmt"
	"testing"
	"user-ms/src/auth0"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"
	"user-ms/src/repository"
	"user-ms/src/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type UserServiceUnitTestsSuite struct {
	suite.Suite
	userRepositoryMock *repository.UserRepositoryMock
	auth0ClientMock    *auth0.Auth0ClientMock
	service            IUserService
}

func TestUserServiceUnitTestsSuite(t *testing.T) {
	suite.Run(t, new(UserServiceUnitTestsSuite))
}

func (suite *UserServiceUnitTestsSuite) SetupSuite() {
	suite.userRepositoryMock = new(repository.UserRepositoryMock)
	suite.auth0ClientMock = new(auth0.Auth0ClientMock)
	suite.service = NewUserService(suite.userRepositoryMock, suite.auth0ClientMock, utils.Logger())
}

func (suite *UserServiceUnitTestsSuite) TestNewUserService() {
	assert.NotNil(suite.T(), suite.service, "Service is nil")
}

func (suite *UserServiceUnitTestsSuite) TestUserService_Register_PasswordDoesntContainNumber() {
	userEntity := dto.RegistrationRequestDTO{}
	userEntity.Password = "password"

	passwordErr := errors.New("Password must contain at least one number!")
	_, err := suite.service.Register(&userEntity)

	assert.Equal(suite.T(), passwordErr, err)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_Register_PasswordIsLessThan8CharsLong() {
	userEntity := dto.RegistrationRequestDTO{}
	userEntity.Password = "pass"

	passwordErr := errors.New("Password must be at least 8 characters long!")
	_, err := suite.service.Register(&userEntity)

	assert.Equal(suite.T(), passwordErr, err)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_Register_ValidDataProvided() {
	gender := model.Female
	userDTO := dto.RegistrationRequestDTO{
		FirstName:   "Name",
		LastName:    "Surname",
		Email:       "test@test.com",
		Password:    "password123",
		Username:    "username",
		PhoneNumber: "06524775859",
		DateOfBirth: 4356785889,
		Gender:      &gender,
	}

	user := mapper.RegistrationRequestDTOToUser(&userDTO)
	user.Auth0ID = "123"

	forReturn := mapper.UserToDTO(user)

	suite.userRepositoryMock.On("AddUser", mock.AnythingOfType("*model.User")).Return(1, nil).Once()
	suite.auth0ClientMock.On("Register", userDTO.Email, userDTO.Password).Return("123", nil).Once()
	suite.userRepositoryMock.On("Update", mock.AnythingOfType("*model.User")).Return(forReturn, nil).Once()
	userID, err := suite.service.Register(&userDTO)

	assert.Equal(suite.T(), 1, userID)
	assert.Equal(suite.T(), nil, err)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByEmail_UserDoesNotExist() {
	email := "mail@mail.com"
	err := errors.New(fmt.Sprintf("User with email %s not found", email))

	suite.userRepositoryMock.On("GetByEmail", email).Return(nil, err).Once()

	_, userErr := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), err, userErr)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetByEmail_UserExists() {
	email := "mail@mail.com"
	userEntity := dto.UserResponseDTO{}

	suite.userRepositoryMock.On("GetByEmail", email).Return(&userEntity, nil).Once()

	retUser, _ := suite.service.GetByEmail(email)

	assert.Equal(suite.T(), userEntity, *retUser)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_Update_UserToUpdateNotFound() {
	userEntity := dto.UserUpdateDTO{}
	userEntity.ID = 1

	err := errors.New(fmt.Sprintf("User with ID %d not found", userEntity.ID))

	suite.userRepositoryMock.On("GetByID", userEntity.ID).Return(nil, err).Once()

	_, e := suite.service.Update(&userEntity)

	assert.Equal(suite.T(), err, e)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_Update_UserUpdated() {
	gender := model.Female
	user := dto.UserUpdateDTO{
		ID:          1,
		FirstName:   "novo ime",
		LastName:    "novo prezime",
		Email:       "test@test.com",
		PhoneNumber: "0612456373",
		Gender:      &gender,
		Username:    "username",
	}

	userEntity := model.User{
		ID:          1,
		FirstName:   "ime",
		LastName:    "prezime",
		Email:       "test@test.com",
		PhoneNumber: "0612456373",
		Gender:      &gender,
		Username:    "username",
		Auth0ID:     "123",
		Password:    "123",
	}

	suite.userRepositoryMock.On("GetByID", user.ID).Return(&userEntity, nil).Once()

	toUpdate := mapper.UserUpdateDTOToUser(&user)
	toUpdate.Auth0ID = userEntity.Auth0ID
	toUpdate.Password = "123"

	forReturn := dto.UserResponseDTO{
		ID:          1,
		FirstName:   "novo ime",
		LastName:    "novo prezime",
		Email:       "test@test.com",
		PhoneNumber: "0612456373",
		Gender:      &gender,
		Username:    "username",
		Auth0ID:     "123",
	}
	suite.userRepositoryMock.On("Update", toUpdate).Return(&forReturn, nil).Once()

	suite.auth0ClientMock.On("Update", forReturn.Email, forReturn.Auth0ID).Return(nil).Once()

	updatedUser, err := suite.service.Update(&user)

	assert.Equal(suite.T(), user.FirstName, updatedUser.FirstName)
	assert.Equal(suite.T(), user.LastName, updatedUser.LastName)
	assert.Equal(suite.T(), user.Email, updatedUser.Email)
	assert.Equal(suite.T(), user.PhoneNumber, updatedUser.PhoneNumber)
	assert.Equal(suite.T(), user.Gender, updatedUser.Gender)
	assert.Equal(suite.T(), user.Username, updatedUser.Username)
	assert.Equal(suite.T(), nil, err)
}

func (suite *UserServiceUnitTestsSuite) TestUserService_GetBlockedUsers_NoBlockedUsersReturnsEmpty() {
	suite.userRepositoryMock.On("GetByAuth0ID", "1").Return(&model.User{ID: 1}, nil).Once()
	suite.userRepositoryMock.On("GetBlockedUsers", 1).Return([]model.User{}).Once()

	blockedUsers := suite.service.GetBlockedUsers("1")

	assert.Equal(suite.T(), 0, len(blockedUsers))
}
