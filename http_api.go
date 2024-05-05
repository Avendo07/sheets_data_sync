package main

import (
	"bytes"
	"net/http"
)

func postCall(url string, payload []byte, headers map[string]string) (int, error) {
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	for key, value := range headers {
		r.Header.Set(key, value)
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	return res.StatusCode, nil
}
