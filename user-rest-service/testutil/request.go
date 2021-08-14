package testutil

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// request file name info
	requestFileNamePrefix = "request"
	requestFileNameSuffix = ".json"
)

func GetRequestJsonFromTestDataV2(t *testing.T, fileNameOpts ...string) string {
	t.Helper()

	requestFilePath := newRequestFilePath(t, fileNameOpts...)

	byteData, err := ioutil.ReadFile(requestFilePath)
	if err != nil {
		t.Fatalf("unexpected error by ioutil.ReadFile '%v'", err)
	}

	return string(byteData)
}

func newRequestFilePath(t *testing.T, fileNameOpts ...string) string {
	var requestFilePath string
	testFuncName := strings.Split(t.Name(), "/")[0]

	if len(fileNameOpts) == 0 {
		// file path example: ./testdata/{testFuncName}/request.json
		requestFilePath = filepath.Join(
			fixtureDir,
			testFuncName,
			fmt.Sprintf("%s%s", requestFileNamePrefix, requestFileNameSuffix),
		)
	} else {
		optsStr := strings.Join(fileNameOpts, "-")

		// file path example: ./testdata/{testFuncName}/request-{fileNameOpts}.json
		requestFilePath = filepath.Join(
			fixtureDir,
			testFuncName,
			fmt.Sprintf("%s-%s%s", requestFileNamePrefix, optsStr, requestFileNameSuffix),
		)
	}

	return requestFilePath
}
