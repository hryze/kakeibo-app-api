package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type GroupBudgetsRepository struct {
	*MySQLHandler
}

func NewGroupBudgetsRepository(mysqlHandler *MySQLHandler) *GroupBudgetsRepository {
	return &GroupBudgetsRepository{mysqlHandler}
}

func (r *GroupBudgetsRepository) PostInitGroupStandardBudgets(groupID int) error {
	query := `
        INSERT INTO group_standard_budgets
            (group_id, big_category_id)
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

	_, err := r.MySQLHandler.conn.Exec(query, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID)

	return err
}

func (r *GroupBudgetsRepository) GetGroupStandardBudgets(groupID int) (*model.GroupStandardBudgets, error) {
	query := `
        SELECT
            group_standard_budgets.big_category_id big_category_id,
            big_categories.category_name big_category_name,
            group_standard_budgets.budget budget
        FROM
            group_standard_budgets
        INNER JOIN
            big_categories
        ON
            group_standard_budgets.big_category_id = big_categories.id
        WHERE
            group_standard_budgets.group_id = ?
        ORDER BY
            group_standard_budgets.big_category_id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupStandardBudgetByCategoryList []model.GroupStandardBudgetByCategory
	for rows.Next() {
		var groupStandardBudgetByCategory model.GroupStandardBudgetByCategory
		if err := rows.StructScan(&groupStandardBudgetByCategory); err != nil {
			return nil, err
		}

		groupStandardBudgetByCategoryList = append(groupStandardBudgetByCategoryList, groupStandardBudgetByCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	groupStandardBudgets := model.NewGroupStandardBudgets(groupStandardBudgetByCategoryList)

	return &groupStandardBudgets, nil
}

func (r *GroupBudgetsRepository) PutGroupStandardBudgets(groupStandardBudgets *model.GroupStandardBudgets, groupID int) error {
	query := `
	    UPDATE
	        group_standard_budgets
	    SET
	        budget = ?
	    WHERE
	        group_id = ?
	    AND
	        big_category_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, groupStandardBudgetByCategory := range groupStandardBudgets.GroupStandardBudgets {
			if _, err := r.MySQLHandler.conn.Exec(query, groupStandardBudgetByCategory.Budget, groupID, groupStandardBudgetByCategory.BigCategoryID); err != nil {
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

func (r *GroupBudgetsRepository) GetGroupCustomBudgets(yearMonth time.Time, groupID int) (*model.GroupCustomBudgets, error) {
	query := `
        SELECT
            group_custom_budgets.big_category_id big_category_id,
            big_categories.category_name big_category_name,
            group_custom_budgets.budget budget
        FROM
            group_custom_budgets
        INNER JOIN
            big_categories
        ON
            group_custom_budgets.big_category_id = big_categories.id
        WHERE
            group_custom_budgets.group_id = ?
        AND
            group_custom_budgets.years_months = ?
        ORDER BY
            group_custom_budgets.big_category_id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, yearMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupCustomBudgetByCategoryList []model.GroupCustomBudgetByCategory
	for rows.Next() {
		var groupCustomBudgetByCategory model.GroupCustomBudgetByCategory
		if err := rows.StructScan(&groupCustomBudgetByCategory); err != nil {
			return nil, err
		}

		groupCustomBudgetByCategoryList = append(groupCustomBudgetByCategoryList, groupCustomBudgetByCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	groupCustomBudgets := model.NewGroupCustomBudgets(groupCustomBudgetByCategoryList)

	return &groupCustomBudgets, nil
}

func (r *GroupBudgetsRepository) PostGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error {
	query := `
        INSERT INTO group_custom_budgets
            (group_id, years_months, big_category_id, budget)
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
	for _, groupCustomBudgetByCategory := range groupCustomBudgets.GroupCustomBudgets {
		queryArgs = append(queryArgs, groupID, yearMonth, groupCustomBudgetByCategory.BigCategoryID, groupCustomBudgetByCategory.Budget)
	}

	_, err := r.MySQLHandler.conn.Exec(query, queryArgs...)

	return err
}

func (r *GroupBudgetsRepository) PutGroupCustomBudgets(groupCustomBudgets *model.GroupCustomBudgets, yearMonth time.Time, groupID int) error {
	query := `
	    UPDATE
	        group_custom_budgets
	    SET
	        budget = ?
	    WHERE
	        group_id = ?
        AND
            years_months = ?
	    AND
	        big_category_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, groupCustomBudgetByCategory := range groupCustomBudgets.GroupCustomBudgets {
			if _, err := r.MySQLHandler.conn.Exec(query, groupCustomBudgetByCategory.Budget, groupID, yearMonth, groupCustomBudgetByCategory.BigCategoryID); err != nil {
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

func (r *GroupBudgetsRepository) DeleteGroupCustomBudgets(yearMonth time.Time, groupID int) error {
	query := `
        DELETE 
        FROM
            group_custom_budgets
        WHERE
            group_id = ?
        AND
            years_months = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupID, yearMonth)

	return err
}

func (r *GroupBudgetsRepository) GetMonthlyGroupStandardBudget(groupID int) (model.MonthlyGroupBudget, error) {
	query := `
        SELECT
            SUM(budget) total_budget
        FROM
            group_standard_budgets
        WHERE
            group_id = ?`

	monthlyGroupStandardBudget := model.MonthlyGroupBudget{BudgetType: "StandardBudget"}
	if err := r.MySQLHandler.conn.QueryRowx(query, groupID).StructScan(&monthlyGroupStandardBudget); err != nil {
		return monthlyGroupStandardBudget, err
	}

	return monthlyGroupStandardBudget, nil
}

func (r *GroupBudgetsRepository) GetMonthlyGroupCustomBudgets(year time.Time, groupID int) ([]model.MonthlyGroupBudget, error) {
	query := `
        SELECT
            years_months,
            SUM(budget) total_budget
        FROM
            group_custom_budgets
        WHERE
            group_id = ?
        AND
            years_months >= ?
        AND
            years_months < ?
        GROUP BY
            years_months`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, year, year.AddDate(1, 0, 0))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyGroupCustomBudgets []model.MonthlyGroupBudget
	for rows.Next() {
		monthlyGroupCustomBudget := model.MonthlyGroupBudget{BudgetType: "CustomBudget"}
		if err := rows.StructScan(&monthlyGroupCustomBudget); err != nil {
			return nil, err
		}

		monthlyGroupCustomBudgets = append(monthlyGroupCustomBudgets, monthlyGroupCustomBudget)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return monthlyGroupCustomBudgets, nil
}
