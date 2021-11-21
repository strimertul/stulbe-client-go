package stulbe

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"

	kvclient "github.com/strimertul/kilovolt-client-go/v6"
	"github.com/strimertul/stulbe/api"
)

type Client struct {
	Endpoint string
	Logger   logrus.FieldLogger
	KV       *kvclient.Client

	client *http.Client
	token  string
}

func NewClient(options ClientOptions) (*Client, error) {
	if options.Logger == nil {
		options.Logger = logrus.New()
	}

	client := &Client{
		Endpoint: options.Endpoint,
		Logger:   options.Logger,
		token:    "",
		client:   &http.Client{},
	}

	err := client.Authenticate(options.Username, options.AuthKey)
	if err != nil {
		return nil, err
	}
	options.Logger.Debug("client authenticated")

	// Create kilovolt client
	client.KV, err = kvclient.NewClient(options.Endpoint+"/ws", kvclient.ClientOptions{
		Logger: client.Logger,
		Headers: http.Header{
			"Authorization": []string{"Bearer " + client.token},
		},
	})
	options.Logger.Debug("kv client connected")

	return client, err
}

func (s *Client) Close() {
	if s.KV != nil {
		s.KV.Close()
	}
}

func (s *Client) Authenticate(user string, authKey string) error {
	body := new(bytes.Buffer)
	err := jsoniter.ConfigFastest.NewEncoder(body).Encode(api.AuthRequest{User: user, AuthKey: authKey})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(fmt.Sprintf("%s/api/auth", s.Endpoint), "application/json", body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return getAPIError(resp)
	}

	var reply api.AuthResponse
	err = jsoniter.ConfigFastest.NewDecoder(resp.Body).Decode(&reply)
	if err != nil {
		return err
	}

	s.token = reply.Token
	return nil
}

func (s *Client) authenticated() bool {
	return s.token != ""
}

func (s *Client) NewAuthRequest(method string, path string, body io.Reader) (*http.Request, error) {
	if !s.authenticated() {
		return nil, ErrNotAuthenticated
	}

	req, err := http.NewRequest(method, s.Endpoint+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	return req, nil
}

func getAPIError(r *http.Response) error {
	var apiError api.ResponseError
	err := jsoniter.ConfigFastest.NewDecoder(r.Body).Decode(&apiError)
	if err != nil {
		return err
	}
	return errors.New(apiError.Error)
}
