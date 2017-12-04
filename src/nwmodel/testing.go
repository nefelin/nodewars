package nwmodel

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// mirrored struct definitions from TestBox
// Challenge does this need to be exported? TODO
type Challenge struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type SubmissionRequest struct {
	Id       string `json:"id"`
	Language string `json:"language"`
	Code     string `json:"code"`
}

type ChallengeResponse struct {
	PassFail map[string]bool `json:"passFail"`
}

func (c ChallengeResponse) passed() int {
	var passed int
	for _, res := range c.PassFail {
		if res == true {
			passed++
		}
	}
	return passed
}

func getRandomTest() Challenge {
	// address := os.Getenv("TEST_BOX_ADDRESS")
	// port := os.Getenv("TEST_BOX_PORT")

	address := "http://localhost"
	port := "31337"
	r, err := http.Get(address + ":" + port)

	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(r.Body)
	var chal Challenge
	err = decoder.Decode(&chal)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return chal
}

// returns a map of inputs to test pass/fail
func submitTest(id, language, code string) ChallengeResponse {
	// address := os.Getenv("TEST_BOX_ADDRESS")
	// port := os.Getenv("TEST_BOX_PORT")
	address := "http://localhost"
	port := "31337"

	submission := SubmissionRequest{id, language, code}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	r, err := http.Post(address+":"+port+"/submit/", "application/json", buf)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var response ChallengeResponse
	err = decoder.Decode(&response)
	log.Printf("submitTest response: %v", response)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return response
}
