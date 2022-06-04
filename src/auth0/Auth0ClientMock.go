package auth0

import (
	"github.com/stretchr/testify/mock"
)

type Auth0ClientMock struct {
	mock.Mock
}

func (a *Auth0ClientMock) Register(email string, password string) (string, error) {
	args := a.Called(email, password)
	if args.Get(1) == nil {
		return args.Get(0).(string), nil
	}
	return "", args.Get(1).(error)
}

func (a *Auth0ClientMock) getAPIToken() (string, error) {
	panic("implement me")
}

func (a *Auth0ClientMock) setRole(s string, s2 string) error {
	panic("implement me")
}
