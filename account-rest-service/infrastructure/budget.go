package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type BudgetsRepository struct {
	*MySQLHandler
}

func NewBudgetsRepository(mysqlHandler *MySQLHandler) *BudgetsRepository {
	return &BudgetsRepository{mysqlHandler}
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
	query := `
	   UPDATE
	       standard_budgets
	   SET
	       budget = ?
	   WHERE
	       user_id = ?
	   AND
	       big_category_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, standardBudgetByCategory := range standardBudgets.StandardBudgets {
			if _, err := r.MySQLHandler.conn.Exec(query, standardBudgetByCategory.Budget, userID, standardBudgetByCategory.BigCategoryID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := transactions(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
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

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, customBudgetByCategory := range customBudgets.CustomBudgets {
			if _, err := r.MySQLHandler.conn.Exec(query, customBudgetByCategory.Budget, userID, yearMonth, customBudgetByCategory.BigCategoryID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := transactions(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *BudgetsRepository) DeleteCustomBudgets(yearMonth time.Time, userID string) error {
	query := `
        DELETE 
        FROM
            custom_budgets
        WHERE
            user_id = ?
        AND
            years_months = ?`
	_, err := r.MySQLHandler.conn.Exec(query, userID, yearMonth)
	return err
}

func (r *BudgetsRepository) GetMonthlyStandardBudget(userID string) (model.MonthlyBudget, error) {
	query := `
        SELECT
            SUM(budget) total_budget
        FROM
            standard_budgets
        WHERE
            user_id = ?`

	monthlyStandardBudget := model.MonthlyBudget{BudgetType: "StandardBudget"}
	if err := r.MySQLHandler.conn.QueryRowx(query, userID).StructScan(&monthlyStandardBudget); err != nil {
		return monthlyStandardBudget, err
	}
	return monthlyStandardBudget, nil
}

func (r *BudgetsRepository) GetMonthlyCustomBudgets(year time.Time, userID string) ([]model.MonthlyBudget, error) {
	query := `
        SELECT
            years_months,
            SUM(budget) total_budget
        FROM
            custom_budgets
        WHERE
            user_id = ?
        AND
            years_months >= ?
        AND
            years_months < ?
        GROUP BY
            years_months`

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, year, year.AddDate(1, 0, 0))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyCustomBudgets []model.MonthlyBudget
	for rows.Next() {
		monthlyCustomBudget := model.MonthlyBudget{BudgetType: "CustomBudget"}
		if err := rows.StructScan(&monthlyCustomBudget); err != nil {
			return nil, err
		}
		monthlyCustomBudgets = append(monthlyCustomBudgets, monthlyCustomBudget)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return monthlyCustomBudgets, nil
}
