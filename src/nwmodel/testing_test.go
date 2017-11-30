package nwmodel

import (
	"log"
	"testing"
)

// These tests will fail if TestBox is not running with environment variables
// pointing to it.
func TestGetTest(t *testing.T) {
	id, description := getRandomTest()

	if id != "1" {
		log.Fatalf("id did not equal '1', instead: %s", id)
	}

	if description != "echo" {
		log.Fatalf("equal did not equal 'echo', instead: %s", description)
	}
}

func TestSubmitTest(t *testing.T) {
	passMap := submitTest("1", "Python", "import sys\nprint(sys.stdin.read())")

	if passMap[""] != true || passMap["test"] != true {
		log.Fatal("passMap did not match expected value:")
		log.Println(passMap)
	}
}
