package service

import (
	"errors"
	"fmt"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"
	"user-ms/src/rabbitmq"
	"user-ms/src/repository"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type FollowingService struct {
	FollowerRepository         repository.IFollowerRepository
	FollowingRequestRepository repository.IFollowingRequestRepository
	UserRepository             repository.IUserRepository
	RabbitMQChannel            *amqp.Channel
	Logger                     *logrus.Entry
}

type IFollowingService interface {
	CreateRequest(*dto.FollowingRequestDTO) (int, error)
	UpdateRequest(int, *dto.FollowingRequestDTO) (*dto.FollowingRequestDTO, error)
	CreateFollower(*dto.FollowingRequestDTO) (int, error)
}

func NewFollowingService(followerRepository repository.IFollowerRepository, followingRequestRepository repository.IFollowingRequestRepository, userRepository repository.IUserRepository, channel *amqp.Channel, logger *logrus.Entry) IFollowingService {
	return &FollowingService{
		followerRepository,
		followingRequestRepository,
		userRepository,
		channel,
		logger,
	}
}

func (service *FollowingService) CreateRequest(request *dto.FollowingRequestDTO) (int, error) {
	service.Logger.Info("Requesting for follow")
	followingRequestId, err := service.FollowingRequestRepository.AddFollowingRequest(mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		service.Logger.Debug(err.Error())
		return -1, errors.New("can't create the request")
	}

	follower, _ := service.UserRepository.GetByID(request.FollowerId)
	following, _ := service.UserRepository.GetByID(request.FollowingId)

	followType := dto.Follow
	notification := dto.NotificationDTO{Message: fmt.Sprintf("%s requested to follow you.", follower.Username), UserAuth0ID: following.Auth0ID, NotificationType: &followType}

	service.Logger.Info("Adding following request notification to rabbitMQ")
	rabbitmq.AddNotification(&notification, service.RabbitMQChannel)

	return followingRequestId, nil
}

func (service *FollowingService) UpdateRequest(reqId int, request *dto.FollowingRequestDTO) (*dto.FollowingRequestDTO, error) {
	service.Logger.Info("Updating following request with id %d", reqId)
	followingRequest, err := service.FollowingRequestRepository.UpdateFollowingRequest(reqId, mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		service.Logger.Debug(err.Error())
		return mapper.RequestToFollowingDTO(followingRequest), errors.New("can't create the request")
	}
	status := model.RequestStatus(request.RequestStatus)
	if model.ACCEPTED == status {
		service.Logger.Info("Accepting following request with id %d", reqId)
		_, _ = service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	}
	return mapper.RequestToFollowingDTO(followingRequest), nil
}

func (service *FollowingService) GetRequests() ([]model.FollowingRequest, error) {
	service.Logger.Info("Getting all requests")
	requests := service.FollowingRequestRepository.GetRequests()
	return requests, nil
}

func (service *FollowingService) GetRequestsByFollowingID(id int) ([]model.FollowingRequest, error) {
	service.Logger.Info("Getting following requests for id %d", id)
	requests := service.FollowingRequestRepository.GetRequestsByFollowingID(id)
	return requests, nil
}

func (service *FollowingService) CreateFollower(request *dto.FollowingRequestDTO) (int, error) {
	service.Logger.Info("Creating follower")
	followerId, err := service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	if err != nil {
		service.Logger.Debug(err.Error())
		return -1, errors.New("can't create the request")
	}

	follower, _ := service.UserRepository.GetByID(request.FollowerId)
	following, _ := service.UserRepository.GetByID(request.FollowingId)

	followType := dto.Follow
	notification := dto.NotificationDTO{Message: fmt.Sprintf("%s started following you.", follower.Username), UserAuth0ID: following.Auth0ID, NotificationType: &followType}

	service.Logger.Info("Adding following notification to rabbitMQ")
	rabbitmq.AddNotification(&notification, service.RabbitMQChannel)

	return followerId, nil
}

func (service *FollowingService) GetFollowers(id int) ([]model.Follower, error) {
	service.Logger.Info("Getting followers for user with id %d", id)
	followers := service.FollowerRepository.GetFollowers(id)
	return followers, nil
}

func (service *FollowingService) GetFollowing(id int) ([]model.Follower, error) {
	service.Logger.Info("Getting following for user with id %d", id)
	following := service.FollowerRepository.GetFollowing(id)
	return following, nil
}

func (service *FollowingService) RemoveFollowing(id int, followingId int) error {
	service.Logger.Info("User with id %d unfollowed user with id %d", id, followingId)
	service.FollowerRepository.RemoveFollowing(id, followingId)
	return nil
}
