package externalapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
)

type accountApi struct {
	*config.AccountApiHandler
}

func NewAccountApi(accountApiHandler *config.AccountApiHandler) *accountApi {
	return &accountApi{accountApiHandler}
}

func (a *accountApi) PostInitStandardBudgets(userID string) error {
	accountHost := os.Getenv("ACCOUNT_HOST")
	requestURL := fmt.Sprintf("http://%s:8081/standard-budgets", accountHost)

	request, err := http.NewRequest(
		"POST",
		requestURL,
		bytes.NewBuffer([]byte(fmt.Sprintf(`{ "user_id" : "%s" }`, userID))),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response, err := a.HttpClient.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusCreated {
		return nil
	}

	return xerrors.New("couldn't create a standard budget")
}
