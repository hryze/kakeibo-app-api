package infrastructure

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type ShoppingListRepository struct {
	*MySQLHandler
}

func NewShoppingListRepository(mysqlHandler *MySQLHandler) *ShoppingListRepository {
	return &ShoppingListRepository{mysqlHandler}
}

func (r *ShoppingListRepository) GetRegularShoppingList(userID string) (model.RegularShoppingList, error) {
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
            transaction_auto_add
        FROM
            regular_shopping_list
        WHERE
            user_id = ?`

	regularShoppingList := model.RegularShoppingList{
		RegularShoppingList: make([]model.RegularShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID)
	if err != nil {
		return regularShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var regularShoppingItem model.RegularShoppingItem
		if err := rows.StructScan(&regularShoppingItem); err != nil {
			return regularShoppingList, err
		}

		regularShoppingList.RegularShoppingList = append(regularShoppingList.RegularShoppingList, regularShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return regularShoppingList, err
	}

	return regularShoppingList, nil
}

func (r *ShoppingListRepository) GetRegularShoppingItem(regularShoppingItemID int) (model.RegularShoppingItem, error) {
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
            transaction_auto_add
        FROM
            regular_shopping_list
        WHERE
            id = ?`

	var regularShoppingItem model.RegularShoppingItem
	if err := r.MySQLHandler.conn.QueryRowx(query, regularShoppingItemID).StructScan(&regularShoppingItem); err != nil {
		return regularShoppingItem, err
	}

	return regularShoppingItem, nil
}

func (r *ShoppingListRepository) GetShoppingListRelatedToPostedRegularShoppingItem(todayShoppingItemID int, laterThanTodayShoppingItemID int) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            id = ?`

	var queryArgs []interface{}
	if todayShoppingItemID != 0 {
		query += " UNION " + query
		queryArgs = append(queryArgs, todayShoppingItemID, laterThanTodayShoppingItemID)
	} else {
		queryArgs = append(queryArgs, laterThanTodayShoppingItemID)
	}

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, queryArgs...)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) PostRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, userID string, today time.Time) (sql.Result, sql.Result, sql.Result, error) {
	regularShoppingItemQuery := `
        INSERT INTO regular_shopping_list
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
            user_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?,?)`

	shoppingItemQuery := `
        INSERT INTO shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            user_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return nil, nil, nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, sql.Result, sql.Result, error) {
		var regularShoppingItemResult, todayShoppingItemResult, laterThanTodayShoppingItemResult sql.Result
		nextExpectedPurchaseDate := regularShoppingItem.ExpectedPurchaseDate.Time

		if today.Equal(regularShoppingItem.ExpectedPurchaseDate.Time) {
			if regularShoppingItem.CycleType == "daily" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
			} else if regularShoppingItem.CycleType == "weekly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
			} else if regularShoppingItem.CycleType == "monthly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
			} else if regularShoppingItem.CycleType == "custom" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, regularShoppingItem.Cycle.Int)
			}
		}

		regularShoppingItemResult, err = tx.Exec(
			regularShoppingItemQuery,
			nextExpectedPurchaseDate,
			regularShoppingItem.CycleType,
			regularShoppingItem.Cycle,
			regularShoppingItem.Purchase,
			regularShoppingItem.Shop,
			regularShoppingItem.Amount,
			regularShoppingItem.BigCategoryID,
			regularShoppingItem.MediumCategoryID,
			regularShoppingItem.CustomCategoryID,
			userID,
			regularShoppingItem.TransactionAutoAdd,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		regularShoppingItemId, err := regularShoppingItemResult.LastInsertId()
		if err != nil {
			return nil, nil, nil, err
		}

		laterThanTodayShoppingItemResult, err = tx.Exec(
			shoppingItemQuery,
			nextExpectedPurchaseDate,
			regularShoppingItem.Purchase,
			regularShoppingItem.Shop,
			regularShoppingItem.Amount,
			regularShoppingItem.BigCategoryID,
			regularShoppingItem.MediumCategoryID,
			regularShoppingItem.CustomCategoryID,
			regularShoppingItemId,
			userID,
			regularShoppingItem.TransactionAutoAdd,
		)
		if err != nil {
			return nil, nil, nil, err
		}

		if today.Equal(regularShoppingItem.ExpectedPurchaseDate.Time) {
			todayShoppingItemResult, err = tx.Exec(
				shoppingItemQuery,
				today,
				regularShoppingItem.Purchase,
				regularShoppingItem.Shop,
				regularShoppingItem.Amount,
				regularShoppingItem.BigCategoryID,
				regularShoppingItem.MediumCategoryID,
				regularShoppingItem.CustomCategoryID,
				regularShoppingItemId,
				userID,
				regularShoppingItem.TransactionAutoAdd,
			)
			if err != nil {
				return nil, nil, nil, err
			}
		}

		return regularShoppingItemResult, todayShoppingItemResult, laterThanTodayShoppingItemResult, nil
	}

	regularShoppingItemResult, todayShoppingItemResult, laterThanTodayShoppingItemResult, err := transactions(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, nil, nil, err
		}

		return nil, nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, nil, err
	}

	return regularShoppingItemResult, todayShoppingItemResult, laterThanTodayShoppingItemResult, nil
}

