package actor

import (
	"jasdel/explore/entity/location"
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/inventory"
	"jasdel/explore/util/messaging"
	"sync"
)

type Interface interface {
	thing.Interface
	inventory.Interface
	location.Locateable
	messaging.Broadcaster
}

type Actor struct {
	*thing.Thing
	inventory.Inventory
	atLoc location.Interface

	locateMtx sync.Mutex
}

// Create a new actor
func New(t *thing.Thing, atLoc location.Interface) *Actor {
	return &Actor{
		Thing: t,
		atLoc: atLoc,
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
func (a *Actor) Relocate(loc location.Interface) {
	a.locateMtx.Lock()
	defer a.locateMtx.Unlock()

	a.atLoc = loc
}

// Processes the actors specific commands
func (a *Actor) Process(cmd *command.Command) bool {
	return false
}

// Broadcasts the message to the location, if the actor is in that location.
func (a *Actor) Broadcast(omit []thing.Interface, format string, any ...interface{}) {
	if a.atLoc != nil {
		a.atLoc.Broadcast(omit, format, any...)
	}
}
