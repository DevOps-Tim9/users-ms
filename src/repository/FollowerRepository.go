package repository

import (
	"errors"
	"user-ms/src/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type IFollowerRepository interface {
	AddFollower(*model.Follower) (int, error)
	DeleteFollower(int) error
	GetFollowing(int) []model.Follower
	GetFollowers(int) []model.Follower
	RemoveFollowing(int, int) error
}

func NewFollowerRepository(database *gorm.DB) IFollowerRepository {
	return &FollowerRepository{
		database,
	}
}

type FollowerRepository struct {
	Database *gorm.DB
}

func (repo *FollowerRepository) AddFollower(follower *model.Follower) (int, error) {
	var check []model.Follower
	repo.Database.Where("follower_id = ? and following_id = ? ",
		follower.FollowerId, follower.FollowingId).Find(&check)
	if len(check) > 0 {
		return -1, errors.New("the couple already exists")
	}
	result := repo.Database.Create(&follower)
	if result.Error != nil {
		return -1, result.Error
	}

	return follower.ID, nil
}

func (repo *FollowerRepository) DeleteFollower(id int) error {
	result := repo.Database.Delete(&model.Follower{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *FollowerRepository) GetFollowing(id int) []model.Follower {
	var req []model.Follower
	repo.Database.Where("follower_id = ?", id).Find(&req)
	return req
}

func (repo *FollowerRepository) GetFollowers(id int) []model.Follower {
	var req []model.Follower
	repo.Database.Where("following_id = ?", id).Find(&req)
	return req
}

func (repo *FollowerRepository) RemoveFollowing(id int, followingId int) error {
	var follower model.Follower
	result := repo.Database.Where("follower_id = ? and following_id = ?", id, followingId).Delete(&follower)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
