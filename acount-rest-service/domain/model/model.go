package model

import (
	"encoding/json"
)

type CategoriesList struct {
	CategoriesList []BigCategory `json:"categories_list"`
}

type BigCategory struct {
	Type                     string               `json:"type"`
	ID                       int                  `json:"id"   db:"id"`
	Name                     string               `json:"name" db:"category_name"`
	AssociatedCategoriesList []AssociatedCategory `json:"associated_categories_list"`
}

type MediumCategory struct {
	Type          string `json:"type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

type CustomCategory struct {
	Type          string `json:"type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

type AssociatedCategory interface {
	showCategory() (string, error)
}

func (c MediumCategory) showCategory() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}

func (c CustomCategory) showCategory() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return string(b), err
	}
	return string(b), nil
}

func NewCategoriesList(bigCategoriesList []BigCategory) CategoriesList {
	return CategoriesList{
		CategoriesList: bigCategoriesList,
	}
}

func NewBigCategory() BigCategory {
	return BigCategory{
		Type: "BigCategory",
	}
}

func NewMediumCategory() MediumCategory {
	return MediumCategory{
		Type: "MediumCategory",
	}
}

func NewCustomCategory() CustomCategory {
	return CustomCategory{
		Type: "CustomCategory",
	}
}
