package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type ShoppingListRepository struct {
	*MySQLHandler
}

func NewShoppingListRepository(mysqlHandler *MySQLHandler) *ShoppingListRepository {
	return &ShoppingListRepository{mysqlHandler}
}

func (r *ShoppingListRepository) PostShoppingItem(shoppingItem *model.ShoppingItem, userID string) (sql.Result, error) {
	query := `
        INSERT INTO shopping_list
            (expected_purchase_date, purchase, shop, amount, big_category_id, medium_category_id, custom_category_id, user_id, transaction_auto_add)
        VALUES
            (?,?,?,?,?,?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, shoppingItem.ExpectedPurchaseDate, shoppingItem.Purchase, shoppingItem.Shop, shoppingItem.Amount, shoppingItem.BigCategoryID, shoppingItem.MediumCategoryID, shoppingItem.CustomCategoryID, userID, shoppingItem.TransactionAutoAdd)
	return result, err
}

func (r *ShoppingListRepository) GetShoppingItem(shoppingItemID int) (model.ShoppingItem, error) {
	query := `
        SELECT
            id,
            posted_date,
            updated_date,
            expected_purchase_date,
            complete_date,
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
