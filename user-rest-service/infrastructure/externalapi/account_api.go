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

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
	"golang.org/x/xerrors"
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

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	response, err := a.fetch(ctx, request)
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

func (a *accountApi) fetch(ctx context.Context, request *http.Request) (*http.Response, error) {
	request = request.WithContext(ctx)

	responseCh := make(chan *http.Response)
	errorCh := make(chan error)

	go func() {
		response, err := a.HttpClient.Do(request)
		if err != nil {
			errorCh <- err
			return
		}

		responseCh <- response
	}()

	select {
	case response := <-responseCh:
		return response, nil

	case err := <-errorCh:
		return nil, err

	case <-ctx.Done():
		return nil, xerrors.New("HTTP request cancelled")
	}
}
