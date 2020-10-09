package handler

import (
	"testing"

	"github.com/paypay3/kakeibo-app-api/account-rest-service/testutil"
)

func TestMain(m *testing.M) {
	tearDown := testutil.SetUpMockServer()
	defer tearDown()

	m.Run()
}
