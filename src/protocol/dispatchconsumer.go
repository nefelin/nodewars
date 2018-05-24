package protocol

import (
	"fmt"
	"model/player"
	"nwmessage"
)

func dispatchConsumer(d *Dispatcher) {
	for {
		select {
		// if we get a new player, register and pass back
		// to the connection handler
		case regReq := <-d.registrationQueue:
			// fmt.Println("Lobby relaying reg")
			d.handleRegRequest(regReq)

		// if we get a player command, handle that
		case m := <-d.clientMessages:
			p := m.Sender.(*player.Player) // TODO use Client instead of player everywhere we can....

			switch m.Type {
			case "editorState":
				p.SetEditor(m.Data, false)
			case "stdinState":
				p.SetStdin(m.Data, false)
			case "playerCmd":
				if room, ok := d.locations[p]; ok {
					err := room.Recv(m)
					if err != nil {
						err := d.Recv(m)
						if err != nil {
							m.Sender.Outgoing(nwmessage.PsError(err))

						}
					}
				} else {
					err := d.Recv(m)
					if err != nil {
						m.Sender.Outgoing(nwmessage.PsError(err))

					}
				}
				p.Outgoing(nwmessage.PsPrompt(p.Prompt()))
			default:
				errStr := fmt.Sprintf("Unknown message type, '%s'", m.Type)
				fmt.Println(errStr)
				p.Outgoing(nwmessage.ServerError(errStr))
			}

		}
	}
}

func (d *Dispatcher) Recv(m nwmessage.ClientMessage) error {
	if m.Data == "" {
		return nil
	}

	return dispatchCommands.Exec(d, m)
}
