package model

type CategoriesList struct {
	IncomeBigCategoriesList  []IncomeBigCategory  `json:"income_categories_list"`
	ExpenseBigCategoriesList []ExpenseBigCategory `json:"expense_categories_list"`
}

type IncomeBigCategory struct {
	CategoryType             string               `json:"category_type"`
	TransactionType          string               `json:"transaction_type"`
	ID                       int                  `json:"id"`
	Name                     string               `json:"name"`
	AssociatedCategoriesList []AssociatedCategory `json:"associated_categories_list"`
}

type ExpenseBigCategory struct {
	CategoryType             string               `json:"category_type"`
	TransactionType          string               `json:"transaction_type"`
	ID                       int                  `json:"id"`
	Name                     string               `json:"name"`
	AssociatedCategoriesList []AssociatedCategory `json:"associated_categories_list"`
}

type BigCategory struct {
	ID                              int    `db:"id"`
	Name                            string `db:"category_name"`
	TransactionType                 string `db:"transaction_type"`
	IncomeAssociatedCategoriesList  []AssociatedCategory
	ExpenseAssociatedCategoriesList []AssociatedCategory
}

type AssociatedCategory struct {
	CategoryType  string `json:"category_type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

type CustomCategory struct {
	CategoryType  string `json:"category_type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

type CategoriesID struct {
	MediumCategoryID NullInt64 `json:"medium_category_id"`
	CustomCategoryID NullInt64 `json:"custom_category_id"`
}

type CategoriesName struct {
	BigCategoryName    NullString `json:"big_category_name"    db:"big_category_name"`
	MediumCategoryName NullString `json:"medium_category_name" db:"medium_category_name"`
	CustomCategoryName NullString `json:"custom_category_name" db:"custom_category_name"`
}

func NewIncomeBigCategory(bigCategory *BigCategory) IncomeBigCategory {
	return IncomeBigCategory{
		CategoryType:             "IncomeBigCategory",
		TransactionType:          bigCategory.TransactionType,
		ID:                       bigCategory.ID,
		Name:                     bigCategory.Name,
		AssociatedCategoriesList: bigCategory.IncomeAssociatedCategoriesList,
	}
}

func NewExpenseBigCategory(bigCategory *BigCategory) ExpenseBigCategory {
	return ExpenseBigCategory{
		CategoryType:             "ExpenseBigCategory",
		TransactionType:          bigCategory.TransactionType,
		ID:                       bigCategory.ID,
		Name:                     bigCategory.Name,
		AssociatedCategoriesList: bigCategory.ExpenseAssociatedCategoriesList,
	}
}

func NewCustomCategory() CustomCategory {
	return CustomCategory{
		CategoryType: "CustomCategory",
	}
}
