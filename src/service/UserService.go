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

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo    repository.IUserRepository
	Auth0Client auth0.Auth0Client
}

type IUserService interface {
	Register(*dto.RegistrationRequestDTO) (int, error)
	GetByEmail(string) (*dto.UserResponseDTO, error)
	Update(*dto.UserUpdateDTO) (*dto.UserResponseDTO, error)
	BlockUser(int, string) error
	UnblockUser(int, string) error
	GetBlockedUsers(string) []dto.BlockedUserDTO
	SetNotifications(*dto.NotificationsUpdateDTO, string) error
	GetNotifications(string) *dto.NotificationsUpdateDTO
}

func NewUserService(userRepository repository.IUserRepository, auth0Client auth0.Auth0Client) IUserService {
	return &UserService{
		userRepository,
		auth0Client,
	}
}

func (service *UserService) Register(userToRegister *dto.RegistrationRequestDTO) (int, error) {
	if strings.TrimSpace(userToRegister.Password) == "" || len(userToRegister.Password) < 8 {
		return -1, errors.New("Password must be at least 8 characters long!")
	}
	if match, _ := regexp.MatchString(".*\\d.*", userToRegister.Password); !match {
		return -1, errors.New("Password must contain at least one number!")
	}

	user := mapper.RegistrationRequestDTOToUser(userToRegister)

	err := user.Validate()
	if err != nil {
		return -1, err
	}

	user.Password, err = HashPassword(user.Password)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	addedUserID, err := service.UserRepo.AddUser(user)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	if auth0ID, err := service.Auth0Client.Register(userToRegister.Email, userToRegister.Password); err != nil {
		fmt.Println(err)
		if err = service.UserRepo.DeleteUser(addedUserID); err != nil {
			return -1, err
		}
		return -1, err
	} else {
		user.Auth0ID = auth0ID
		service.UserRepo.Update(user)
	}

	return addedUserID, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (service *UserService) GetByEmail(email string) (*dto.UserResponseDTO, error) {
	return service.UserRepo.GetByEmail(email)
}

func (service *UserService) GetByID(id int) (*dto.UserResponseDTO, error) {
	user, err := service.UserRepo.GetByID(id)
	userResponse := *mapper.UserToDTO(user)
	return &userResponse, err
}

func (service *UserService) Update(userToUpdate *dto.UserUpdateDTO) (*dto.UserResponseDTO, error) {
	userEntity, errr := service.UserRepo.GetByID(userToUpdate.ID)
	if errr != nil {
		return nil, errr
	}

	user := mapper.UserUpdateDTOToUser(userToUpdate)
	user.Password = userEntity.Password
	user.Auth0ID = userEntity.Auth0ID

	err := user.Validate()
	if err != nil {
		return nil, err
	}

	userDTO, err := service.UserRepo.Update(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := service.Auth0Client.Update(user.Email, user.Auth0ID); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return userDTO, nil
}

func (service *UserService) BlockUser(blockingID int, userAuth0ID string) error {
	blockedUserEntity, err := service.UserRepo.GetByID(blockingID)
	if err != nil {
		return err
	}

	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	userEntity.Blocked = append(userEntity.Blocked, *blockedUserEntity)

	service.UserRepo.Update(userEntity)

	return nil
}

func (service *UserService) UnblockUser(blockingID int, userAuth0ID string) error {
	_, err := service.UserRepo.GetByID(blockingID)
	if err != nil {
		return err
	}

	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	service.UserRepo.UnblockUser(blockingID, userEntity.ID)

	return nil
}

func (service *UserService) GetBlockedUsers(userAuth0ID string) []dto.BlockedUserDTO {
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)
	blockedUsers := service.UserRepo.GetBlockedUsers(userEntity.ID)

	res := make([]dto.BlockedUserDTO, len(blockedUsers))
	for i := 0; i < len(blockedUsers); i++ {
		res[i] = *mapper.UserToBlockedUserDTO(&blockedUsers[i])
	}
	return res
}

func (service *UserService) SetNotifications(notificationSettings *dto.NotificationsUpdateDTO, userAuth0ID string) error {
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	userEntity.MessageNotifications = notificationSettings.MessageNotifications
	userEntity.FollowNotifications = notificationSettings.FollowNotifications
	userEntity.CommentNotifications = notificationSettings.CommentNotifications
	userEntity.LikeNotifications = notificationSettings.LikeNotifications

	_, err := service.UserRepo.Update(userEntity)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (service *UserService) GetNotifications(userAuth0ID string) *dto.NotificationsUpdateDTO {
	userEntity, _ := service.UserRepo.GetByAuth0ID(userAuth0ID)

	notificationSettings := mapper.UserToNotificationsDTO(userEntity)

	return notificationSettings
}
