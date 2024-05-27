package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestParseResponse(t *testing.T) {

	tests := map[string]struct {
		input        http.Response
		expectedBody string
		expectedErr  RequestError
	}{
		"429":     {input: http.Response{StatusCode: 429, Header: http.Header{"Retry-After": []string{"5"}}, Body: io.NopCloser(bytes.NewBufferString(""))}, expectedBody: "", expectedErr: RequestError{SleepTime: time.Duration(5), ShouldRetry: true, Cause: fmt.Errorf("server is busy, retrying in few seconds")}},
		"200":     {input: http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewBufferString("Today is sunny!"))}, expectedBody: "Today is sunny!", expectedErr: RequestError{SleepTime: time.Duration(0), ShouldRetry: false, Cause: nil}},
		"a-while": {input: http.Response{StatusCode: 429, Header: http.Header{"Retry-After": []string{"a-while"}}, Body: io.NopCloser(bytes.NewBufferString(""))}, expectedBody: "", expectedErr: RequestError{SleepTime: time.Duration(0), ShouldRetry: false, Cause: fmt.Errorf("server is busy, retrying in few seconds")}},
		"500":     {input: http.Response{StatusCode: 429, Body: io.NopCloser(bytes.NewBufferString(""))}, expectedBody: "", expectedErr: RequestError{SleepTime: time.Duration(0), ShouldRetry: false, Cause: fmt.Errorf("server is busy, retrying in few seconds")}},
	}

	for _, tc := range tests {
		gotBody, gotError := parseResponse(tc.input)
		assertResponse(t, gotBody, tc.expectedBody)
		assertResponse(t, gotError.SleepTime, tc.expectedErr.SleepTime)
		assertResponse(t, gotError.ShouldRetry, tc.expectedErr.ShouldRetry)
	}
}

func assertResponse(t testing.TB, got, want interface{}) {
	t.Helper()
	if got != want {
		t.Errorf("got %v but expected %v\n", got, want)
	}

}
