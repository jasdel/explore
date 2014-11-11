package actor

import (
	"fmt"
	"github.com/jasdel/explore/entity/location"
	"github.com/jasdel/explore/entity/thing"
	"github.com/jasdel/explore/util/command"
	"github.com/jasdel/explore/util/inventory"
	"github.com/jasdel/explore/util/messaging"
	"github.com/jasdel/explore/util/uid"
	"strings"
	"sync"
)

type Interface interface {
	thing.Interface
	inventory.Interface
	location.Locatable
	messaging.Broadcaster
}

type Actor struct {
	*thing.Thing
	inventory.Inventory
	atLoc location.Interface

	locateMtx sync.Mutex
}

// Create a new actor
func New(id uid.UID, name, desc string, aliases []string) *Actor {
	return &Actor{
		Thing: thing.New(id, name, desc, aliases),
	}
}

// Returns the location the Actor is currently at
func (a *Actor) Locate() location.Interface {
	a.locateMtx.Lock()
	defer a.locateMtx.Unlock()

	loc := a.atLoc

	return loc
}

// Updates the location of the actor
func (a *Actor) Relocate(loc location.Interface) location.Interface {
	a.locateMtx.Lock()
	defer a.locateMtx.Unlock()

	origLoc := a.atLoc
	a.atLoc = loc

	return origLoc
}

// Processes the actors specific commands
func (a *Actor) Process(cmd *command.Command) bool {
	switch cmd.Verb {
	case "inventory", "inv", "i":
		inv := thing.SliceToString(a.List(cmd.Issuer))
		cmd.Respond(inv)
	case "say":
		statement := strings.TrimSpace(strings.Join(cmd.Nouns, " "))
		if statement != "" {
			a.Broadcast(a.SelfOmit(), "%s says %s", a.Name(), statement)
		}
	case "tell", "whisper":
		// TODO handle cross thing communication.
		// start with location first, region, then branch out to world
	default:
		// Give the inventory a chance to process the command
		if a.Inventory.Process(cmd) {
			return true
		}

		// No match, so this command is not handled by the actor
		return false
	}

	return true
}

// Broadcasts the message to the location, if the actor is in that location.
func (a *Actor) Broadcast(omit []thing.Interface, format string, any ...interface{}) {
	if a.atLoc != nil {
		a.atLoc.Broadcast(omit, format, any...)
	} else {
		fmt.Println("Actor.Broadcast: DEBUG:", a.Name(), "tried to broadcast, but is not at a location")
	}
}
