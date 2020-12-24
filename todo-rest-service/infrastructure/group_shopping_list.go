package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/todo-rest-service/domain/model"
)

type GroupShoppingListRepository struct {
	*MySQLHandler
}

func NewGroupShoppingListRepository(mysqlHandler *MySQLHandler) *GroupShoppingListRepository {
	return &GroupShoppingListRepository{mysqlHandler}
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
