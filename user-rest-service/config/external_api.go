package config

import "net/http"

type AccountApiHandler struct {
	HttpClient *http.Client
}

func NewAccountApiHandler() *AccountApiHandler {
	return &AccountApiHandler{
		HttpClient: http.DefaultClient,
	}
}
