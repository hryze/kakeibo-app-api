package infrastructure

import (
	"database/sql"
	"time"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/domain/model"
)

type GroupTransactionsRepository struct {
	*MySQLHandler
}

func NewGroupTransactionsRepository(mysqlHandler *MySQLHandler) *GroupTransactionsRepository {
	return &GroupTransactionsRepository{mysqlHandler}
}

func (r *GroupTransactionsRepository) GetMonthlyGroupTransactionsList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.GroupTransactionSender, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.posted_date posted_date,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            group_custom_categories.category_name custom_category_name
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
            group_custom_categories
        ON
            group_transactions.custom_category_id = group_custom_categories.id
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

func (r *GroupTransactionsRepository) Get10LatestGroupTransactionsList(groupID int) (*model.GroupTransactionsList, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.posted_date posted_date,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            group_custom_categories.category_name custom_category_name
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
            group_custom_categories
        ON
            group_transactions.custom_category_id = group_custom_categories.id
        WHERE
            group_transactions.group_id = ?
        ORDER BY
            group_transactions.updated_date DESC
        LIMIT
            10`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groupTransactionsList := model.GroupTransactionsList{
		GroupTransactionsList: make([]model.GroupTransactionSender, 0),
	}
	for rows.Next() {
		var groupTransactionSender model.GroupTransactionSender
		if err := rows.StructScan(&groupTransactionSender); err != nil {
			return nil, err
		}

		groupTransactionsList.GroupTransactionsList = append(groupTransactionsList.GroupTransactionsList, groupTransactionSender)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &groupTransactionsList, nil
}

func (r *GroupTransactionsRepository) GetGroupTransaction(groupTransactionID int) (*model.GroupTransactionSender, error) {
	query := `
        SELECT
            group_transactions.id id,
            group_transactions.transaction_type transaction_type,
            group_transactions.posted_date posted_date,
            group_transactions.updated_date updated_date,
            group_transactions.transaction_date transaction_date,
            group_transactions.shop shop,
            group_transactions.memo memo,
            group_transactions.amount amount,
            group_transactions.user_id user_id,
            big_categories.category_name big_category_name,
            medium_categories.category_name medium_category_name,
            group_custom_categories.category_name custom_category_name
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
            group_custom_categories
        ON
            group_transactions.custom_category_id = group_custom_categories.id
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

func (r *GroupTransactionsRepository) SearchGroupTransactionsList(query string) ([]model.GroupTransactionSender, error) {
	rows, err := r.MySQLHandler.conn.Queryx(query)
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

func (r *GroupTransactionsRepository) GetUserPaymentAmountList(groupID int, firstDay time.Time, lastDay time.Time) ([]model.UserPaymentAmount, error) {
	query := `
        SELECT
            user_id,
            SUM(amount) total_payment_amount
        FROM
            group_transactions
        WHERE
            group_id = ?
        AND
            transaction_date >= ?
        AND
            transaction_date < ?
        GROUP BY
            user_id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, firstDay, lastDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userPaymentAmountList []model.UserPaymentAmount
	for rows.Next() {
		var userPaymentAmount model.UserPaymentAmount
		if err := rows.StructScan(&userPaymentAmount); err != nil {
			return nil, err
		}

		userPaymentAmountList = append(userPaymentAmountList, userPaymentAmount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userPaymentAmountList, nil
}

func (r *GroupTransactionsRepository) GetGroupAccountsList(yearMonth time.Time, groupID int) ([]model.GroupAccount, error) {
	query := `
        SELECT
            id,
            years_months,
            payer_user_id,
            recipient_user_id,
            payment_amount,
            payment_confirmation,
            receipt_confirmation,
            group_id
        FROM
            group_accounts
        WHERE
            group_id = ?
        AND
            years_months = ?
        ORDER BY
            id`

	rows, err := r.MySQLHandler.conn.Queryx(query, groupID, yearMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groupAccountsList []model.GroupAccount
	for rows.Next() {
		var groupAccount model.GroupAccount
		if err := rows.StructScan(&groupAccount); err != nil {
			return nil, err
		}

		groupAccountsList = append(groupAccountsList, groupAccount)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupAccountsList, nil
}

func (r *GroupTransactionsRepository) PostGroupAccountsList(groupAccountsList []model.GroupAccount, yearMonth time.Time, groupID int) error {
	query := `
        INSERT INTO group_accounts
            (years_months, payer_user_id, recipient_user_id, payment_amount, group_id)
        VALUES
            (?,?,?,?,?)`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, groupAccount := range groupAccountsList {
			if _, err := r.MySQLHandler.conn.Exec(query, yearMonth, groupAccount.Payer, groupAccount.Recipient, groupAccount.PaymentAmount, groupID); err != nil {
				return err
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

func (r *GroupTransactionsRepository) PutGroupAccountsList(groupAccountsList []model.GroupAccount) error {
	query := `
        UPDATE
            group_accounts
        SET 
            payment_confirmation = ?,
            receipt_confirmation = ?
        WHERE
            id = ?`

	tx, err := r.MySQLHandler.conn.Begin()
	if err != nil {
		return err
	}

	transactions := func(tx *sql.Tx) error {
		for _, groupAccount := range groupAccountsList {
			if _, err := r.MySQLHandler.conn.Exec(query, groupAccount.PaymentConfirmation, groupAccount.ReceiptConfirmation, groupAccount.ID); err != nil {
				return err
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

func (r *GroupTransactionsRepository) DeleteGroupAccountsList(yearMonth time.Time, groupID int) error {
	query := `
        DELETE
        FROM 
            group_accounts
        WHERE 
            group_id = ?
        AND
            years_months = ?`

	_, err := r.MySQLHandler.conn.Exec(query, groupID, yearMonth)

	return err
}
