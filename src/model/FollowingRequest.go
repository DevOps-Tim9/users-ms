package model

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

type RequestStatus int

const (
	PENDING RequestStatus = iota
	ACCEPTED
	REJECTED
)

type FollowingRequest struct {
	ID            int           `json:"id" `
	FollowerId    int           `json:"followers_id" gorm:"TYPE:integer REFERENCES users" validate:"required"`
	FollowingId   int           `json:"following_id" gorm:"TYPE:integer REFERENCES users" validate:"required"`
	RequestStatus RequestStatus `json:"request_status" validate:"required"`
}

func (f *FollowingRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(f)
}

func (f *FollowingRequest) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(f)
}

func (f *FollowingRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}
