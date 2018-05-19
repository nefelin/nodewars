package receiver

import "nwmessage"

// Receiver interface is implemented by any struct able to process client messsages
type Receiver interface {
	Recv(msg nwmessage.ClientMessage) error
}
