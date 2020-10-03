package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type TransactionsRepository struct {
	*MySQLHandler
}

func NewTransactionsRepository(mysqlHandler *MySQLHandler) *TransactionsRepository {
	return &TransactionsRepository{mysqlHandler}
}

func (r *TransactionsRepository) GetMonthlyTransactionsList(userID string, firstDay time.Time, lastDay time.Time) ([]model.TransactionSender, error) {
	query := `
        SELECT
            transactions.id id,
            transactions.transaction_type transaction_type,
            transactions.updated_date updated_date,
            transactions.transaction_date transaction_date,
            transactions.shop shop,
            transactions.memo memo,
            transactions.amount amount,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            custom_categories.category_name custom_category_name
        FROM
            transactions
        LEFT JOIN
            big_categories
        ON
            transactions.big_category_id = big_categories.id
        LEFT JOIN
            medium_categories
        ON
            transactions.medium_category_id = medium_categories.id
        LEFT JOIN
            custom_categories
        ON
            transactions.custom_category_id = custom_categories.id
        WHERE
            transactions.user_id = ?
        AND
            transactions.transaction_date >= ?
        AND
            transactions.transaction_date <= ?
        ORDER BY
            transactions.transaction_date DESC, transactions.updated_date DESC`
	rows, err := r.MySQLHandler.conn.Queryx(query, userID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionsList []model.TransactionSender
	for rows.Next() {
		var transactionSender model.TransactionSender
		if err := rows.StructScan(&transactionSender); err != nil {
			return nil, err
		}
		transactionsList = append(transactionsList, transactionSender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactionsList, nil
}

func (r *TransactionsRepository) GetTransaction(transactionSender *model.TransactionSender, transactionID int) (*model.TransactionSender, error) {
	query := `
        SELECT
            transactions.id id,
            transactions.transaction_type transaction_type,
            transactions.updated_date updated_date,
            transactions.transaction_date transaction_date,
            transactions.shop shop,
            transactions.memo memo,
            transactions.amount amount,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            custom_categories.category_name custom_category_name
        FROM
            transactions
        LEFT JOIN
            big_categories
        ON
            transactions.big_category_id = big_categories.id
        LEFT JOIN
            medium_categories
        ON
            transactions.medium_category_id = medium_categories.id
        LEFT JOIN
            custom_categories
        ON
            transactions.custom_category_id = custom_categories.id
        WHERE
            transactions.id = ?`
	if err := r.MySQLHandler.conn.QueryRowx(query, transactionID).StructScan(transactionSender); err != nil {
		return nil, err
	}
	return transactionSender, nil
}

func (r *TransactionsRepository) PostTransaction(transaction *model.TransactionReceiver, userID string) (sql.Result, error) {
	query := `
        INSERT INTO transactions
            (transaction_type, transaction_date, shop, memo, amount, user_id, big_category_id, medium_category_id, custom_category_id)
        VALUES
            (?,?,?,?,?,?,?,?,?)`
	result, err := r.MySQLHandler.conn.Exec(query, transaction.TransactionType, transaction.TransactionDate, transaction.Shop, transaction.Memo, transaction.Amount, userID, transaction.BigCategoryID, transaction.MediumCategoryID, transaction.CustomCategoryID)
	return result, err
}

func (r *TransactionsRepository) PutTransaction(transaction *model.TransactionReceiver, transactionID int) error {
	query := `
        UPDATE
            transactions
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
	_, err := r.MySQLHandler.conn.Exec(query, transaction.TransactionType, transaction.TransactionDate, transaction.Shop, transaction.Memo, transaction.Amount, transaction.BigCategoryID, transaction.MediumCategoryID, transaction.CustomCategoryID, transactionID)
	return err
}

func (r *TransactionsRepository) DeleteTransaction(transactionID int) error {
	query := `
        DELETE
        FROM 
            transactions
        WHERE 
            id = ?`
	_, err := r.MySQLHandler.conn.Exec(query, transactionID)
	return err
}

func (r *TransactionsRepository) SearchTransactionsList(query string) ([]model.TransactionSender, error) {
	rows, err := r.MySQLHandler.conn.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactionsList []model.TransactionSender
	for rows.Next() {
		var transactionSender model.TransactionSender
		if err := rows.StructScan(&transactionSender); err != nil {
			return nil, err
		}
		transactionsList = append(transactionsList, transactionSender)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactionsList, nil
}
