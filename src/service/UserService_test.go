package service

import (
	"errors"
	"testing"
	"user-ms/src/auth0"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"
	"user-ms/src/repository"

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
	suite.service = NewUserService(suite.userRepositoryMock, suite.auth0ClientMock)
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
