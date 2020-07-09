package model

type CategoriesList struct {
	CategoriesList []BigCategory `json:"categories_list"`
}

type BigCategory struct {
	Type                     string        `json:"type"`
	ID                       int           `json:"id"   db:"id"`
	Name                     string        `json:"name" db:"category_name"`
	AssociatedCategoriesList []interface{} `json:"associated_categories_list"`
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
