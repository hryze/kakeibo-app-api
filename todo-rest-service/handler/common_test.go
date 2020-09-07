package handler

type MockAuthRepository struct{}

type MockSqlResult struct{}

func (t MockAuthRepository) GetUserID(sessionID string) (string, error) {
	return "userID1", nil
}

func (r MockSqlResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (r MockSqlResult) RowsAffected() (int64, error) {
	return 1, nil
}
