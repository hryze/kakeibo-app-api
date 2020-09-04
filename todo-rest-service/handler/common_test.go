package handler

type MockAuthRepository struct{}

func (t MockAuthRepository) GetUserID(sessionID string) (string, error) {
	return "userID1", nil
}
