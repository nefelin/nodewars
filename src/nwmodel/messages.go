package nwmodel

import (
	"fmt"
	"log"
	"strings"
)

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
	confirmStr = "confirmation"
	// yesnoStr   = "(y/n)"

	editStateStr   = "editorState"
	promptStateStr = "promptState"

	pseudoStr    = "pseudoServer"
	serverStr    = "server"
	noConnectStr = "No connection"

	terminatorStr = "\n"
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

func psPrompt(p *Player, m string) string {
	question := Message{
		Type:   confirmStr,
		Sender: pseudoStr,
		Data:   m,
	}

	// pose question
	p.outgoing <- question

	// wait for response
	var res Message

	err := p.socket.ReadJSON(&res)
	if err != nil {
		log.Printf("error: %v", err)
		return "error"
	}

	return strings.ToLower(res.Data)
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
		Data:   "Unknown command '" + cmd + "'" + terminatorStr,
	}
}

func psMessage(msg string) Message {
	return Message{
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func psAlert(msg string) Message {
	return Message{
		Type:   alertStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func psSuccess(msg string) Message {
	return Message{
		Type:   successStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func psBegin(msg string) Message {
	return Message{
		Type:   beginStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func editStateMsg(msg string) Message {
	return Message{
		Type:   editStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func promptStateMsg(msg string) Message {
	return Message{
		Type:   promptStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

// var nwmessages = interface{
// 	NoTeam: Message{"error:", "pseudoServer", "You have no team"},
// }
