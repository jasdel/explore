// Structure and comments originally copied from Andrew 'Diddymus' Rolfe WolfMUD (2012)

// Definition of the messaging interface
package entity

import (
	"fmt"
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
func Respond(t ThingInterface, format string, any ...interface{}) {
	if responder, ok := t.(Responder); ok {
		responder.Respond(format, any...)
	} else {
		fmt.Printf("messaging.Respond: DEBUG: %s %d, is not a Responder. %#v\n", t.Name(), t.UniqueId(), t)
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
//	cmd.Broadcast([]entity.ThingInterface{p}, "You see %s sneeze.", cmd.Issuer.Name())
//	PlayerList.Broadcast(p.Locate().List(), "You hear a loud sneeze.")
type Broadcaster interface {
	Broadcast(omit []ThingInterface, format string, any ...interface{})
}

// Helper method to simplify broadcasts to a thing
func Broadcast(t ThingInterface, omit []ThingInterface, format string, any ...interface{}) {
	if broadcaster, ok := t.(Broadcaster); ok {
		broadcaster.Broadcast(omit, format, any...)
	}
}
