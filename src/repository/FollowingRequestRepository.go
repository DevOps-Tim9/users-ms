package repository

import (
	"user-ms/src/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type IFollowingRequestRepository interface {
	AddFollowingRequest(*model.FollowingRequest) (int, error)
	UpdateFollowingRequest(int, *model.FollowingRequest) (*model.FollowingRequest, error)
	DeleteFollowingRequest(int) error
	GetRequests() []model.FollowingRequest
	GetRequestsByFollowingID(int) []model.FollowingRequest
}

func NewFollowingRequestRepository(database *gorm.DB) IFollowingRequestRepository {
	return &FollowingRequestRepository{
		database,
	}
}

type FollowingRequestRepository struct {
	Database *gorm.DB
}

func (repo *FollowingRequestRepository) AddFollowingRequest(followingRequest *model.FollowingRequest) (int, error) {
	result := repo.Database.Create(followingRequest)

	if result.Error != nil {
		return -1, result.Error
	}

	return followingRequest.ID, nil
}

func (repo *FollowingRequestRepository) UpdateFollowingRequest(reqId int, followingRequest *model.FollowingRequest) (*model.FollowingRequest, error) {
	followingRequest.ID = reqId
	result := repo.Database.Save(followingRequest)

	if result.Error != nil {
		return nil, result.Error
	}

	return followingRequest, nil
}

func (repo *FollowingRequestRepository) DeleteFollowingRequest(id int) error {
	result := repo.Database.Delete(&model.FollowingRequest{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *FollowingRequestRepository) GetRequests() []model.FollowingRequest {
	var req []model.FollowingRequest
	repo.Database.Find(&req)
	return req
}

func (repo *FollowingRequestRepository) GetRequestsByFollowingID(id int) []model.FollowingRequest {
	var req []model.FollowingRequest
	repo.Database.Where("following_id = ? and request_status = 0", id).Find(&req)
	return req
}
