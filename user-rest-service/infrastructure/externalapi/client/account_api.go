package client

import "net/http"

type AccountApiHandler struct {
	Client *http.Client
}

func NewAccountApiHandler() *AccountApiHandler {
	return &AccountApiHandler{
		Client: http.DefaultClient,
	}
}
