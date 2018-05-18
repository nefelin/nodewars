package protocol

import (
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
			if room, ok := d.locations[m.Sender]; ok {
				err := room.Recv(m)
				if err != nil {
					err := d.Recv(m)
					if err != nil {
						m.Sender.Outgoing <- nwmessage.PsError(err)
					}
				}
			} else {
				err := d.Recv(m)
				if err != nil {
					m.Sender.Outgoing <- nwmessage.PsError(err)
				}
			}
			m.Sender.Outgoing <- nwmessage.PsPrompt(m.Sender.Prompt())

		}
	}
}
