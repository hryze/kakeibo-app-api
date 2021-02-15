package testutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/mux"
)

func SetUpMockServer() func() {
	if err := os.Setenv("USER_HOST", "localhost"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	verifyGroupAffiliationHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	getGroupUserIDListHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDList := []string{"userID1", "userID4", "userID5", "userID3", "userID2"}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&userIDList); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	router := mux.NewRouter()
	router.HandleFunc("/groups/{group_id:[0-9]+}/users/{user_id:[\\S]{1,10}}/verify", verifyGroupAffiliationHandler).Methods("GET")
	router.HandleFunc("/groups/{group_id:[0-9]+}/users", getGroupUserIDListHandler).Methods("GET")

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ts := httptest.Server{
		Listener: listener,
		Config:   &http.Server{Handler: router},
	}

	ts.Start()

	return func() {
		ts.Close()
	}
}

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
