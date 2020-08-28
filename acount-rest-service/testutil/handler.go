package testutil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func GetJsonFromTestData(t *testing.T, path string) string {
	t.Helper()

	byteData, err := ioutil.ReadFile(path)
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

func AssertResponseBody(t *testing.T, res *http.Response, path string) {
	t.Helper()

	want := GetJsonFromTestData(t, path)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadAll() '%#v'", err)
	}

	var got bytes.Buffer
	err = json.Indent(&got, body, "", "  ")
	if err != nil {
		t.Fatalf("unexpected error by json.Indent() '%#v'", err)
	}

	if diff := cmp.Diff(want, got.String()); len(diff) != 0 {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}
