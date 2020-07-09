package repository

import "github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

type DBRepository interface {
	CategoriesRepository
}

type CategoriesRepository interface {
	GetUserID(sessionID string) (string, error)
	GetBigCategoriesList() ([]model.BigCategory, error)
	GetMediumCategoriesList() ([]model.MediumCategory, error)
	GetCustomCategoriesList(userID string) ([]model.CustomCategory, error)
}
