package infrastructure

import (
	"database/sql"
	"time"

	"github.com/hryze/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupShoppingListRepository struct {
	*MySQLHandler
}

func NewGroupShoppingListRepository(mysqlHandler *MySQLHandler) *GroupShoppingListRepository {
	return &GroupShoppingListRepository{mysqlHandler}
}

func (r *GroupShoppingListRepository) GetGroupRegularShoppingList(groupID int) (model.GroupRegularShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            cycle_type,
            cycle,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            payment_user_id,
            transaction_auto_add
        FROM
            group_regular_shopping_list
        WHERE
            group_id = ?`

	groupRegularShoppingList := model.GroupRegularShoppingList{
		GroupRegularShoppingList: make([]model.GroupRegularShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return groupRegularShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupRegularShoppingItem model.GroupRegularShoppingItem
		if err := rows.StructScan(&groupRegularShoppingItem); err != nil {
			return groupRegularShoppingList, err
		}

		groupRegularShoppingList.GroupRegularShoppingList = append(groupRegularShoppingList.GroupRegularShoppingList, groupRegularShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupRegularShoppingList, err
	}

	return groupRegularShoppingList, nil
}

func (r *GroupShoppingListRepository) GetGroupRegularShoppingItem(groupRegularShoppingItemID int) (model.GroupRegularShoppingItem, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            cycle_type,
            cycle,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            payment_user_id,
            transaction_auto_add
        FROM
            group_regular_shopping_list
        WHERE
            id = ?`

	var groupRegularShoppingItem model.GroupRegularShoppingItem
	if err := r.MySQLHandler.conn.QueryRowx(query, groupRegularShoppingItemID).StructScan(&groupRegularShoppingItem); err != nil {
		return groupRegularShoppingItem, err
	}

	return groupRegularShoppingItem, nil
}

