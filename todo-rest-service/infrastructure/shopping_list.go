package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type ShoppingListRepository struct {
	*MySQLHandler
}

func NewShoppingListRepository(mysqlHandler *MySQLHandler) *ShoppingListRepository {
	return &ShoppingListRepository{mysqlHandler}
}

func (r *ShoppingListRepository) GetRegularShoppingItem(regularShoppingItemID int64) (model.RegularShoppingItem, error) {
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

func (r *ShoppingListRepository) GetShoppingListRelatedToRegularShoppingItemID(regularShoppingItemID int64) (model.ShoppingList, error) {
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
            regular_shopping_list_id = ?`

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

func (r *ShoppingListRepository) PostRegularShoppingItem(regularShoppingItem *model.RegularShoppingItem, userID string, today time.Time) (sql.Result, error) {
	regularShoppingItemQuery := `
        INSERT INTO regular_shopping_list
            (expected_purchase_date, cycle_type, cycle, purchase, shop, amount, big_category_id, medium_category_id, custom_category_id, user_id, transaction_auto_add)
        VALUES
            (?,?,?,?,?,?,?,?,?,?,?)`

	shoppingItemQuery := `
        INSERT INTO shopping_list
            (expected_purchase_date, purchase, shop, amount, big_category_id, medium_category_id, custom_category_id, regular_shopping_list_id, user_id, transaction_auto_add)
        VALUES
            (?,?,?,?,?,?,?,?,?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return nil, err
	}

	transactions := func(tx *sql.Tx) (sql.Result, error) {
		nextExpectedPurchaseDate := regularShoppingItem.ExpectedPurchaseDate.Time

		if !today.Before(regularShoppingItem.ExpectedPurchaseDate.Time) {
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

		result, err := tx.Exec(
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
			return nil, err
		}

		regularShoppingItemLastInsertId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		if _, err := tx.Exec(
			shoppingItemQuery,
			regularShoppingItem.ExpectedPurchaseDate,
			regularShoppingItem.Purchase,
			regularShoppingItem.Shop,
			regularShoppingItem.Amount,
			regularShoppingItem.BigCategoryID,
			regularShoppingItem.MediumCategoryID,
			regularShoppingItem.CustomCategoryID,
			regularShoppingItemLastInsertId,
			userID,
			regularShoppingItem.TransactionAutoAdd,
		); err != nil {
			return nil, err
		}

		if !today.Before(regularShoppingItem.ExpectedPurchaseDate.Time) {
			if _, err := tx.Exec(
				shoppingItemQuery,
				nextExpectedPurchaseDate,
				regularShoppingItem.Purchase,
				regularShoppingItem.Shop,
				regularShoppingItem.Amount,
				regularShoppingItem.BigCategoryID,
				regularShoppingItem.MediumCategoryID,
				regularShoppingItem.CustomCategoryID,
				regularShoppingItemLastInsertId,
				userID,
				regularShoppingItem.TransactionAutoAdd,
			); err != nil {
				return nil, err
			}
		}

		return result, nil
	}

	result, err := transactions(tx)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}

		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
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
            (expected_purchase_date, purchase, shop, amount, big_category_id, medium_category_id, custom_category_id, user_id, transaction_auto_add)
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
