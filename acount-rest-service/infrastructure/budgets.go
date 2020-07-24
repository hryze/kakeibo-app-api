package infrastructure

type BudgetsRepository struct {
	*MySQLHandler
}

func (r *BudgetsRepository) PostInitStandardBudgets(userID string) error {
	query := `
        INSERT INTO standard_budgets
            (user_id, big_category_id)
        VALUES
            (?,2),
            (?,3),
            (?,4),
            (?,5),
            (?,6),
            (?,7),
            (?,8),
            (?,9),
            (?,10),
            (?,11),
            (?,12),
            (?,13),
            (?,14),
            (?,15),
            (?,16),
            (?,17)`
	_, err := r.MySQLHandler.conn.Exec(query, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID, userID)
	return err
}
