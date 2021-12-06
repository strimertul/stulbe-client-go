package stulbe

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/nicklaw5/helix/v2"
)

func (s *Client) TwitchUserInfo() (helix.User, error) {
	req, err := s.NewAuthRequest("GET", "/api/twitch/user", nil)
	if err != nil {
		return helix.User{}, err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return helix.User{}, err
	}
	defer res.Body.Close()
	var user helix.User
	err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&user)
	return user, err
}

func (s *Client) TwitchGetAuthenticationURL() (string, error) {
	req, err := s.NewAuthRequest("GET", "/api/twitch/authorize", nil)
	if err != nil {
		return "", err
	}
	res, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	var response struct {
		AuthenticationURL string `json:"auth_url"`
	}
	err = jsoniter.ConfigFastest.NewDecoder(res.Body).Decode(&response)
	return response.AuthenticationURL, err
}
