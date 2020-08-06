package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

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

type YearlyBudget struct {
	Year              time.Time       `json:"year"`
	YearlyTotalBudget int             `json:"yearly_total_budget"`
	MonthlyBudgets    []MonthlyBudget `json:"monthly_budgets"`
}

type MonthlyBudget struct {
	Month              Months `json:"month"                db:"years_months"`
	BudgetType         string `json:"budget_type"`
	MonthlyTotalBudget int    `json:"monthly_total_budget" db:"total_budget"`
}

type Months struct {
	time.Time
	String string
}

func NewStandardBudgets(standardBudgetByCategoryList []StandardBudgetByCategory) StandardBudgets {
	return StandardBudgets{StandardBudgets: standardBudgetByCategoryList}
}

func NewCustomBudgets(customBudgetByCategoryList []CustomBudgetByCategory) CustomBudgets {
	return CustomBudgets{CustomBudgets: customBudgetByCategoryList}
}

func NewYearlyBudget(year time.Time) YearlyBudget {
	return YearlyBudget{
		Year:           year,
		MonthlyBudgets: make([]MonthlyBudget, 12),
	}
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

func (m *Months) Scan(value interface{}) error {
	month, ok := value.(time.Time)
	if !ok {
		return errors.New("type assertion error")
	}
	m.Time = month
	return nil
}

func (m Months) Value() (driver.Value, error) {
	return m.Time, nil
}

func (m Months) MarshalJSON() ([]byte, error) {
	months := [...]string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"}
	m.String = months[m.Time.Month()-1]

	return json.Marshal(m.String)
}