func (r *GroupShoppingListRepository) GetGroupShoppingListRelatedToPostedGroupRegularShoppingItem(todayGroupShoppingItemID int, laterThanTodayGroupShoppingItemID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            id = ?`

	var queryArgs []interface{}
	if todayGroupShoppingItemID != 0 {
		query += " UNION " + query
		queryArgs = append(queryArgs, todayGroupShoppingItemID, laterThanTodayGroupShoppingItemID)
	} else {
		queryArgs = append(queryArgs, laterThanTodayGroupShoppingItemID)
	}

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, queryArgs...)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) PostGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupID int, today time.Time) (sql.Result, sql.Result, sql.Result, error) {
	groupRegularShoppingItemQuery := `
        INSERT INTO group_regular_shopping_list
        (
            expected_purchase_date,
            cycle_type,
            cycle,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            payment_user_id,
            group_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?,?,?)`

	groupShoppingItemQuery := `
        INSERT INTO group_shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            group_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return nil, nil, nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, sql.Result, sql.Result, error) {
		var groupRegularShoppingItemResult, todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult sql.Result
		nextExpectedPurchaseDate := groupRegularShoppingItem.ExpectedPurchaseDate.Time

		if today.Equal(groupRegularShoppingItem.ExpectedPurchaseDate.Time) {
			if groupRegularShoppingItem.CycleType == "daily" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
			} else if groupRegularShoppingItem.CycleType == "weekly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
			} else if groupRegularShoppingItem.CycleType == "monthly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
			} else if groupRegularShoppingItem.CycleType == "custom" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, groupRegularShoppingItem.Cycle.Int)
			}
		}

		groupRegularShoppingItemResult, err = tx.Exec(
			groupRegularShoppingItemQuery,
			nextExpectedPurchaseDate,
			groupRegularShoppingItem.CycleType,
			groupRegularShoppingItem.Cycle,
			groupRegularShoppingItem.Purchase,
			groupRegularShoppingItem.Shop,
			groupRegularShoppingItem.Amount,
			groupRegularShoppingItem.BigCategoryID,
			groupRegularShoppingItem.MediumCategoryID,
			groupRegularShoppingItem.CustomCategoryID,
			groupRegularShoppingItem.PaymentUserID,
			groupID,
			groupRegularShoppingItem.TransactionAutoAdd,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		groupRegularShoppingItemId, err := groupRegularShoppingItemResult.LastInsertId()
		if err != nil {
			return nil, nil, nil, err
		}

		laterThanTodayGroupShoppingItemResult, err = tx.Exec(
			groupShoppingItemQuery,
			nextExpectedPurchaseDate,
			groupRegularShoppingItem.Purchase,
			groupRegularShoppingItem.Shop,
			groupRegularShoppingItem.Amount,
			groupRegularShoppingItem.BigCategoryID,
			groupRegularShoppingItem.MediumCategoryID,
			groupRegularShoppingItem.CustomCategoryID,
			groupRegularShoppingItemId,
			groupRegularShoppingItem.PaymentUserID,
			groupID,
			groupRegularShoppingItem.TransactionAutoAdd,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		if today.Equal(groupRegularShoppingItem.ExpectedPurchaseDate.Time) {
			todayGroupShoppingItemResult, err = tx.Exec(
				groupShoppingItemQuery,
				today,
				groupRegularShoppingItem.Purchase,
				groupRegularShoppingItem.Shop,
				groupRegularShoppingItem.Amount,
				groupRegularShoppingItem.BigCategoryID,
				groupRegularShoppingItem.MediumCategoryID,
				groupRegularShoppingItem.CustomCategoryID,
				groupRegularShoppingItemId,
				groupRegularShoppingItem.PaymentUserID,
				groupID,
				groupRegularShoppingItem.TransactionAutoAdd,
			)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		return groupRegularShoppingItemResult, todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult, nil
	}

	groupRegularShoppingItemResult, todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult, err := transactions(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, nil, nil, err
		}

		return nil, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	return groupRegularShoppingItemResult, todayGroupShoppingItemResult, laterThanTodayGroupShoppingItemResult, nil
}

func (r *GroupShoppingListRepository) GetGroupShoppingListRelatedToUpdatedGroupRegularShoppingItem(groupRegularShoppingItemID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            regular_shopping_list_id = ?
        ORDER BY
            expected_purchase_date`

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupRegularShoppingItemID)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) PutGroupRegularShoppingItem(groupRegularShoppingItem *model.GroupRegularShoppingItem, groupRegularShoppingItemID int, groupID int, today time.Time) error {
	deleteGroupShoppingItemQuery := `
        DELETE
        FROM
            group_shopping_list
        WHERE
            regular_shopping_list_id = ?
        AND
            complete_flag = false`

	updateGroupRegularShoppingItemQuery := `
        UPDATE
            group_regular_shopping_list
        SET
            expected_purchase_date = ?,
            cycle_type = ?,
            cycle = ?,
            purchase = ?,
            shop = ?,
            amount = ?,
            big_category_id = ?,
            medium_category_id = ?,
            custom_category_id = ?,
            payment_user_id = ?,
            transaction_auto_add = ?
        WHERE
            id = ?`

	insertGroupShoppingItemQuery := `
        INSERT INTO group_shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            group_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			deleteGroupShoppingItemQuery,
			groupRegularShoppingItemID,
		); err != nil {
			return err
		}

		if _, err = tx.Exec(
			insertGroupShoppingItemQuery,
			groupRegularShoppingItem.ExpectedPurchaseDate,
			groupRegularShoppingItem.Purchase,
			groupRegularShoppingItem.Shop,
			groupRegularShoppingItem.Amount,
			groupRegularShoppingItem.BigCategoryID,
			groupRegularShoppingItem.MediumCategoryID,
			groupRegularShoppingItem.CustomCategoryID,
			groupRegularShoppingItemID,
			groupRegularShoppingItem.PaymentUserID,
			groupID,
			groupRegularShoppingItem.TransactionAutoAdd,
		); err != nil {
			return err
		}

		nextExpectedPurchaseDate := groupRegularShoppingItem.ExpectedPurchaseDate.Time

		for !today.Before(nextExpectedPurchaseDate) {
			if groupRegularShoppingItem.CycleType == "daily" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
			} else if groupRegularShoppingItem.CycleType == "weekly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
			} else if groupRegularShoppingItem.CycleType == "monthly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
			} else if groupRegularShoppingItem.CycleType == "custom" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, groupRegularShoppingItem.Cycle.Int)
			}

			if _, err = tx.Exec(
				insertGroupShoppingItemQuery,
				nextExpectedPurchaseDate,
				groupRegularShoppingItem.Purchase,
				groupRegularShoppingItem.Shop,
				groupRegularShoppingItem.Amount,
				groupRegularShoppingItem.BigCategoryID,
				groupRegularShoppingItem.MediumCategoryID,
				groupRegularShoppingItem.CustomCategoryID,
				groupRegularShoppingItemID,
				groupRegularShoppingItem.PaymentUserID,
				groupID,
				groupRegularShoppingItem.TransactionAutoAdd,
			); err != nil {
				return err
			}
		}

		if _, err := tx.Exec(
			updateGroupRegularShoppingItemQuery,
			nextExpectedPurchaseDate,
			groupRegularShoppingItem.CycleType,
			groupRegularShoppingItem.Cycle,
			groupRegularShoppingItem.Purchase,
			groupRegularShoppingItem.Shop,
			groupRegularShoppingItem.Amount,
			groupRegularShoppingItem.BigCategoryID,
			groupRegularShoppingItem.MediumCategoryID,
			groupRegularShoppingItem.CustomCategoryID,
			groupRegularShoppingItem.PaymentUserID,
			groupRegularShoppingItem.TransactionAutoAdd,
			groupRegularShoppingItemID,
		); err != nil {
			return err
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

func (r *GroupShoppingListRepository) PutGroupRegularShoppingList(groupRegularShoppingList model.GroupRegularShoppingList, groupID int, today time.Time) error {
	updateGroupRegularShoppingItemQuery := `
        UPDATE
            group_regular_shopping_list
        SET
            expected_purchase_date = ?
        WHERE
            id = ?`

	insertGroupShoppingItemQuery := `
        INSERT INTO 
            group_shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            group_id,
            transaction_auto_add
        )
        VALUES
        (
            ?,?,?,?,?,?,?,?,?,?,?
        )`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, groupRegularShoppingItem := range groupRegularShoppingList.GroupRegularShoppingList {
			nextExpectedPurchaseDate := groupRegularShoppingItem.ExpectedPurchaseDate.Time

			for !today.Before(nextExpectedPurchaseDate) {
				if groupRegularShoppingItem.CycleType == "daily" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
				} else if groupRegularShoppingItem.CycleType == "weekly" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
				} else if groupRegularShoppingItem.CycleType == "monthly" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
				} else if groupRegularShoppingItem.CycleType == "custom" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, groupRegularShoppingItem.Cycle.Int)
				}

				if _, err = tx.Exec(
					insertGroupShoppingItemQuery,
					nextExpectedPurchaseDate,
					groupRegularShoppingItem.Purchase,
					groupRegularShoppingItem.Shop,
					groupRegularShoppingItem.Amount,
					groupRegularShoppingItem.BigCategoryID,
					groupRegularShoppingItem.MediumCategoryID,
					groupRegularShoppingItem.CustomCategoryID,
					groupRegularShoppingItem.ID,
					groupRegularShoppingItem.PaymentUserID,
					groupID,
					groupRegularShoppingItem.TransactionAutoAdd,
				); err != nil {
					return err
				}
			}

			if !today.Before(groupRegularShoppingItem.ExpectedPurchaseDate.Time) {
				if _, err := tx.Exec(
					updateGroupRegularShoppingItemQuery,
					nextExpectedPurchaseDate,
					groupRegularShoppingItem.ID,
				); err != nil {
					return err
				}
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

func (r *GroupShoppingListRepository) DeleteGroupRegularShoppingItem(groupRegularShoppingItemID int) error {
	deleteGroupShoppingItemQuery := `
        DELETE
        FROM
            group_shopping_list
        WHERE
            regular_shopping_list_id = ?
        AND
            complete_flag = false`

	deleteGroupRegularShoppingItemQuery := `
        DELETE
        FROM
            group_regular_shopping_list
        WHERE
            id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteGroupShoppingItemQuery, groupRegularShoppingItemID); err != nil {
			return err
		}

		if _, err := tx.Exec(deleteGroupRegularShoppingItemQuery, groupRegularShoppingItemID); err != nil {
			return err
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

func (r *GroupShoppingListRepository) GetDailyGroupShoppingListByDay(date time.Time, groupID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            group_id = ?
        AND
            expected_purchase_date = ?`

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) GetDailyGroupShoppingListByCategory(date time.Time, groupID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            group_id = ?
        AND
            expected_purchase_date = ?
        ORDER BY
            big_category_id`

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, date)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) GetMonthlyGroupShoppingListByDay(firstDay time.Time, lastDay time.Time, groupID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            group_id = ?
        AND
            expected_purchase_date >= ?
        AND
            expected_purchase_date <= ?
        ORDER BY
            expected_purchase_date`

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) GetMonthlyGroupShoppingListByCategory(firstDay time.Time, lastDay time.Time, groupID int) (model.GroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            group_id = ?
        AND
            expected_purchase_date >= ?
        AND
            expected_purchase_date <= ?
        ORDER BY
            big_category_id, expected_purchase_date`

	groupShoppingList := model.GroupShoppingList{
		GroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return groupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var groupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&groupShoppingItem); err != nil {
			return groupShoppingList, err
		}

		groupShoppingList.GroupShoppingList = append(groupShoppingList.GroupShoppingList, groupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return groupShoppingList, err
	}

	return groupShoppingList, nil
}

func (r *GroupShoppingListRepository) GetExpiredGroupShoppingList(dueDate time.Time, groupID int) (model.ExpiredGroupShoppingList, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            group_id = ?
        AND
            complete_flag = false
        AND
            expected_purchase_date <= ?
        ORDER BY
            expected_purchase_date`

	expiredGroupShoppingList := model.ExpiredGroupShoppingList{
		ExpiredGroupShoppingList: make([]model.GroupShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, dueDate)
	if err != nil {
		return expiredGroupShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var expiredGroupShoppingItem model.GroupShoppingItem
		if err := rows.StructScan(&expiredGroupShoppingItem); err != nil {
			return expiredGroupShoppingList, err
		}

		expiredGroupShoppingList.ExpiredGroupShoppingList = append(expiredGroupShoppingList.ExpiredGroupShoppingList, expiredGroupShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return expiredGroupShoppingList, err
	}

	return expiredGroupShoppingList, nil
}

func (r *GroupShoppingListRepository) GetGroupShoppingItem(groupShoppingItemID int) (model.GroupShoppingItem, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_flag,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            payment_user_id,
            transaction_auto_add,
            transaction_id
        FROM
            group_shopping_list
        WHERE
            id = ?`

	var groupShoppingItem model.GroupShoppingItem
	if err := r.MySQLHandler.conn.QueryRowx(query, groupShoppingItemID).StructScan(&groupShoppingItem); err != nil {
		return groupShoppingItem, err
	}

	return groupShoppingItem, nil
}

func (r *GroupShoppingListRepository) PostGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem, groupID int) (sql.Result, error) {
	query := `
        INSERT INTO group_shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            payment_user_id,
            group_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(
		query,
		groupShoppingItem.ExpectedPurchaseDate,
		groupShoppingItem.Purchase,
		groupShoppingItem.Shop,
		groupShoppingItem.Amount,
		groupShoppingItem.BigCategoryID,
		groupShoppingItem.MediumCategoryID,
		groupShoppingItem.CustomCategoryID,
		groupShoppingItem.PaymentUserID,
		groupID,
		groupShoppingItem.TransactionAutoAdd,
	)

	return result, err
}

func (r *GroupShoppingListRepository) PutGroupShoppingItem(groupShoppingItem *model.GroupShoppingItem) (sql.Result, error) {
	query := `
        UPDATE
            group_shopping_list
        SET 
            expected_purchase_date = ?,
            complete_flag = ?,
            purchase = ?,
            shop = ?,
            amount = ?,
            big_category_id = ?,
            medium_category_id = ?,
            custom_category_id = ?,
            payment_user_id = ?,
            transaction_auto_add = ?,
            transaction_id = ?
        WHERE
            id = ?`

	relatedTransactionID := func(relatedTransactionData *model.GroupTransactionData) *int64 {
		if relatedTransactionData != nil {
			return &relatedTransactionData.ID.Int64
		}

		return nil
	}

	result, err := r.MySQLHandler.conn.Exec(
		query,
		groupShoppingItem.ExpectedPurchaseDate,
		groupShoppingItem.CompleteFlag,
		groupShoppingItem.Purchase,
		groupShoppingItem.Shop,
		groupShoppingItem.Amount,
		groupShoppingItem.BigCategoryID,
		groupShoppingItem.MediumCategoryID,
		groupShoppingItem.CustomCategoryID,
		groupShoppingItem.PaymentUserID,
		groupShoppingItem.TransactionAutoAdd,
		relatedTransactionID(groupShoppingItem.RelatedTransactionData),
		groupShoppingItem.ID,
	)

	return result, err
}

func (r *GroupShoppingListRepository) DeleteGroupShoppingItem(groupShoppingItemID int) error {
	query := `
        DELETE
        FROM
            group_shopping_list
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupShoppingItemID)

	return err
}

func (r *GroupShoppingListRepository) PutGroupShoppingListCustomCategoryIdToMediumCategoryId(mediumCategoryID int, customCategoryID int) error {
	updateGroupShoppingListQuery := `
        UPDATE
            group_shopping_list
        SET 
            medium_category_id = ?,
            custom_category_id = NULL
        WHERE
            custom_category_id = ?`

	updateGroupRegularShoppingListQuery := `
        UPDATE
            group_regular_shopping_list
        SET 
            medium_category_id = ?,
            custom_category_id = NULL
        WHERE
            custom_category_id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			updateGroupShoppingListQuery,
			mediumCategoryID,
			customCategoryID,
		); err != nil {
			return err
		}

		if _, err := tx.Exec(
			updateGroupRegularShoppingListQuery,
			mediumCategoryID,
			customCategoryID,
		); err != nil {
			return err
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
