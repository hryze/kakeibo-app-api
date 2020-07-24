package model

type StandardBudgets struct {
	StandardBudgets []StandardBudgetByCategory `json:"standard_budgets"`
}

type StandardBudgetByCategory struct {
	BigCategoryID   int    `json:"big_category_id"   db:"big_category_id"`
	BigCategoryName string `json:"big_category_name" db:"big_category_name"`
	Budget          int    `json:"budget"            db:"budget"`
}

func NewStandardBudgets(standardBudgetByCategoryList []StandardBudgetByCategory) StandardBudgets {
	return StandardBudgets{StandardBudgets: standardBudgetByCategoryList}
}
