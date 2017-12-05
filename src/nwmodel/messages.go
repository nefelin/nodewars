package nwmodel

import "fmt"

// Message is our basic message struct
type Message struct {
	Type   string `json:"type"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
	Code   string `json:"code"`
} // TODO fix code submission to append to other data, this is unnecessary

const (
	errorStr     = "error:"
	successStr   = "succese:"
	beginStr     = "begin:"
	pseudoStr    = "pseudoServer"
	noConnectStr = "No connection"
)

var msgNoTeam = Message{
	Type:   errorStr,
	Sender: pseudoStr,
	Data:   "No team assignment",
}

var msgNoConnection = Message{
	Type:   errorStr,
	Sender: pseudoStr,
	Data:   "No connection",
}

func psError(e error) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   fmt.Sprintln(e),
	}
}

func psUnknown(cmd string) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "Unknown command '" + cmd + "'",
	}
}

func psSuccess(msg string) Message {
	return Message{
		Type:   successStr,
		Sender: pseudoStr,
		Data:   msg,
	}
}

func psBegin(msg string) Message {
	return Message{
		Type:   beginStr,
		Sender: pseudoStr,
		Data:   msg,
	}
}

// var nwmessages = interface{
// 	NoTeam: Message{"error:", "pseudoServer", "You have no team"},
// }
