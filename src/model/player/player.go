package player

import (
	"challenges"
	"fmt"
	"nwmessage"
	"strconv"

	"github.com/gorilla/websocket"
)

// TODO un export all but route
type PlayerID = int

var playerIDCount PlayerID

// Player ...
type Player struct {
	ID       PlayerID `json:"id"`
	name     string   `json:"name"`
	TeamName string   `json:"team"`
	// Route    *route                 `json:"route"`
	socket   *websocket.Conn        `json:"-"`
	outgoing chan nwmessage.Message `json:"-"`
	language string                 // current working language

	macAddress string // address of the machine player is current attached to

	compiling bool // this is used to block player action while submitted code is compiling
	chatMode  bool // track whether player is in chatmode or not (for use in lobby)
	inGame    bool // is player in a game?

	editorState string
	stdinState  string // stdin buffer for testing
	// termState terminalState
}

// implement client interface

func (p *Player) Outgoing(m nwmessage.Message) {
	p.outgoing <- m
}

func (p *Player) ChatMode() bool {
	return p.chatMode
}

func (p *Player) ToggleChat() {
	p.chatMode = !p.chatMode
}

func (p *Player) Socket() *websocket.Conn {
	return p.socket
}

func (p *Player) Cleanup() {
	p.socket.Close()
	close(p.outgoing)
}

// player methods -------------------------------------------------------------------------------
// TODO this is in the wrong place
func NewPlayer(ws *websocket.Conn) *Player {
	ret := &Player{
		ID:          playerIDCount,
		name:        "",
		socket:      ws,
		outgoing:    make(chan nwmessage.Message),
		editorState: "",
		stdinState:  "",
	}

	go outgoingRelay(ret)

	// log.Println("New player created, setting language...")
	playerIDCount++
	return ret
}

func outgoingRelay(p *Player) {
	for {
		if msg, ok := <-p.outgoing; ok { // if channel is open...
			if err := p.socket.WriteJSON(msg); err != nil { // try writing message to player, complain if we have problems
				// fmt.Printf("error dispatching message: '%v',\n to player '%s'\n", msg, p.GetName())
			}
		} else { // if channel is closed, player is gone.
			return
		}
	}
}

func (p *Player) Stdin() string {
	return p.stdinState
}

func (p *Player) Editor() string {
	return p.editorState
}

func (p *Player) SetStdin(s string, send bool) {
	p.stdinState = s
	if send {
		p.Outgoing(nwmessage.StdinState(p.stdinState))
	}
}

func (p *Player) SetEditor(s string, send bool) {
	p.editorState = s
	if send {
		p.Outgoing(nwmessage.EditState(p.editorState))
	}
}

// this is odd under player as it's really about FE messaging
func (p *Player) SetChallenge(c challenges.Challenge) {
	p.Outgoing(nwmessage.ChallengeState(c.Redacted()))

}

// func (p *Player) TeamName() string {
// 	// TODO send to front end
// 	return p.teamName
// }

// func (p *Player) SetTeamName(n string) {
// 	// TODO send to front end
// 	p.teamName = n
// }

func (p *Player) Language() string {
	// TODO send to front end
	return p.language
}

func (p *Player) SetLanguage(l string) {
	// TODO send to front end
	p.language = l
}

func (p *Player) MacAddress() string {
	// TODO send to front end
	return p.macAddress
}

func (p *Player) SetMacAddress(a string) {
	// TODO send to front end
	p.macAddress = a
}

func (p *Player) SendPrompt() {
	// p.Outgoing(nwmessage.PsPrompt(p.Prompt()))
	p.Outgoing(nwmessage.PsPrompt(">"))

}

// Prompt should be generated by the ROOM the player is in...
func (p *Player) Prompt() string {
	promptEndChar := ">"
	prompt := p.Name()

	if p.macAddress != "" {
		prompt += fmt.Sprintf(":%s", p.macAddress)
	}
	prompt += fmt.Sprintf("[%s]", p.language)

	if p.ChatMode() {
		prompt += "[CHATMODE]"
	}

	prompt += promptEndChar

	return prompt
}

func (p *Player) SubmitCode(id challenges.ChallengeID) (challenges.GradedResult, error) {
	response := challenges.SubmitTest(id, p.language, p.Editor())

	p.Outgoing(nwmessage.ResultState(response))

	if response.Message.Type == "error" {
		return response, fmt.Errorf("Compiled failed")
	}

	if response.Passed() == 0 {
		return response, fmt.Errorf("Solution failed all tests")
	}

	return response, nil
}

// GetName returns the players name if they have one, assigns one if they don't
func (p *Player) Name() string {
	if p.name == "" {
		p.SetName("player_" + strconv.Itoa(p.ID))
	}

	return p.name
}

func (p *Player) SetName(n string) {
	p.name = n
	p.Outgoing(nwmessage.PsSuccess("Name set to '" + n + "'"))

}

// hasTeam is deprecated I think TOD
func (p Player) HasTeam() bool {
	if p.TeamName == "" {
		return false
	}
	return true
}

func (p Player) String() string {
	return fmt.Sprintf("( <player> {Name: %v, team: %v} )", p.Name(), p.TeamName)
}

// Methods relying on coupling -------------------------------------------------------------------------------------------

// func (p *Player) Location() *node {
// 	if p.Route == nil {
// 		return nil
// 	}
// 	return p.Route.Endpoint()
// }

// func (p *Player) CanSubmit() error {
// 	mac := p.currentMachine()
// 	switch {
// 	case p.EditorState == "":
// 		return errors.New("No code to submit")
// 	case mac == nil:
// 		return errors.New("Not attached to a machine")
// 	case !mac.isNeutral() && !mac.belongsTo(p.TeamName) && mac.language != p.language:
// 		return fmt.Errorf("This machine is written in %s, your code must also be written in %s", mac.language, mac.language)
// 	}
// 	return nil

// }

// func (p *Player) MacDetach() {
// 	mac := p.currentMachine()
// 	p.challengeState(Challenge{})

// 	if mac != nil {
// 		mac.remPlayer(p)
// 		p.macAddress = ""
// 	}
// }

// func (p *Player) BreakConnection(forced bool) {
// 	if p.Route == nil {
// 		return
// 	}

// 	p.macDetach()
// 	p.Route = nil

// 	if forced {
// 		p.Outgoing(nwmessage.PsError(errors.New("Connection interrupted!")))

// 	}
// }

// // TODO refactor this, modify how slots are tracked, probably with IDs
// func (p *Player) CurrentMachine() *machine {
// 	if p.Route == nil || p.macAddress == "" {
// 		return nil
// 	}

// 	return p.Route.Endpoint().addressMap[p.macAddress]
// }
