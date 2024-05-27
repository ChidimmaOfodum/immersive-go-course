package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type RequestError struct {
	SleepTime time.Duration
	ShouldRetry     bool
	Cause       error
}
func main() {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	for i := 0; i < 3; i++ {
		resp, err := c.Get("http://localhost:8080")

		if err != nil {
			fmt.Fprint(os.Stderr, "cannot get weather at the moment. Please try again later\n")
			os.Exit(2)
		}
		defer resp.Body.Close()
		
		body, reqErr := parseRequest(*resp)

		if reqErr.Cause != nil {
			fmt.Fprintln(os.Stderr, reqErr.Cause)

			if reqErr.ShouldRetry {
				sleep(reqErr.SleepTime)
				continue
			} else {
				os.Exit(28)
			}
		} else {
			fmt.Fprintln(os.Stdout, body)
			break
		}
	}
}

func parseRequest(resp http.Response) (reqBody string, reqError RequestError) {

	switch resp.StatusCode {
	case 429:
		retryTime := resp.Header.Get("Retry-After")
		formattedRetryTime, err := ParseRetryValue(retryTime)
		fmt.Print(err)
		fmt.Print(formattedRetryTime)

		if err != nil {
			// cannot determine how long to sleep for, give up
			reqError.Cause = errors.New("cannot get weather at the moment try again later")
			return
		}
		if formattedRetryTime > 5 {
			// sleep time is long give up
			reqError.Cause = errors.New("server is busy at the moment try again later")

		} else {
			reqError.Cause = errors.New("server is busy, retrying in few seconds")
			reqError.SleepTime = time.Duration(formattedRetryTime)
			reqError.ShouldRetry = true
		}
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			reqError.Cause = err
		}
		reqBody = string(body)
	}
	return

}

func ParseRetryValue(v string) (int, error) {
	if value, err := strconv.Atoi(v); err == nil {
		return value, err
	}
	if value, err := time.Parse(time.RFC1123, v); err == nil {
		return int(time.Until(value).Seconds()), err
	}
	return 0, fmt.Errorf("couldn't parse header as int or timestamp")
}

// sleep for a given time
func sleep(t time.Duration) {
	if t > 1 {
		fmt.Fprint(os.Stderr, "We are currently retrying your request. Things might take a bit longer than usual\n")
	}
	time.Sleep(time.Duration(t) * time.Second)
}
