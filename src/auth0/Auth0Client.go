package auth0

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	USER_ROLE_ID = "rol_IweSr3VWzIst8EDw"
)

type Auth0Client interface {
	Register(email string, password string) (string, error)
	getAPIToken() (string, error)
	setRole(string, string) error
	Update(email string, auth0ID string) error
}

type auth0Client struct {
	domain       string
	clientId     string
	clientSecret string
	audience     string
}

func NewAuth0Client(domain string, clientId string, clientSecret string, audience string) Auth0Client {
	return &auth0Client{
		domain,
		clientId,
		clientSecret,
		audience,
	}
}

type ApiTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type RegistrationRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Connection string `json:"connection"`
}

type RegistrationResponse struct {
	UserId string `json:"user_id"`
}

type RoleRequest struct {
	Roles []string `json:"roles"`
}

type UpdateRequest struct {
	Email      string `json:"email"`
	Connection string `json:"connection"`
}

func (c *auth0Client) Register(email string, password string) (string, error) {
	apiToken, err := c.getAPIToken()

	if err != nil {
		return "", err
	}

	// with this endpoint, user role cannot be set
	endpoint := fmt.Sprintf("https://%s/api/v2/users", c.domain)

	b, _ := json.Marshal(&RegistrationRequest{Email: email, Password: password, Connection: "Dislinkt-User"})

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		fmt.Println(res.StatusCode)
		b, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(b))
		return "", errors.New("Failed to register user on Auth0")
	}

	userRegistrationResponse := &RegistrationResponse{}
	json.NewDecoder(res.Body).Decode(userRegistrationResponse)

	userId := userRegistrationResponse.UserId

	if err := c.setRole(userId, apiToken); err != nil {
		return "", err
	}

	return userId, nil
}

func (c *auth0Client) getAPIToken() (string, error) {
	endpoint := fmt.Sprintf("https://%s/oauth/token", c.domain)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.clientId)
	data.Set("client_secret", c.clientSecret)
	data.Set("audience", c.audience)

	req, _ := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	apiTokenResponse := &ApiTokenResponse{}
	json.NewDecoder(res.Body).Decode(apiTokenResponse)

	return apiTokenResponse.AccessToken, nil
}

func (c *auth0Client) setRole(userId string, apiToken string) error {
	url := fmt.Sprintf("https://%s/api/v2/users/%s/roles", c.domain, userId)

	b, _ := json.Marshal(&RoleRequest{[]string{USER_ROLE_ID}})

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 204 {
		fmt.Println(res.StatusCode)
		fmt.Println(res.Body)
		return errors.New("Failed to assign role to user on Auth0")
	}

	return nil
}

func (c *auth0Client) Update(email string, auth0ID string) error {
	apiToken, err := c.getAPIToken()

	if err != nil {
		return err
	}

	// with this endpoint, user role cannot be set
	fmt.Println(auth0ID)
	endpoint := fmt.Sprintf("https://%s/api/v2/users/%s", c.domain, auth0ID)

	b, _ := json.Marshal(&UpdateRequest{Email: email, Connection: "Dislinkt-User"})

	req, _ := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode)
		b, err := io.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(b))
		return errors.New("Failed to update user on Auth0")
	}

	return nil
}
