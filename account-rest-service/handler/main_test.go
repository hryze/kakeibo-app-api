package handler

import (
	"os"
	"testing"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

func TestMain(m *testing.M) {
	tearDown := testutil.SetUpMockServer()
	status := m.Run()
	tearDown()

	os.Exit(status)
}
