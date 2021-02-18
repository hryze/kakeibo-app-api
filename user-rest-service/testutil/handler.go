package testutil

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/paypay3/kakeibo-app-api/user-rest-service/config"
)

func GetRequestJsonFromTestData(t *testing.T) string {
	t.Helper()

	requestFilePath := filepath.Join("testdata", t.Name(), "request.json")

	byteData, err := ioutil.ReadFile(requestFilePath)
	if err != nil {
		t.Fatalf("unexpected error while opening file '%#v'", err)
	}

	return string(byteData)
}

func AssertResponseHeader(t *testing.T, res *http.Response, wantStatusCode int) {
	t.Helper()

	if diff := cmp.Diff(wantStatusCode, res.StatusCode); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff("application/json; charset=UTF-8", res.Header.Get("Content-Type")); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func AssertResponseBody(t *testing.T, res *http.Response, wantStruct interface{}, gotStruct interface{}) {
	t.Helper()

	goldenFilePath := filepath.Join("testdata", t.Name(), "response.json.golden")

	wantData, err := ioutil.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadFile '%#v'", err)
	}

	gotData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadAll() '%#v'", err)
	}

	if err := json.Unmarshal(wantData, wantStruct); err != nil {
		t.Fatalf("unexpected error by json.Unmarshal() '%#v'", err)
	}

	if err := json.Unmarshal(gotData, gotStruct); err != nil {
		t.Fatalf("unexpected error by json.Unmarshal() '%#v'", err)
	}

	if diff := cmp.Diff(wantStruct, gotStruct); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func AssertSetResponseCookie(t *testing.T, res *http.Response) {
	t.Helper()

	cookie := res.Cookies()[0]

	if diff := cmp.Diff(config.Env.Cookie.Name, cookie.Name); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(36, len(cookie.Value)); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(true, time.Now().UTC().Before(cookie.Expires)); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(config.Env.Cookie.Domain, cookie.Domain); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(config.Env.Cookie.Secure, cookie.Secure); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(true, cookie.HttpOnly); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

}

func AssertDeleteResponseCookie(t *testing.T, res *http.Response) {
	t.Helper()

	cookie := res.Cookies()[0]

	if diff := cmp.Diff(0, len(cookie.Value)); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

	if diff := cmp.Diff(true, !time.Now().UTC().Before(cookie.Expires)); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
