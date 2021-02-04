package interfaces

type AccountApi interface {
	PostInitStandardBudgets(userID string) error
}
