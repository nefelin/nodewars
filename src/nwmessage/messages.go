package nwmessage

import (
	"fmt"
)

// Message is our basic message struct
type Message struct {
	Type   string `json:"type"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
	Code   string `json:"code"`
} // TODO fix code submission to append to other data, this is unnecessary

const (
	alertStr       = "alert"
	errorStr       = "error"
	successStr     = "success"
	beginStr       = "begin"
	dialogueMsgStr = "dialogue"

	editStateStr   = "editorState"
	promptStateStr = "promptState"
	graphStateStr  = "graphState"
	teamStateStr   = "teamState"
	graphResetStr  = "graphReset"
	// startDialogueStr = "startDialogue"
	// endDialogueStr   = "endDialogue"

	pseudoStr = "pseudoServer"
	serverStr = "server"

	noConnectStr = "No connection"

	terminatorStr = "\n"
)

func PsDialogue(msg string) Message {
	return Message{
		Type:   dialogueMsgStr,
		Sender: pseudoStr,
		Data:   msg,
	}
}

func PsError(e error) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   fmt.Sprint(e) + terminatorStr,
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
		Sender: pseudoStr + terminatorStr,
	}
}

func PsNoTeam() Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "No team assignment" + terminatorStr,
	}
}

func PsNoConnection() Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   "No connection" + terminatorStr,
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

func GraphState(msg string) Message {
	return Message{
		Type:   graphStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func GraphReset() Message {
	return Message{
		Type:   graphResetStr,
		Sender: serverStr,
		Data:   "",
	}
}

func PromptState(msg string) Message {
	return Message{
		Type:   promptStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

// pause and resume are not behaving client-side
// func StartDialogue() Message {
// 	return Message{
// 		Type:   startDialogueStr,
// 		Sender: serverStr,
// 	}
// }

// func EndDialogue() Message {
// 	return Message{
// 		Type:   startDialogueStr,
// 		Sender: serverStr,
// 	}
// }

func TeamState(msg string) Message {
	return Message{
		Type:   teamStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}
