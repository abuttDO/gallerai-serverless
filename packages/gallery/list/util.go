package main

import (
	"encoding/json"
	"fmt"
)

func makeResponse(statusCode int, body interface{}, err error) *Response {
	bodyString, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("error marshalling body: %s\n", err)
		return &Response{}
	}
	var response Response
	response.StatusCode = statusCode
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	response.Body = string(bodyString)
	return &response
}
