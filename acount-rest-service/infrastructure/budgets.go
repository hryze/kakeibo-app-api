package infrastructure

import "github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

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
	for _, budget := range standardBudgets.StandardBudgets {
		query := `
	   UPDATE
	       standard_budgets
	   SET
	       budget = ?
	   WHERE
	       user_id = ?
	   AND
	       big_category_id = ?`

		_, err := r.MySQLHandler.conn.Exec(query, budget.Budget, userID, budget.BigCategoryID)
		if err != nil {
			return err
		}
	}

	return nil
}
