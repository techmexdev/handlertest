package handlertest

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCase is used to test http handlers
type TestCase struct {
	Name       string
	Request    *http.Request
	StatusCode int
	BodyAssert func(b []byte) error
}

// Test sends a fake request to the handler, and compares
// tests the response code, and body
func Test(t *testing.T, tc TestCase, h http.Handler) {
	var errMsg string
	// Store Request body's value
	var buf bytes.Buffer
	tee := io.TeeReader(tc.Request.Body, &buf)
	tc.Request.Body = ioutil.NopCloser(&buf)
	reqBody, err := ioutil.ReadAll(tee)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, tc.Request)
	res := rec.Result()

	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if tc.StatusCode != 0 && res.StatusCode != tc.StatusCode {
		errMsg += fmt.Sprintf("Wrong status code: have %s, want %s. ", http.StatusText(res.StatusCode), http.StatusText(tc.StatusCode))
	}

	if tc.BodyAssert != nil {
		if err := tc.BodyAssert(resBody); err != nil {
			errMsg += fmt.Sprintf("Body assertion failed: %s. ", err.Error())
		}
	}

	if errMsg != "" {
		t.Errorf("%s - %s\n  Request: %s %s %#v %s\n  Response: %v %s - %s\n\n",
			tc.Name, errMsg, tc.Request.Method, tc.Request.URL.Path, tc.Request.Header, reqBody, res.StatusCode, http.StatusText(res.StatusCode), resBody)
	} else {
		t.Logf("%s\n  Request: %s %s %#v %s\n  Response: %v %s - %s\n\n",
			tc.Name, tc.Request.Method, tc.Request.URL.Path, tc.Request.Header, reqBody, res.StatusCode, http.StatusText(res.StatusCode), resBody)
	}
}

// Assert errors returns error when condition is false
func Assert(condition bool, msg string) error {
	if condition {
		return nil
	}
	return errors.New(msg)
}
