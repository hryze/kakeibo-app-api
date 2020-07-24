package infrastructure

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type BudgetsRepository struct {
	*MySQLHandler
}

func (r *BudgetsRepository) PostInitStandardBudgets(userID string) error {
	query := `
        INSERT INTO standard_budgets
            (user_id, big_category_id)
        VALUES
            (?,2),
            (?,3),
            (?,4),
            (?,5),
            (?,6),
            (?,7),
            (?,8),
            (?,9),
            (?,10),
            (?,11),
            (?,12),
            (?,13),
            (?,14),
            (?,15),
            (?,16),
            (?,17)`
	_, err := r.MySQLHandler.conn.Exec(query, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID)
	return err
}

func (r *BudgetsRepository) GetStandardBudgets(userID string) (*model.StandardBudgets, error) {
	query := `
        SELECT
            standard_budgets.big_category_id big_category_id,
            big_categories.category_name big_category_name,
            standard_budgets.budget budget
        FROM
            standard_budgets
        LEFT JOIN
            big_categories
        ON
            standard_budgets.big_category_id = big_categories.id
        WHERE
            standard_budgets.user_id = ?
        ORDER BY
            standard_budgets.big_category_id`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var standardBudgetByCategoryList []model.StandardBudgetByCategory
	for rows.Next() {
		var standardBudgetByCategory model.StandardBudgetByCategory
		if err := rows.StructScan(&standardBudgetByCategory); err != nil {
			return nil, err
		}
		standardBudgetByCategoryList = append(standardBudgetByCategoryList, standardBudgetByCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	standardBudgets := model.NewStandardBudgets(standardBudgetByCategoryList)

	return &standardBudgets, nil
}

func (r *BudgetsRepository) PutStandardBudgets(standardBudgets *model.StandardBudgets, userID string) error {
	for _, standardBudgetByCategory := range standardBudgets.StandardBudgets {
		query := `
	   UPDATE
	       standard_budgets
	   SET
	       budget = ?
	   WHERE
	       user_id = ?
	   AND
	       big_category_id = ?`

		_, err := r.MySQLHandler.conn.Exec(query, standardBudgetByCategory.Budget, userID, standardBudgetByCategory.BigCategoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *BudgetsRepository) GetCustomBudgets(yearMonth time.Time, userID string) (*model.CustomBudgets, error) {
	query := `
        SELECT
            custom_budgets.big_category_id big_category_id,
            big_categories.category_name big_category_name,
            custom_budgets.budget budget
        FROM
            custom_budgets
        LEFT JOIN
            big_categories
        ON
            custom_budgets.big_category_id = big_categories.id
        WHERE
            custom_budgets.user_id = ?
        AND
            custom_budgets.years_months = ?
        ORDER BY
            custom_budgets.big_category_id`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, yearMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customBudgetByCategoryList []model.CustomBudgetByCategory
	for rows.Next() {
		var customBudgetByCategory model.CustomBudgetByCategory
		if err := rows.StructScan(&customBudgetByCategory); err != nil {
			return nil, err
		}
		customBudgetByCategoryList = append(customBudgetByCategoryList, customBudgetByCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	customBudgets := model.NewCustomBudgets(customBudgetByCategoryList)

	return &customBudgets, nil
}

func (r *BudgetsRepository) PostCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error {
	query := `
        INSERT INTO custom_budgets
            (user_id, years_months, big_category_id, budget)
        VALUES
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?),
            (?,?,?,?)`

	var queryArgs []interface{}
	for _, customBudgetByCategory := range customBudgets.CustomBudgets {
		queryArgs = append(queryArgs, userID, yearMonth, customBudgetByCategory.BigCategoryID, customBudgetByCategory.Budget)
	}

	_, err := r.MySQLHandler.conn.Exec(query, queryArgs...)
	return err
}

func (r *BudgetsRepository) PutCustomBudgets(customBudgets *model.CustomBudgets, yearMonth time.Time, userID string) error {
	for _, customBudgetByCategory := range customBudgets.CustomBudgets {
		query := `
	   UPDATE
	       custom_budgets
	   SET
	       budget = ?
	   WHERE
	       user_id = ?
       AND
           years_months = ?
	   AND
	       big_category_id = ?`

		_, err := r.MySQLHandler.conn.Exec(query, customBudgetByCategory.Budget, userID, yearMonth, customBudgetByCategory.BigCategoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
