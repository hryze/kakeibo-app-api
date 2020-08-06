package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type GroupTransactionsRepository struct {
	*MySQLHandler
}

func (r *GroupTransactionsRepository) GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            custom_categories.category_name custom_category_name
        FROM
            group_transactions
        LEFT JOIN
            big_categories
        ON
            group_transactions.big_category_id = big_categories.id
        LEFT JOIN
            medium_categories
        ON
            group_transactions.medium_category_id = medium_categories.id
        LEFT JOIN
            custom_categories
        ON
            group_transactions.custom_category_id = custom_categories.id
        WHERE
            group_transactions.group_id = ?
        AND
            group_transactions.transaction_date >= ?
        AND
            group_transactions.transaction_date <= ?
        ORDER BY
            group_transactions.transaction_date DESC, group_transactions.updated_date DESC`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupTransactionsList []model.GroupTransactionSender
	for rows.Next() {
		var groupTransactionSender model.GroupTransactionSender
		if err := rows.StructScan(&groupTransactionSender); err != nil {
			return nil, err
		}
		groupTransactionsList = append(groupTransactionsList, groupTransactionSender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupTransactionsList, nil
}

func (r *GroupTransactionsRepository) GetGroupTransaction(groupTransactionID int) (*model.GroupTransactionSender, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            custom_categories.category_name custom_category_name
        FROM
            group_transactions
        LEFT JOIN
            big_categories
        ON
            group_transactions.big_category_id = big_categories.id
        LEFT JOIN
            medium_categories
        ON
            group_transactions.medium_category_id = medium_categories.id
        LEFT JOIN
            custom_categories
        ON
            group_transactions.custom_category_id = custom_categories.id
        WHERE
            group_transactions.id = ?`

	var groupTransactionSender model.GroupTransactionSender
	if err := r.MySQLHandler.conn.QueryRowx(query, groupTransactionID).StructScan(&groupTransactionSender); err != nil {
		return nil, err
	}

	return &groupTransactionSender, nil
}

func (r *GroupTransactionsRepository) PostGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupID int, userID string) (sql.Result, error) {
	query := `
        INSERT INTO group_transactions
            (transaction_type, transaction_date, shop, memo, amount, group_id, user_id, big_category_id, medium_category_id, custom_category_id)
        VALUES
            (?,?,?,?,?,?,?,?,?,?)`

	result, err := r.MySQLHandler.conn.Exec(query, groupTransaction.TransactionType, groupTransaction.TransactionDate, groupTransaction.Shop, groupTransaction.Memo, groupTransaction.Amount, groupID, userID, groupTransaction.BigCategoryID, groupTransaction.MediumCategoryID, groupTransaction.CustomCategoryID)

	return result, err
}

func (r *GroupTransactionsRepository) PutGroupTransaction(groupTransaction *model.GroupTransactionReceiver, groupTransactionID int) error {
	query := `
        UPDATE
            group_transactions
        SET 
            transaction_type = ?,
            transaction_date = ?,
            shop = ?,
            memo = ?,
            amount = ?,
            big_category_id = ?,
            medium_category_id = ?,
            custom_category_id = ?
        WHERE
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTransaction.TransactionType, groupTransaction.TransactionDate, groupTransaction.Shop, groupTransaction.Memo, groupTransaction.Amount, groupTransaction.BigCategoryID, groupTransaction.MediumCategoryID, groupTransaction.CustomCategoryID, groupTransactionID)

	return err
}

func (r *GroupTransactionsRepository) DeleteGroupTransaction(groupTransactionID int) error {
	query := `
        DELETE
        FROM 
            group_transactions
        WHERE 
            id = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupTransactionID)

	return err
}