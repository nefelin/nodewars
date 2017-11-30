package nwmodel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

// mirrored struct definitions from TestBox
type TestResponse struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

type SubmissionRequest struct {
	Id       string `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type SubmissionResponse struct {
	PassedTests map[string]bool `json:"passedTests"`
}

func getRandomTest() (id, description string) {
	address := os.Getenv("TEST_BOX_ADDRESS")
	port := os.Getenv("TEST_BOX_PORT")

	r, err := http.Get(address + port)
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(r.Body)
	var test TestResponse
	err = decoder.Decode(&test)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return test.Id, test.Description
}

// returns a map of inputs to test pass/fail
func submitTest(id, language, code string) map[string]bool {
	address := os.Getenv("TEST_BOX_ADDRESS")
	port := os.Getenv("TEST_BOX_PORT")

	submission := SubmissionRequest{id, language, code}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	r, err := http.Post(address+port+"/submit/", "application/json", buf)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var passed SubmissionResponse
	err = decoder.Decode(&passed)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return passed.PassedTests
}
