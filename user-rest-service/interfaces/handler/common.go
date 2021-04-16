package handler

import "github.com/paypay3/kakeibo-app-api/user-rest-service/domain/repository"

type DBHandler struct {
	HealthRepo repository.HealthRepository
	AuthRepo   repository.AuthRepository
	UserRepo   repository.UserRepository
	GroupRepo  repository.GroupRepository
}
