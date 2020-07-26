package model

type StandardBudgets struct {
	StandardBudgets []StandardBudgetByCategory `json:"standard_budgets"`
}

type StandardBudgetByCategory struct {
	BigCategoryID   int    `json:"big_category_id"   db:"big_category_id"`
	BigCategoryName string `json:"big_category_name" db:"big_category_name"`
	Budget          int    `json:"budget"            db:"budget"`
}

type CustomBudgets struct {
	CustomBudgets []CustomBudgetByCategory `json:"custom_budgets"`
}

type CustomBudgetByCategory struct {
	BigCategoryID   int    `json:"big_category_id"   db:"big_category_id"`
	BigCategoryName string `json:"big_category_name" db:"big_category_name"`
	Budget          int    `json:"budget"            db:"budget"`
}

func NewStandardBudgets(standardBudgetByCategoryList []StandardBudgetByCategory) StandardBudgets {
	return StandardBudgets{StandardBudgets: standardBudgetByCategoryList}
}

func NewCustomBudgets(customBudgetByCategoryList []CustomBudgetByCategory) CustomBudgets {
	return CustomBudgets{CustomBudgets: customBudgetByCategoryList}
}

func (b StandardBudgets) ShowBudgetsList() []int {
	budgetsList := make([]int, len(b.StandardBudgets))

	for i := 0; i < len(b.StandardBudgets); i++ {
		budgetsList[i] = b.StandardBudgets[i].Budget
	}

	return budgetsList
}

func (b CustomBudgets) ShowBudgetsList() []int {
	budgetsList := make([]int, len(b.CustomBudgets))

	for i := 0; i < len(b.CustomBudgets); i++ {
		budgetsList[i] = b.CustomBudgets[i].Budget
	}

	return budgetsList
}
