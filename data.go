package stulbe

import (
	"errors"

	"go.uber.org/zap"
)

type StatusResponse struct {
	Ok bool `json:"ok"`
}

type ResponseError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type AuthRequest struct {
	User    string `json:"user"`
	AuthKey string `json:"key"`
}

type AuthResponse struct {
	Ok    bool   `json:"ok"`
	User  string `json:"username"`
	Level string `json:"level"`
	Token string `json:"token"`
}

type ClientOptions struct {
	Endpoint string
	Username string
	AuthKey  string

	Logger *zap.Logger
}

var (
	ErrNotAuthenticated     = errors.New("not authenticated")
	ErrSubscriptionNotFound = errors.New("subscription not found")
)
