package nwmessage

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
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
	graphStateStr  = "graphState"

	pseudoStr    = "pseudoServer"
	serverStr    = "server"
	noConnectStr = "No connection"

	terminatorStr = "\n"
)

func PsPrompt(c chan Message, ws *websocket.Conn, m string) string {
	question := Message{
		Type:   confirmStr,
		Sender: pseudoStr,
		Data:   m,
	}

	// pose question
	c <- question

	// wait for response
	var res Message

	err := ws.ReadJSON(&res)
	if err != nil {
		log.Printf("error: %v", err)
		return "error"
	}

	return strings.ToLower(res.Data)
}

func PsError(e error) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   fmt.Sprint(e),
	}
}

// PS prefixed messages are printed to the users pseudoterminal
func PsUnknown(cmd string) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "Unknown command '" + cmd + "'" + terminatorStr,
	}
}

// PsNeutral returns a typeless message, pseudo terminal then prints without prefix
func PsNeutral(msg string) Message {
	return Message{
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func PsAlert(msg string) Message {
	return Message{
		Type:   alertStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func PsSuccess(msg string) Message {
	return Message{
		Type:   successStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func PsBegin(msg string) Message {
	return Message{
		Type:   beginStr,
		Sender: pseudoStr,
		Data:   msg + terminatorStr,
	}
}

func PsChat(msg string, context string) Message {
	return Message{
		Type:   context,
		Data:   msg,
		Sender: pseudoStr,
	}
}

func PsNoTeam() Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "No team assignment",
	}
}

func PsNoConnection() Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "No connection",
	}
}

// messages with server as Sender trigger action in the front end but are not show in the pseudoterminal

func AlertFlash(color string) Message {
	return Message{
		Type:   "alertFlash",
		Sender: serverStr,
		Data:   color,
	}
}

func EditState(msg string) Message {
	return Message{
		Type:   editStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func PromptState(msg string) Message {
	return Message{
		Type:   promptStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func GraphState(msg string) Message {
	return Message{
		Type:   graphStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func GraphReset() Message {
	return Message{
		Type:   "graphReset",
		Sender: serverStr,
		Data:   "",
	}
}
