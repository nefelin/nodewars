package nwmodel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"nwmessage"
	"os"
)

// mirrored struct definitions from TestBox
// Challenge does this need to be exported? TODO
type Challenge struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	SampleIO    string            `json:"sampleIO"`
	IO          map[string]string `json:"io"`
}

func (c Challenge) String() string {
	return fmt.Sprintf("( <Challenge> {ID: %s, Desc: %s, SampleIO: %s, IO: %s} )", c.ID, c.Description, c.SampleIO, c.IO)
}

// type CodeSubmission struct {
// 	ID       string   `json:"id"`
// 	Language string   `json:"language"`
// 	Code     string   `json:"code"`
// 	Stdins   []string `json:"stdins,omitempty"`
// }

type CodeSubmission struct {
	Language    string   `json:"language"`
	Code        string   `json:"code"`
	Stdins      []string `json:"stdins"`
	ChallengeId string   `json:"challengeId,omitempty`
}

func (s CodeSubmission) String() string {
	return fmt.Sprintf("( <CodeSubmission> {ChallengeId: %s, Language: %s, Code: Hidden, Stdin: %v} )", s.ChallengeId, s.Language, s.Stdins)
}

// type ExecutionResult struct {
// 	PassFail map[string]string `json:"passFail"`
// 	Message  nwmessage.Message `json:"message"`
// }

type gradeMap map[string]string

func (g gradeMap) String() string {

	var results string
	for k, v := range g {
		results += fmt.Sprintf("%s: %s\n", k, v)
	}
	log.Printf("gradeMap stringer results: %s", results)
	return results
}

type ExecutionResult struct {
	Stdouts []string          `json:"stdouts"`
	Graded  gradeMap          `json:"graded,omitempty"`
	Message nwmessage.Message `json:"message"`
}

func (c ExecutionResult) passed() int {
	var passed int
	for _, res := range c.Graded {
		if res == "Pass" {
			passed++
		}
	}
	return passed
}

func (c ExecutionResult) String() string {
	return fmt.Sprintf("( <ExecutionResult> {Stdouts: %s, Graded: %s, Message: %s} )", c.Stdouts, c.Graded, c.Message)
}

type Language struct {
	// Name          string `json:"name"`
	Boilerplate   string `json:"boilerplate"`
	CommentPrefix string `json:"commentPrefix"`
}

type LanguagesResponse map[string]Language

func getRandomChallenge() Challenge {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	r, err := http.Get(address + ":" + port + "/get_challenge/")

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

	// log.Printf("getRandomChallenge> challenge: %s", chal)
	return chal
}

// returns a map of inputs to test pass/fail
func submitTest(id, language, code string) ExecutionResult {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	submission := CodeSubmission{ChallengeId: id, Language: language, Code: code}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	fmt.Printf("Submitting SubReq: %s", submission)
	r, err := http.Post(address+":"+port+"/submit/", "application/json", buf)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var result ExecutionResult
	err = decoder.Decode(&result)

	log.Printf("submitTest result: %s", result)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return result
}

func getOutput(language, code, stdin string) ExecutionResult {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	submission := CodeSubmission{Language: language, Code: code, Stdins: []string{stdin}}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	r, err := http.Post(address+":"+port+"/stdout/", "application/json", buf)
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(r.Body)
	var response ExecutionResult
	err = decoder.Decode(&response)
	// log.Printf("getOutput response: %v", response)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	return response
}

func getLanguages() map[string]Language {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	langPoint := address + ":" + port + "/languages/"
	// fmt.Printf("testbox at: %s\n", langPoint)

	r, err := http.Get(langPoint)
	if err != nil {
		panic(err)
	}

	// buf := new(bytes.Buffer)
	// buf.ReadFrom(r.Body)
	// // s := buf.String()
	// fmt.Printf("body: %s\n", buf.String())
	decoder := json.NewDecoder(r.Body)
	var langRes LanguagesResponse
	err = decoder.Decode(&langRes)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// map language names into the objects
	// for k := range langRes.Languages {
	// 	langRes.Languages[k].Name = k
	// }

	return langRes
}
