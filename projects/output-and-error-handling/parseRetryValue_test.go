package main

import (
	"strings"
	"testing"
	"time"
)

func TestParseRetryValue(t *testing.T) {

	t.Run("parses valid numbers", func(t *testing.T) {

		gotValue, gotErr := ParseRetryValue("5")
		expectedValue := 5

		if gotValue != expectedValue {
			t.Errorf("expected %v but got %v", expectedValue, gotValue)
		}

		if gotErr != nil {
			t.Errorf("expected %v as error but got %v", nil, gotErr)
		}
	})

	t.Run("parses invalid numbers", func(t *testing.T) {

		gotValue, gotErr := ParseRetryValue("a while")
		expectedValue := 0
		expectedErr := "invalid syntax"

		if gotValue != expectedValue {
			t.Errorf("expected %v but got %v", expectedValue, gotValue)
		}

		if !strings.Contains(gotErr.Error(), expectedErr) {
			t.Errorf("expected %v to be contained in %v", expectedErr, gotErr)
		}
	})

	t.Run("parses http time format", func(t *testing.T) {
		duration := 5
		httpTime := time.Now().Add(time.Duration(duration) * time.Second).Format(time.RFC1123)
		expected := duration - 1

		gotValue, gotErr := ParseRetryValue(httpTime)

		

		if gotValue != expected {
			t.Errorf("got %v expected %v", expected, gotValue)
		}

		if gotErr != nil {
			t.Errorf("expected %v as error but got %v", nil, gotErr)
		}

	})
}
