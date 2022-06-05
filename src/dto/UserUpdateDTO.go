package dto

import "user-ms/src/model"

type UserUpdateDTO struct {
	ID             int
	FirstName      string
	LastName       string
	Email          string
	PhoneNumber    string
	Gender         *model.Gender
	Username       string
	DateOfBirth    float32
	Biography      string
	Education      string
	WorkExperience string
	Skills         string
	Interests      string
	Public         bool
}
