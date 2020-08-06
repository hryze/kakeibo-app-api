package infrastructure

import "github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"

type GroupBudgetsRepository struct {
	*MySQLHandler
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
        LEFT JOIN
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
