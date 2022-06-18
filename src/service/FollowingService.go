package service

import (
	"errors"
	"fmt"
	"user-ms/src/dto"
	"user-ms/src/mapper"
	"user-ms/src/model"
	"user-ms/src/rabbitmq"
	"user-ms/src/repository"

	"github.com/streadway/amqp"
)

type FollowingService struct {
	FollowerRepository         repository.IFollowerRepository
	FollowingRequestRepository repository.IFollowingRequestRepository
	UserRepository             repository.IUserRepository
	RabbitMQChannel            *amqp.Channel
}

type IFollowingService interface {
	CreateRequest(*dto.FollowingRequestDTO) (int, error)
	UpdateRequest(int, *dto.FollowingRequestDTO) (*dto.FollowingRequestDTO, error)
	CreateFollower(*dto.FollowingRequestDTO) (int, error)
}

func NewFollowingService(followerRepository repository.IFollowerRepository, followingRequestRepository repository.IFollowingRequestRepository, userRepository repository.IUserRepository, channel *amqp.Channel) IFollowingService {
	return &FollowingService{
		followerRepository,
		followingRequestRepository,
		userRepository,
		channel,
	}
}

func (service *FollowingService) CreateRequest(request *dto.FollowingRequestDTO) (int, error) {
	followingRequestId, err := service.FollowingRequestRepository.AddFollowingRequest(mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		return -1, errors.New("can't create the request")
	}

	follower, _ := service.UserRepository.GetByID(request.FollowerId)
	following, _ := service.UserRepository.GetByID(request.FollowingId)

	followType := dto.Follow
	notification := dto.NotificationDTO{Message: fmt.Sprintf("%s requested to follow you.", follower.Username), UserAuth0ID: following.Auth0ID, NotificationType: &followType}

	rabbitmq.AddNotification(&notification, service.RabbitMQChannel)

	return followingRequestId, nil
}

func (service *FollowingService) UpdateRequest(reqId int, request *dto.FollowingRequestDTO) (*dto.FollowingRequestDTO, error) {
	followingRequest, err := service.FollowingRequestRepository.UpdateFollowingRequest(reqId, mapper.FollowingDTOToRequestFollower(request))
	if err != nil {
		return mapper.RequestToFollowingDTO(followingRequest), errors.New("can't create the request")
	}
	status := model.RequestStatus(request.RequestStatus)
	if model.ACCEPTED == status {
		_, _ = service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	}
	return mapper.RequestToFollowingDTO(followingRequest), nil
}

func (service *FollowingService) GetRequests() ([]model.FollowingRequest, error) {
	requests := service.FollowingRequestRepository.GetRequests()
	return requests, nil
}

func (service *FollowingService) GetRequestsByFollowingID(id int) ([]model.FollowingRequest, error) {
	requests := service.FollowingRequestRepository.GetRequestsByFollowingID(id)
	return requests, nil
}

func (service *FollowingService) CreateFollower(request *dto.FollowingRequestDTO) (int, error) {
	followerId, err := service.FollowerRepository.AddFollower(mapper.FollowingDTOToFollower(request))
	if err != nil {
		return -1, errors.New("can't create the request")
	}

	follower, _ := service.UserRepository.GetByID(request.FollowerId)
	following, _ := service.UserRepository.GetByID(request.FollowingId)

	followType := dto.Follow
	notification := dto.NotificationDTO{Message: fmt.Sprintf("%s started following you.", follower.Username), UserAuth0ID: following.Auth0ID, NotificationType: &followType}

	rabbitmq.AddNotification(&notification, service.RabbitMQChannel)

	return followerId, nil
}

func (service *FollowingService) GetFollowers(id int) ([]model.Follower, error) {
	followers := service.FollowerRepository.GetFollowers(id)
	return followers, nil
}

func (service *FollowingService) GetFollowing(id int) ([]model.Follower, error) {
	following := service.FollowerRepository.GetFollowing(id)
	return following, nil
}

func (service *FollowingService) RemoveFollowing(id int, followingId int) error {
	service.FollowerRepository.RemoveFollowing(id, followingId)
	return nil
}
