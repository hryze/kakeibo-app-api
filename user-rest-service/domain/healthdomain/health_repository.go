package healthdomain

type Repository interface {
	PingDataStore() error
}
