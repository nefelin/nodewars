package nwmessage

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Message is our basic message struct
type Message struct {
	Type   string `json:"type"`
	Sender string `json:"sender"`
	Data   string `json:"data"`
	// Code   string `json:"code"`
}

const (
	alertStr       = "alert"
	errorStr       = "error"
	successStr     = "success"
	beginStr       = "begin"
	dialogueMsgStr = "dialogue"

	challengeStateStr   = "challengeState"
	compOutStateStr     = "compOutState"
	editLangStateStr    = "editorLangState"
	langSupportStateStr = "langSupportState"
	editStateStr        = "editorState"
	graphStateStr       = "graphState"
	scoreStateStr       = "scoreState"
	StdinStateStr       = "stdinState"
	teamStateStr        = "teamState"
	graphResetStr       = "graphReset"
	resultStateStr      = "resultState"
	terminalPauseStr    = "pauseTerm"
	terminalUnpauseStr  = "unpauseTerm"

	pseudoStr = "pseudoServer"
	serverStr = "server"

	noConnectStr = "No connection"

	terminatorStr = "\n"
	preStr        = ""
)

func PsDialogue(msg string) Message {
	return Message{
		Type:   dialogueMsgStr,
		Sender: pseudoStr,
		Data:   msg,
	}
}

func PsPrompt(msg string) Message {
	return Message{
		Sender: pseudoStr,
		Data:   "\n" + msg + " ",
	}
}

func PsError(e error) Message {
	return Message{
		Type:   errorStr,
		Sender: pseudoStr,
		Data:   preStr + fmt.Sprint(e) + terminatorStr,
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
		Data:   preStr + msg + terminatorStr,
	}
}

func PsAlert(msg string) Message {
	return Message{
		Type:   alertStr,
		Sender: pseudoStr,
		Data:   preStr + msg + terminatorStr,
	}
}

func PsSuccess(msg string) Message {
	return Message{
		Type:   successStr,
		Sender: pseudoStr,
		Data:   preStr + msg + terminatorStr,
	}
}

func PsBegin(msg string) Message {
	return Message{
		Type:   beginStr,
		Sender: pseudoStr,
		Data:   preStr + msg + terminatorStr,
	}
}

func PsChat(sender, context, msg string) Message {
	return Message{
		Type:   "",
		Data:   preStr + fmt.Sprintf("%s (%s): %s", sender, context, msg) + terminatorStr,
		Sender: pseudoStr,
	}
}

func PsNoTeam() Message {
	return PsError(errors.New("No team assignment"))
}

func PsCompileFail() Message {
	return PsError(errors.New("Compile failed"))
}

func PsNoConnection() Message {
	return PsError(errors.New("No connection"))
}

// messages with server as Sender trigger action in the front end but are not show in the pseudoterminal
func ChallengeState(c interface{}) Message {
	msg, _ := json.Marshal(c)
	return Message{
		Type:   challengeStateStr,
		Sender: serverStr,
		Data:   string(msg),
	}
}

func EditState(msg string) Message {
	return Message{
		Type:   editStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func EditLangState(msg string) Message {
	return Message{
		Type:   editLangStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

func LangSupportState(msg []string) Message {
	bytes, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return Message{
		Type:   langSupportStateStr,
		Sender: serverStr,
		Data:   string(bytes),
	}
}

func StdinState(msg string) Message {
	return Message{
		Type:   StdinStateStr,
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

func ResultState(r interface{}) Message {
	msg, _ := json.Marshal(r)
	return Message{
		Type:   resultStateStr,
		Sender: serverStr,
		Data:   string(msg),
	}
}

func ScoreState(msg string) Message {
	return Message{
		Type:   scoreStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}

// pause and resume are not behaving client-side
func TerminalPause() Message {
	return Message{
		Type:   terminalPauseStr,
		Sender: serverStr,
	}
}

func TerminalUnpause() Message {
	return Message{
		Type:   terminalUnpauseStr,
		Sender: serverStr,
	}
}

func TeamState(msg string) Message {
	return Message{
		Type:   teamStateStr,
		Sender: serverStr,
		Data:   msg,
	}
}