func (r *ShoppingListRepository) GetShoppingListRelatedToUpdatedRegularShoppingItem(regularShoppingItemID int) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            regular_shopping_list_id = ?
        ORDER BY
            expected_purchase_date`

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, regularShoppingItemID)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) PutRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, regularShoppingItemID int, userID string, today time.Time) error {
	deleteShoppingItemQuery := `
        DELETE
        FROM
            shopping_list
        WHERE
            regular_shopping_list_id = ?
        AND
            complete_flag = false`

	updateRegularShoppingItemQuery := `
        UPDATE
            regular_shopping_list
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
            transaction_auto_add = ?
        WHERE
            id = ?`

	insertShoppingItemQuery := `
        INSERT INTO 
            shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            user_id,
            transaction_auto_add
        )
        VALUES
        (
            ?,?,?,?,?,?,?,?,?,?
        )`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(
			deleteShoppingItemQuery,
			regularShoppingItemID,
		); err != nil {
			return err
		}

		if _, err = tx.Exec(
			insertShoppingItemQuery,
			regularShoppingItem.ExpectedPurchaseDate,
			regularShoppingItem.Purchase,
			regularShoppingItem.Shop,
			regularShoppingItem.Amount,
			regularShoppingItem.BigCategoryID,
			regularShoppingItem.MediumCategoryID,
			regularShoppingItem.CustomCategoryID,
			regularShoppingItemID,
			userID,
			regularShoppingItem.TransactionAutoAdd,
		); err != nil {
			return err
		}

		nextExpectedPurchaseDate := regularShoppingItem.ExpectedPurchaseDate.Time

		for !today.Before(nextExpectedPurchaseDate) {
			if regularShoppingItem.CycleType == "daily" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
			} else if regularShoppingItem.CycleType == "weekly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
			} else if regularShoppingItem.CycleType == "monthly" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
			} else if regularShoppingItem.CycleType == "custom" {
				nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, regularShoppingItem.Cycle.Int)
			}

			if _, err = tx.Exec(
				insertShoppingItemQuery,
				nextExpectedPurchaseDate,
				regularShoppingItem.Purchase,
				regularShoppingItem.Shop,
				regularShoppingItem.Amount,
				regularShoppingItem.BigCategoryID,
				regularShoppingItem.MediumCategoryID,
				regularShoppingItem.CustomCategoryID,
				regularShoppingItemID,
				userID,
				regularShoppingItem.TransactionAutoAdd,
			); err != nil {
				fmt.Println(err)
				return err
			}
		}

		if _, err := tx.Exec(
			updateRegularShoppingItemQuery,
			nextExpectedPurchaseDate,
			regularShoppingItem.CycleType,
			regularShoppingItem.Cycle,
			regularShoppingItem.Purchase,
			regularShoppingItem.Shop,
			regularShoppingItem.Amount,
			regularShoppingItem.BigCategoryID,
			regularShoppingItem.MediumCategoryID,
			regularShoppingItem.CustomCategoryID,
			regularShoppingItem.TransactionAutoAdd,
			regularShoppingItemID,
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

func (r *ShoppingListRepository) PutRegularShoppingList(regularShoppingList model.RegularShoppingList, userID string, today time.Time) error {
	updateRegularShoppingItemQuery := `
        UPDATE
            regular_shopping_list
        SET
            expected_purchase_date = ?
        WHERE
            id = ?`

	insertShoppingItemQuery := `
        INSERT INTO 
            shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            regular_shopping_list_id,
            user_id,
            transaction_auto_add
        )
        VALUES
        (
            ?,?,?,?,?,?,?,?,?,?
        )`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, regularShoppingItem := range regularShoppingList.RegularShoppingList {
			nextExpectedPurchaseDate := regularShoppingItem.ExpectedPurchaseDate.Time

			for !today.Before(nextExpectedPurchaseDate) {
				if regularShoppingItem.CycleType == "daily" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 1)
				} else if regularShoppingItem.CycleType == "weekly" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, 7)
				} else if regularShoppingItem.CycleType == "monthly" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 1, 0)
				} else if regularShoppingItem.CycleType == "custom" {
					nextExpectedPurchaseDate = nextExpectedPurchaseDate.AddDate(0, 0, regularShoppingItem.Cycle.Int)
				}

				if _, err = tx.Exec(
					insertShoppingItemQuery,
					nextExpectedPurchaseDate,
					regularShoppingItem.Purchase,
					regularShoppingItem.Shop,
					regularShoppingItem.Amount,
					regularShoppingItem.BigCategoryID,
					regularShoppingItem.MediumCategoryID,
					regularShoppingItem.CustomCategoryID,
					regularShoppingItem.ID,
					userID,
					regularShoppingItem.TransactionAutoAdd,
				); err != nil {
					return err
				}
			}

			if !today.Before(regularShoppingItem.ExpectedPurchaseDate.Time) {
				if _, err := tx.Exec(
					updateRegularShoppingItemQuery,
					nextExpectedPurchaseDate,
					regularShoppingItem.ID,
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

func (r *ShoppingListRepository) DeleteRegularShoppingItem(regularShoppingItemID int) error {
	deleteShoppingItemQuery := `
        DELETE
        FROM
            shopping_list
        WHERE
            regular_shopping_list_id = ?
        AND
            complete_flag = false`

	deleteRegularShoppingItemQuery := `
        DELETE
        FROM
            regular_shopping_list
        WHERE
            id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		if _, err := tx.Exec(deleteShoppingItemQuery, regularShoppingItemID); err != nil {
			return err
		}

		if _, err := tx.Exec(deleteRegularShoppingItemQuery, regularShoppingItemID); err != nil {
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

func (r *ShoppingListRepository) GetDailyShoppingListByDay(date time.Time, userID string) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            user_id = ?
        AND
            expected_purchase_date = ?`

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) GetDailyShoppingListByCategory(date time.Time, userID string) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            user_id = ?
        AND
            expected_purchase_date = ?
        ORDER BY
            big_category_id`

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, date)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) GetMonthlyShoppingListByDay(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            user_id = ?
        AND
            expected_purchase_date >= ?
        AND
            expected_purchase_date <= ?
        ORDER BY
            expected_purchase_date`

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) GetMonthlyShoppingListByCategory(firstDay time.Time, lastDay time.Time, userID string) (model.ShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            user_id = ?
        AND
            expected_purchase_date >= ?
        AND
            expected_purchase_date <= ?
        ORDER BY
            big_category_id, expected_purchase_date`

	shoppingList := model.ShoppingList{
		ShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return shoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var shoppingItem model.ShoppingItem
		if err := rows.StructScan(&shoppingItem); err != nil {
			return shoppingList, err
		}

		shoppingList.ShoppingList = append(shoppingList.ShoppingList, shoppingItem)
	}

	if err := rows.Err(); err != nil {
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *ShoppingListRepository) GetExpiredShoppingList(dueDate time.Time, userID string) (model.ExpiredShoppingList, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            user_id = ?
        AND
            complete_flag = false
        AND
            expected_purchase_date <= ?
        ORDER BY
            expected_purchase_date`

	expiredShoppingList := model.ExpiredShoppingList{
		ExpiredShoppingList: make([]model.ShoppingItem, 0),
	}

	rows, err := r.MySQLHandler.conn.Queryx(query, userID, dueDate)
	if err != nil {
		return expiredShoppingList, err
	}
	defer rows.Close()

	for rows.Next() {
		var expiredShoppingItem model.ShoppingItem
		if err := rows.StructScan(&expiredShoppingItem); err != nil {
			return expiredShoppingList, err
		}

		expiredShoppingList.ExpiredShoppingList = append(expiredShoppingList.ExpiredShoppingList, expiredShoppingItem)
	}

	if err := rows.Err(); err != nil {
		return expiredShoppingList, err
	}

	return expiredShoppingList, nil
}

func (r *ShoppingListRepository) GetShoppingItem(shoppingItemID int) (model.ShoppingItem, error) {
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
            transaction_auto_add,
            transaction_id
        FROM
            shopping_list
        WHERE
            id = ?`

	var shoppingItem model.ShoppingItem
	if err := r.MySQLHandler.conn.QueryRowx(query, shoppingItemID).StructScan(&shoppingItem); err != nil {
		return shoppingItem, err
	}

	return shoppingItem, nil
}

func (r *ShoppingListRepository) PostShoppingItem(shoppingItem *model.ShoppingItem, userID string) (sql.Result, error) {
	query := `
        INSERT INTO shopping_list
        (
            expected_purchase_date,
            purchase,
            shop,
            amount,
            big_category_id,
            medium_category_id,
            custom_category_id,
            user_id,
            transaction_auto_add
        )
        VALUES
            (?,?,?,?,?,?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(
		query,
		shoppingItem.ExpectedPurchaseDate,
		shoppingItem.Purchase,
		shoppingItem.Shop,
		shoppingItem.Amount,
		shoppingItem.BigCategoryID,
		shoppingItem.MediumCategoryID,
		shoppingItem.CustomCategoryID,
		userID,
		shoppingItem.TransactionAutoAdd,
	)

	return result, err
}

func (r *ShoppingListRepository) PutShoppingItem(shoppingItem *model.ShoppingItem) (sql.Result, error) {
	query := `
        UPDATE
            shopping_list
        SET 
            expected_purchase_date = ?,
            complete_flag = ?,
            purchase = ?,
            shop = ?,
            amount = ?,
            big_category_id = ?,
            medium_category_id = ?,
            custom_category_id = ?,
            regular_shopping_list_id = ?,
            transaction_auto_add = ?,
            transaction_id = ?
        WHERE
            id = ?`

	relatedTransactionID := func(relatedTransactionData *model.TransactionData) *int64 {
		if relatedTransactionData != nil {
			return &relatedTransactionData.ID.Int64
		}

		return nil
	}

	result, err := r.MySQLHandler.conn.Exec(
		query,
		shoppingItem.ExpectedPurchaseDate,
		shoppingItem.CompleteFlag,
		shoppingItem.Purchase,
		shoppingItem.Shop,
		shoppingItem.Amount,
		shoppingItem.BigCategoryID,
		shoppingItem.MediumCategoryID,
		shoppingItem.CustomCategoryID,
		shoppingItem.RegularShoppingListID,
		shoppingItem.TransactionAutoAdd,
		relatedTransactionID(shoppingItem.RelatedTransactionData),
		shoppingItem.ID,
	)

	return result, err
}

func (r *ShoppingListRepository) DeleteShoppingItem(shoppingItemID int) error {
	query := `
        DELETE
        FROM
            shopping_list
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, shoppingItemID)

	return err
}

func (r *ShoppingListRepository) PutShoppingListCustomCategoryIdToMediumCategoryId(mediumCategoryID int, customCategoryID int) error {
	updateShoppingListQuery := `
        UPDATE
            shopping_list
        SET 
            medium_category_id = ?,
            custom_category_id = NULL
        WHERE
            custom_category_id = ?`

	updateRegularShoppingListQuery := `
        UPDATE
            regular_shopping_list
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
			updateShoppingListQuery,
			mediumCategoryID,
			customCategoryID,
		); err != nil {
			return err
		}

		if _, err := tx.Exec(
			updateRegularShoppingListQuery,
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
