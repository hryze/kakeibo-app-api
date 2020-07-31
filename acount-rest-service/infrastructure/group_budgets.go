package infrastructure

type GroupBudgetsRepository struct {
	*MySQLHandler
}

func (r *GroupBudgetsRepository) PostInitGroupStandardBudgets(groupID int) error {
	query := `
        INSERT INTO group_standard_budgets
            (group_id, big_category_id)
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

	_, err := r.MySQLHandler.conn.Exec(query, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID, groupID)
	return err
}
