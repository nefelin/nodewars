package model

// // Broadcast methods

// func (gm *GameModel) TeamBroadcast(t *team.Team, msg nwmessage.Message) {
// 	msg.Sender = "pseudoServer"

// 	for pID := range t.players {
// 		p, err := gm.pFromID(pID)
// 		if err != nil {
// 			panic(err)
// 		}
// 		p.Outgoing(msg)
// 	}
// }
