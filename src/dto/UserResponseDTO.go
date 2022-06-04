package dto

import "user-ms/src/model"

type UserResponseDTO struct {
	ID             int
	Auth0ID        string
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
