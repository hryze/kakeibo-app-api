package externalapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/xerrors"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"github.com/paypay3/kakeibo-app-api/user-rest-service/infrastructure/externalapi/client"
)

type accountApi struct {
	*client.AccountApiHandler
}

func NewAccountApi(accountApiHandler *client.AccountApiHandler) *accountApi {
	return &accountApi{
		accountApiHandler,
	}
}

func (a *accountApi) PostInitStandardBudgets(userID string) error {
	requestURL := fmt.Sprintf(
		"http://%s:%d/standard-budgets",
		config.Env.AccountApi.Host,
		config.Env.AccountApi.Port,
	)

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

	response, err := a.Client.Do(request.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()

	if response.StatusCode == http.StatusCreated {
		return nil
	}

	return xerrors.New("couldn't create a standard budget")
}
