package model

type GroupCategoriesList struct {
	GroupIncomeBigCategoriesList  []GroupIncomeBigCategory  `json:"income_categories_list"`
	GroupExpenseBigCategoriesList []GroupExpenseBigCategory `json:"expense_categories_list"`
}

type GroupIncomeBigCategory struct {
	CategoryType                  string                    `json:"category_type"`
	TransactionType               string                    `json:"transaction_type"`
	ID                            int                       `json:"id"`
	Name                          string                    `json:"name"`
	GroupAssociatedCategoriesList []GroupAssociatedCategory `json:"associated_categories_list"`
}

type GroupExpenseBigCategory struct {
	CategoryType                  string                    `json:"category_type"`
	TransactionType               string                    `json:"transaction_type"`
	ID                            int                       `json:"id"`
	Name                          string                    `json:"name"`
	GroupAssociatedCategoriesList []GroupAssociatedCategory `json:"associated_categories_list"`
}

type GroupBigCategory struct {
	ID                              int    `db:"id"`
	Name                            string `db:"category_name"`
	TransactionType                 string `db:"transaction_type"`
	IncomeAssociatedCategoriesList  []GroupAssociatedCategory
	ExpenseAssociatedCategoriesList []GroupAssociatedCategory
}

type GroupAssociatedCategory struct {
	CategoryType  string `json:"category_type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

type GroupCustomCategory struct {
	CategoryType  string `json:"category_type"`
	ID            int    `json:"id"              db:"id"`
	Name          string `json:"name"            db:"category_name"`
	BigCategoryID int    `json:"big_category_id" db:"big_category_id"`
}

func NewGroupIncomeBigCategory(groupBigCategory *GroupBigCategory) GroupIncomeBigCategory {
	return GroupIncomeBigCategory{
		CategoryType:                  "IncomeBigCategory",
		TransactionType:               groupBigCategory.TransactionType,
		ID:                            groupBigCategory.ID,
		Name:                          groupBigCategory.Name,
		GroupAssociatedCategoriesList: groupBigCategory.IncomeAssociatedCategoriesList,
	}
}

func NewGroupExpenseBigCategory(groupBigCategory *GroupBigCategory) GroupExpenseBigCategory {
	return GroupExpenseBigCategory{
		CategoryType:                  "ExpenseBigCategory",
		TransactionType:               groupBigCategory.TransactionType,
		ID:                            groupBigCategory.ID,
		Name:                          groupBigCategory.Name,
		GroupAssociatedCategoriesList: groupBigCategory.ExpenseAssociatedCategoriesList,
	}
}

func NewGroupCustomCategory() GroupCustomCategory {
	return GroupCustomCategory{
		CategoryType: "CustomCategory",
	}
}
