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
	alertStr   = "alert:"
	errorStr   = "error:"
	successStr = "success:"
	beginStr   = "begin:"

	editStateStr = "editorState"
	pseudoStr    = "pseudoServer"
	serverStr    = "server"
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
		Data:   fmt.Sprint(e),
	}
}

func psUnknown(cmd string) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "Unknown command '" + cmd + "'",
	}
}

func psMessage(msg string) Message {
	return Message{
		Sender: pseudoStr,
		Data:   msg,
	}
}

func psAlert(msg string) Message {
	return Message{
		Type:   alertStr,
		Sender: pseudoStr,
		Data:   msg,
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

func editStateMsg(msg string) Message {
	return Message{
		Type:   editStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

// var nwmessages = interface{
// 	NoTeam: Message{"error:", "pseudoServer", "You have no team"},
// }
