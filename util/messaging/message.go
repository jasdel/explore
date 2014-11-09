// Structure and comments copied from Andrew 'Diddymus' Rolfe WolfMUD (2012)

// Definition of the messaging interface
package messaging

import (
	"fmt"
	"jasdel/explore/entity/thing"
)

// Respond should be implemented by anything that wants to 'respond' to players.
// It is modeled after fmt.Printf so that messages can easily be built with
// parameters. For example:
//
//	cmd.Respond("You go %s.", directionLongNames[d])
type Responder interface {
	Respond(format string, any ...interface{})
}

// Helper method to simplify sending a respond to a thing.
func Respond(t thing.Interface, format string, any ...interface{}) {
	if responder, ok := t.(Responder); ok {
		responder.Respond(format, any...)
	} else {
		fmt.Println("messaging.Respond: DEBUG:", t.Name(), t.UniqueId(), "is not a Responder.")
	}
}

// Broadcast should be implemented by anything that wants to send messages to
// multiple responders. This is usually to everyone currently in the world
// or at a specific location. Like responders the function is modeled after
// fmt.Printf and takes messages formatted in the same way. The omit parameter
// is used to omit certain responders. For example if a player sneezes in a
// location they would have a different message and be omitted from the broadcast
// to the location and the sneezer and people in the location would be omitted
// from the message broadcast to the world:
//
//	cmd.Respond("You sneeze. Aaahhhccchhhooo!")
//	cmd.Broadcast([]thing.Interface{p}, "You see %s sneeze.", cmd.Issuer.Name())
//	PlayerList.Broadcast(p.Locate().List(), "You hear a loud sneeze.")
type Broadcaster interface {
	Broadcast(omit []thing.Interface, format string, any ...interface{})
}
