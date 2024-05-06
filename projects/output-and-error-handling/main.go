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
	SleepTime int
	Retry     bool
	ErrMsg       error
}



func parseRequest(resp http.Response) (reqBody string, reqError RequestError) {

	switch resp.StatusCode {
	case 429:
		retryTime := resp.Header.Get("Retry-After")
		formattedRetryTime, err := ParseRetryValue(retryTime)

		if err != nil {
			// cannot determine how long to sleep for, give up
			reqError.ErrMsg = errors.New("cannot get weather at the moment try again later")
			return
		}
		if formattedRetryTime > 5 {
			// sleep time is long give up
			reqError.ErrMsg = errors.New("server is busy at the moment try again later")

		} else {
			reqError.ErrMsg = errors.New("server is busy, retrying in few seconds")
			reqError.SleepTime = formattedRetryTime
			reqError.Retry = true
		}
	default:
		body, _ := io.ReadAll(resp.Body)
		reqBody = string(body)
	}
	return

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

		if reqErr.ErrMsg != nil {
			fmt.Fprintln(os.Stderr, reqErr.ErrMsg)

			if reqErr.Retry {
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

func ParseRetryValue(v string) (int, error) {

	value, err := strconv.Atoi(v)
	if err != nil {
		// check if retry time is in httpTime format
		httpTime, httpTimeErr := time.Parse(time.RFC1123, v)

		if httpTimeErr == nil {
			value = int(time.Until(httpTime).Seconds())
			err = httpTimeErr
		}
	}
	return value, err

}

// sleep for a given time
func sleep(t int) {
	if t > 1 {
		fmt.Fprint(os.Stderr, "We are currently retrying your request. Things might take a bit longer than usual\n")
	}
	time.Sleep(time.Duration(t) * time.Second)
}
