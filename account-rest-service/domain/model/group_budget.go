package model

import "time"

type GroupStandardBudgets struct {
	GroupStandardBudgets []GroupStandardBudgetByCategory `json:"standard_budgets"`
}

type GroupStandardBudgetByCategory struct {
	BigCategoryID     int    `json:"big_category_id"      db:"big_category_id"`
	BigCategoryName   string `json:"big_category_name"    db:"big_category_name"`
	Budget            int    `json:"budget"               db:"budget"`
	LastMonthExpenses int    `json:"last_month_expenses"`
}

type GroupCustomBudgets struct {
	GroupCustomBudgets []GroupCustomBudgetByCategory `json:"custom_budgets"`
}

type GroupCustomBudgetByCategory struct {
	BigCategoryID     int    `json:"big_category_id"      db:"big_category_id"`
	BigCategoryName   string `json:"big_category_name"    db:"big_category_name"`
	Budget            int    `json:"budget"               db:"budget"`
	LastMonthExpenses int    `json:"last_month_expenses"`
}

type YearlyGroupBudget struct {
	Year                time.Time            `json:"year"`
	YearlyTotalBudget   int                  `json:"yearly_total_budget"`
	GroupMonthlyBudgets []MonthlyGroupBudget `json:"monthly_budgets"`
}

type MonthlyGroupBudget struct {
	Month              Months `json:"month"                db:"years_months"`
	BudgetType         string `json:"budget_type"`
	MonthlyTotalBudget int    `json:"monthly_total_budget" db:"total_budget"`
}

func NewGroupStandardBudgets(groupStandardBudgetByCategoryList []GroupStandardBudgetByCategory) GroupStandardBudgets {
	return GroupStandardBudgets{GroupStandardBudgets: groupStandardBudgetByCategoryList}
}

func NewGroupCustomBudgets(groupCustomBudgetByCategoryList []GroupCustomBudgetByCategory) GroupCustomBudgets {
	return GroupCustomBudgets{GroupCustomBudgets: groupCustomBudgetByCategoryList}
}

func NewYearlyGroupBudget(year time.Time) YearlyGroupBudget {
	return YearlyGroupBudget{
		Year:                year,
		GroupMonthlyBudgets: make([]MonthlyGroupBudget, 12),
	}
}

func (b GroupStandardBudgets) ShowBudgetsList() []int {
	budgetsList := make([]int, len(b.GroupStandardBudgets))
	for i := 0; i < len(b.GroupStandardBudgets); i++ {
		budgetsList[i] = b.GroupStandardBudgets[i].Budget
	}

	return budgetsList
}

func (b GroupCustomBudgets) ShowBudgetsList() []int {
	budgetsList := make([]int, len(b.GroupCustomBudgets))
	for i := 0; i < len(b.GroupCustomBudgets); i++ {
		budgetsList[i] = b.GroupCustomBudgets[i].Budget
	}

	return budgetsList
}
