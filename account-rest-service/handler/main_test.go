package handler

import (
	"testing"

	"github.com/hryze/kakeibo-app-api/account-rest-service/testutil"
)

func TestMain(m *testing.M) { //nolint:staticcheck
	tearDown := testutil.SetUpMockServer()
	m.Run()
	tearDown()
}
