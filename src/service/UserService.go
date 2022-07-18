package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"user-ms/src/auth0"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/repository"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo    repository.IUserRepository
	Auth0Client auth0.Auth0Client
	Logger      *logrus.Entry
}

type IUserService interface {
	Register(*dto.RegistrationRequestDTO) (int, error)
	GetByEmail(string) (*dto.UserResponseDTO, error)
	GetByUsername(string) []dto.UserResponseDTO
	Update(*dto.UserUpdateDTO) (*dto.UserResponseDTO, error)
	BlockUser(int, string) error
	UnblockUser(int, string) error
	GetBlockedUsers(string) []dto.BlockedUserDTO
	SetNotifications(*dto.NotificationsUpdateDTO, string) error
	GetNotifications(string) *dto.NotificationsUpdateDTO
}

func NewUserService(userRepository repository.IUserRepository, auth0Client auth0.Auth0Client, logger *logrus.Entry) IUserService {
	return &UserService{
		userRepository,
		auth0Client,
		logger,
	}
}

func (service *UserService) Register(userToRegister *dto.RegistrationRequestDTO) (int, error) {
	if strings.TrimSpace(userToRegister.Password) == "" || len(userToRegister.Password) < 8 {
		service.Logger.Debug("Password must be at least 8 characters long!")
		return -1, errors.New("Password must be at least 8 characters long!")
	}
	if match, _ := regexp.MatchString(".*\\d.*", userToRegister.Password); !match {
		service.Logger.Debug("Password must contain at least one number!")
		return -1, errors.New("Password must contain at least one number!")
	}

	user := mapper.RegistrationRequestDTOToUser(userToRegister)

	err := user.Validate()
	if err != nil {
		service.Logger.Debug(err.Error())
		return -1, err
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		service.Logger.Debug(err.Error())
		return -1, err
	}

	service.Logger.Info(fmt.Sprintf("Adding user to database with username %s", user.Username))
	addedUserID, err := service.UserRepo.AddUser(user)
	if err != nil {
		service.Logger.Debug(err.Error())
		return -1, err
	}

	service.Logger.Info(fmt.Sprintf("Calling Auth0 to add user with username %s", user.Username))
	if auth0ID, err := service.Auth0Client.Register(userToRegister.Email, userToRegister.Password); err != nil {
		service.Logger.Debug(err.Error())
		service.Logger.Info(fmt.Sprintf("Deleting user with username %s", user.Username))
		if err = service.UserRepo.DeleteUser(addedUserID); err != nil {
			service.Logger.Debug(err.Error())
			return -1, err
		}
		return -1, err
	} else {
		service.Logger.Info(fmt.Sprintf("Updating user with username %s", user.Username))
		user.Auth0ID = auth0ID
		service.UserRepo.Update(user)
	}

	service.Logger.Info(fmt.Sprintf("Succesfully registered user with id %d", addedUserID))
	return addedUserID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (service *UserService) GetByUsername(username string) []dto.UserResponseDTO {
	service.Logger.Info(fmt.Sprintf("Getting users by username %s", username))

	users := service.UserRepo.GetByUsername(username)

	res := make([]dto.UserResponseDTO, len(users))

	for i := 0; i < len(users); i++ {
		res[i] = *mapper.UserToDTO(&users[i])
	}

	return res
}

func (service *UserService) GetByEmail(email string) (*dto.UserResponseDTO, error) {
	service.Logger.Info("Getting user by email")
	return service.UserRepo.GetByEmail(email)
}

func (service *UserService) GetAll() ([]*dto.UserResponseDTO, error) {
	service.Logger.Info("Getting all users")
	return service.UserRepo.GetAll()
}

func (service *UserService) GetByID(id int) (*dto.UserResponseDTO, error) {
	service.Logger.Info(fmt.Sprintf("Getting user by id %d", id))
	user, err := service.UserRepo.GetByID(id)
	if err == nil {
		userResponse := *mapper.UserToDTO(user)
		return &userResponse, err
	}
	return nil, err
}

func (service *UserService) Update(userToUpdate *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	service.Logger.Info(fmt.Sprintf("Updating user with id %d", userToUpdate.ID))
	userEntity, errr := service.UserRepo.GetByID(userToUpdate.ID)
	if errr != nil {
		service.Logger.Debug(errr.Error())
		return nil, errr
	}

	user := mapper.UserUpdateDTOToUser(userToUpdate)
	user.Password = userEntity.Password
	user.Auth0ID = userEntity.Auth0ID

	err := user.Validate()
	if err != nil {
		service.Logger.Debug(err.Error())
		return nil, err
	}

	userDTO, err := service.UserRepo.Update(user)
	if err != nil {
		service.Logger.Debug(err.Error())
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Updating user with id %d in Auth0", userToUpdate.ID))
	if err := service.Auth0Client.Update(user.Email, user.Auth0ID); err != nil {
		service.Logger.Debug(err.Error())
		return nil, err
	}

	service.Logger.Info(fmt.Sprintf("Succesfully updated user with %d", userToUpdate.ID))
	return userDTO, nil
}

func (service *UserService) BlockUser(blockingID int, userAuth0ID string) error {
	service.Logger.Info(fmt.Sprintf("Blocking user with id %d", blockingID))
	blockedUserEntity, err := service.UserRepo.GetByID(blockingID)
	if err != nil {
		return err
	}

	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	userEntity.Blocked = append(userEntity.Blocked, *blockedUserEntity)

	service.UserRepo.Update(userEntity)

	service.Logger.Info(fmt.Sprintf("Succesfully blocked user with %d", blockingID))
	return nil
}

func (service *UserService) UnblockUser(blockingID int, userAuth0ID string) error {
	service.Logger.Info(fmt.Sprintf("Unblocking user with id %d", blockingID))
	_, err := service.UserRepo.GetByID(blockingID)
	if err != nil {
		service.Logger.Debug(err.Error())
		return err
	}

	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	service.UserRepo.UnblockUser(blockingID, userEntity.ID)

	service.Logger.Info(fmt.Sprintf("Succesfully unblocked user with %d", blockingID))
	return nil
}

func (service *UserService) GetBlockedUsers(userAuth0ID string) []dto.BlockedUserDTO {
	service.Logger.Info(fmt.Sprintf("Getting blocked users for user %s", userAuth0ID))
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)
	blockedUsers := service.UserRepo.GetBlockedUsers(userEntity.ID)

	res := make([]dto.BlockedUserDTO, len(blockedUsers))
	for i := 0; i < len(blockedUsers); i++ {
		res[i] = *mapper.UserToBlockedUserDTO(&blockedUsers[i])
	}
	return res
}

func (service *UserService) SetNotifications(notificationSettings *dto.NotificationsUpdateDTO, userAuth0ID string) error {
	service.Logger.Info(fmt.Sprintf("Setting notifications for user with %s", userAuth0ID))
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	userEntity.MessageNotifications = notificationSettings.MessageNotifications
	userEntity.FollowNotifications = notificationSettings.FollowNotifications
	userEntity.CommentNotifications = notificationSettings.CommentNotifications
	userEntity.LikeNotifications = notificationSettings.LikeNotifications

	_, err := service.UserRepo.Update(userEntity)
	if err != nil {
		service.Logger.Debug(err.Error())
		return err
	}

	service.Logger.Info(fmt.Sprintf("Succesfully set notification settings for user %s", userAuth0ID))
	return nil
}

func (service *UserService) GetNotifications(userAuth0ID string) *dto.NotificationsUpdateDTO {
	service.Logger.Info(fmt.Sprintf("Getting notification settings for user %s", userAuth0ID))
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	notificationSettings := mapper.UserToNotificationsDTO(userEntity)

	return notificationSettings
}

func (service *UserService) GetByParam(param string) ([]*dto.UserResponseDTO, error) {
	service.Logger.Info("Getting users by param")
	return service.UserRepo.GetBySearchParam(param)
}
