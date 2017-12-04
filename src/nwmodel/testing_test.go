package nwmodel

import (
	"log"
	"testing"
)

// These tests will fail if TestBox is not running with environment variables
// pointing to it.
func TestGetTest(t *testing.T) {
	challenge := getRandomTest()

	if challenge.ID != "1" {
		log.Fatalf("id did not equal '1', instead: %s", challenge.ID)
	}

	if challenge.Description != "echo stdin" {
		log.Fatalf("equal did not equal 'echo', instead: %s", challenge.Description)
	}
}

func TestSubmitTest(t *testing.T) {
	passMap := submitTest("1", "Python", "import sys\nprint(sys.stdin.read())")

	if passMap.PassFail[""] != true || passMap.PassFail["test"] != true {
		log.Fatal("passMap did not match expected value:")
		log.Println(passMap)
	}
}
