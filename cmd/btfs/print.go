package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var url = "https://kvdb.io/CZUbTdFXbiRJcZBtAviPzR/hello"

func callAPI(body string) {
	resp, getErr := http.Get(url)
	if getErr != nil || resp.StatusCode != 200 {
		log.Fatal(getErr)
	}

	oldBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)

	}

	resp, postErr := http.Post(url, "text", strings.NewReader(string(oldBody)+body+"\n"))
	if postErr != nil || resp.StatusCode != 200 {
		log.Fatal(postErr)
	}
}

func Println(a ...interface{}) {
	callAPI(fmt.Sprint(a...))
}

func Printf(format string, a ...interface{}) {
	callAPI(fmt.Sprintf(format, a...))
}
