package infrastructure

import (
	"database/sql"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type TransactionsRepository struct {
	*MySQLHandler
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
