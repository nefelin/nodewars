package challenges

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"nwmessage"
	"os"
)

type ChallengeID = int64

// Challenge holds info for individual programming challenges
type Challenge struct {
	ID        ChallengeID `json:"id"`
	Name      string      `json:"name"`
	ShortDesc string      `json:"shortDesc"`
	LongDesc  string      `json:"longDesc"`
	Tags      tagList     `json:"tags"`
	Cases     caseList    `json:"cases"`
	SampleIO  caseList    `json:"sampleIO"`
}

type safeChallenge struct {
	Name      string   `json:"name"`
	ShortDesc string   `json:"shortDesc"`
	LongDesc  string   `json:"longDesc"`
	Tags      tagList  `json:"tags"`
	SampleIO  caseList `json:"sampleIO"`
}

func (c *Challenge) Redacted() safeChallenge {
	return safeChallenge{
		Name:      c.Name,
		ShortDesc: c.ShortDesc,
		LongDesc:  c.LongDesc,
		Tags:      c.Tags,
		SampleIO:  c.SampleIO,
	}
}

// TestCase describes an individual test case, of which a challenge may have one->many
type TestCase struct {
	Input  string `json:"input"`
	Expect string `json:"expect"`
	Desc   string `json:"desc,omitempty"`
}

type tagList []string
type caseList []TestCase

// CodeSubmission is the format the testbox like to receive in order to execute code, compare against challenge expectations, and return results
type CodeSubmission struct {
	Language    string   `json:"language"`
	Code        string   `json:"code"`
	Stdins      []string `json:"stdins"`
	ChallengeID int64    `json:"challengeId,omitempty`
}

func (s CodeSubmission) String() string {
	return fmt.Sprintf("( <CodeSubmission> {ChallengeID: %s, Language: %s, Code: Hidden, Stdin: %v} )", s.ChallengeID, s.Language, s.Stdins)
}

type tbAPIResponse struct {
	ErrorMessage string `json:"error,omitempty"`
	ID           int64  `json:"id,omitempty"`
	Result       string `json:"result,omitempty"`
}

// GradedResult hold the graded response from testbox
type GradedResult struct {
	Stdouts []string          `json:"stdouts"`
	Grades  []string          `json:"grades,omitempty"`
	Hints   []string          `json:"hints"`
	Message nwmessage.Message `json:"message"`
}

type Solution struct {
	Author   string
	Code     string
	Language string
	Strength int
	IsDummy  bool
}

// SimpleResult holds an unopinionated output from execution
// type SimpleResult struct {
// 	Stdout  string            `json:"stdout"`
// 	Message nwmessage.Message `json:"message"`
// }

func (r GradedResult) gradeMsg() string {
	var res string
	for i, g := range r.Grades {
		res += fmt.Sprintf("Test #%d: ", i)
		if g == "Pass" {
			res += fmt.Sprintf("PASS\n")
		} else {
			res += fmt.Sprintf("FAIL\n")
		}
	}
	return res
}

func (r GradedResult) Passed() int {
	var passed int
	for _, grade := range r.Grades {
		if grade == "Pass" {
			passed++
		}
	}
	return passed
}

func (r GradedResult) String() string {
	return fmt.Sprintf("( <GradedResult> {Stdouts: %s, Graded: %s, Hints: %s, Message: %s} )", r.Stdouts, r.Grades, r.Hints, r.Message)
}

// Language describes the language details nodewars server needs to hold
type Language struct {
	// Name          string `json:"name"`
	Boilerplate   string `json:"boilerplate"`
	CommentPrefix string `json:"commentPrefix"`
}

// LanguagesResponse describes the data format testbox will describe supported languages in.
type LanguagesResponse map[string]Language

func GetRandomChallenge() Challenge {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	r, err := http.Get(address + ":" + port + "/challenges/rand/")

	if err != nil {
		panic(err)
	}

	var c Challenge
	decodeAPIResponse(r, &c)

	return c
}

// returns a map of inputs to test pass/fail
func SubmitTest(id int64, language, code string) (GradedResult, error) {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	submission := CodeSubmission{ChallengeID: id, Language: language, Code: code}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	// fmt.Printf("Submitting SubReq: %s", submission)
	r, err := http.Post(address+":"+port+"/submit/", "application/json", buf)
	if err != nil {
		fmt.Printf("Error submitting: %s\n", err.Error())
		return GradedResult{}, nwmessage.ErrorExecUnavail()
	}

	var e GradedResult
	err = decodeAPIResponse(r, &e)

	if err != nil {
		fmt.Printf("Error decoding tbAPIResponse: %s\n", err.Error())
		return GradedResult{}, nwmessage.ErrorExecUnavail()
	}

	return e, nil
}

func GetOutput(language, code, stdin string) (GradedResult, error) {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	submission := CodeSubmission{Language: language, Code: code, Stdins: []string{stdin}}
	jsonBytes, _ := json.MarshalIndent(submission, "", "    ")
	buf := bytes.NewBuffer(jsonBytes)

	r, err := http.Post(address+":"+port+"/stdout/", "application/json", buf)

	if err != nil {
		fmt.Printf("Error submitting: %s\n", err.Error())
		return GradedResult{}, nwmessage.ErrorExecUnavail()
	}

	var e GradedResult
	err = decodeAPIResponse(r, &e)

	if err != nil {
		fmt.Printf("Error decoding tbAPIResponse: %s\n", err.Error())
		return GradedResult{}, nwmessage.ErrorExecUnavail()
	}

	return e, nil
}

func GetLanguages() map[string]Language {
	address := os.Getenv("TESTBOX_ADDRESS")
	port := os.Getenv("TESTBOX_PORT")

	langPoint := address + ":" + port + "/languages/"
	// fmt.Printf("testbox at: %s\n", langPoint)

	r, err := http.Get(langPoint)
	if err != nil {
		panic(err)
	}

	var l LanguagesResponse
	decodeAPIResponse(r, &l)

	return l
}

func decodeAPIResponse(r *http.Response, i interface{}) error {
	decoder := json.NewDecoder(r.Body)
	var resp tbAPIResponse
	err := decoder.Decode(&resp)
	if err != nil {
		// panic(err)
		return err
	}
	// fmt.Printf("Got apiresponse: %v\n\n\n", resp.Result)
	defer r.Body.Close()

	err = json.Unmarshal([]byte(resp.Result), &i)
	if err != nil {
		// panic(err)
		return err
	}

	return nil
}
