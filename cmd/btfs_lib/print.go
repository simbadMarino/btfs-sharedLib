package main

import (
	"fmt"
	"io"
	logger "log"
	"net/http"
	"strings"
)

var url = "https://kvdb.io/6hYHSuSjhuX8G1EXsHKrJw/btfs"

func callAPI(body string) {
	resp, getErr := http.Get(url)
	if getErr != nil || resp.StatusCode != 200 {
		logger.Fatal(getErr)
	}

	oldBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)

	}

	resp, postErr := http.Post(url, "text", strings.NewReader(string(oldBody)+body+"\n"))
	if postErr != nil || resp.StatusCode != 200 {
		logger.Fatal(postErr)
	}
}

func Println(a ...interface{}) {
	callAPI(fmt.Sprint(a...))
}

func Sprint(a ...interface{}) string {
	msg := fmt.Sprint(a...)
	callAPI(msg)
	return msg
}

func Errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	callAPI(err.Error())
	return err
}

func Sprintf(format string, a ...interface{}) string {
	msg := fmt.Sprintf(format, a...)
	callAPI(msg)
	return msg
}

func Printf(format string, a ...interface{}) {
	callAPI(fmt.Sprintf(format, a...))
}
