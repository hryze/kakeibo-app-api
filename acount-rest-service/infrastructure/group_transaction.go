package infrastructure

import (
	"time"

	"github.com/paypay3/kakeibo-app-api/acount-rest-service/domain/model"
)

type GroupTransactionsRepository struct {
	*MySQLHandler
}

func (r *TransactionsRepository) GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error) {
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
            group_transactions.transaction_date DESC`

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
