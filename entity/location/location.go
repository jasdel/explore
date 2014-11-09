package location

import (
	"fmt"
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/inventory"
	"jasdel/explore/util/messaging"
	"jasdel/explore/util/uid"
	"sync"
)

// Locatable defines the interface for something that has a location or can be
// moved from or to a location. For example a mobile.
type Locatable interface {
	// Relocates a Locatable to a new location,
	// returning the original location
	Relocate(Interface) Interface

	// Locate gets a Locatable's current location
	Locate() Interface
}

// interface for all locations
type Interface interface {
	thing.Interface
	inventory.Interface
	command.Processor
	messaging.Broadcaster

	Exits() Exits
	LinkExit(e Exit)

	// A Thing is moving into this location
	MoveIn(thing thing.Interface, from Interface, enterMsg string)
	Command(*command.Command)
}

// Specific location on the
type Location struct {
	*thing.Thing
	inventory.Inventory
	exits Exits

	moveCh chan ThingMove
	cmdCh  chan *command.Command

	linkMtx sync.Mutex
}

// Creates a new area
func New(id uid.UID, name, desc string, cmdCh chan *command.Command, moveCh chan ThingMove) *Location {
	return &Location{
		Thing:  thing.New(id, name, desc, []string{}),
		moveCh: moveCh,
		cmdCh:  cmdCh,
	}
}

// Returns the list of exists
func (l *Location) Exits() Exits {
	l.linkMtx.Lock()
	defer l.linkMtx.Unlock()

	e := make(Exits, len(l.exits))
	copy(e, l.exits)
	return e
}

// Adds a new exit to the location
//
// TODO de-dupe exits
//
func (l *Location) LinkExit(e Exit) {
	l.linkMtx.Lock()
	defer l.linkMtx.Unlock()

	l.exits = append(l.exits, e)
}

// Broadcast sends a message to all responders at this location. This
// implements the broadcast.Interface - see that for more details.
func (l *Location) Broadcast(omit []thing.Interface, format string, any ...interface{}) {
	for _, t := range l.Inventory.List(omit...) {
		if responder, ok := t.(messaging.Responder); ok {
			responder.Respond(format, any...)
		}
	}
}

type ThingMove struct {
	Thing    thing.Interface
	ToLoc    Interface
	EnterMsg string
	Spawn    bool
}

// Moves a thing from this location from another
func (l *Location) MoveIn(t thing.Interface, from Interface, enterMsg string) {
	if _, ok := t.(Locatable); !ok {
		fmt.Println("Location.Move: DEBUG:", l.Name(), l.UniqueId(), "thing", t.Name(), t.UniqueId(), "is not locatable")
		return
	}

	l.moveCh <- ThingMove{
		Thing:    t,
		ToLoc:    l,
		EnterMsg: enterMsg,
	}
}

// Spawns the thing into the
func (l *Location) Spawn(t thing.Interface) {
	locatable, ok := t.(Locatable)
	if !ok {
		fmt.Println("Location.Spawn: DEBUG:", l.Name(), l.UniqueId(), "thing", t.Name(), t.UniqueId(), "is not locatable")
		return
	}

	locatable.Relocate(l)

	messaging.Respond(t, "You look around dazed as a swirl of smoke fades around you.")
	l.moveCh <- ThingMove{
		Thing:    t,
		ToLoc:    l,
		EnterMsg: fmt.Sprintf("%s appears in a puff of smoke looking dazed and confused.", t.Name()),
		Spawn:    true,
	}
}

// Sends the command to be processed by this location
func (l *Location) Command(cmd *command.Command) {
	l.cmdCh <- cmd
}

// Processes the command within the scope of the area
func (l *Location) Process(cmd *command.Command) bool {
	// Give exits first chance to run so actors can leave the area
	for _, e := range l.Exits() {
		if e.Process(cmd) {
			return true
		}
	}

	// Give the inventory a chance to process the command
	if l.Inventory.Process(cmd) {
		return true
	}

	// The following commands can only be processed relative to the
	// issuer's location. So we need to check if this location is
	// where the issuer is.
	//
	// TODO is this actually needed? explore's model requires issuer
	// to be in same location, and command, might be useful for
	// error checking.
	//
	if loc, ok := cmd.Issuer.(Locatable); ok {
		if !loc.Locate().IsAlso(l) {
			return false
		}
	}

	switch cmd.Verb {
	case "look":
		return l.look(cmd)
	case "exits":
		return l.listExists(cmd)
	}
	return false
}

// Responds with the location's description and known visible exits
func (l *Location) look(cmd *command.Command) bool {
	inv := thing.SliceToString(l.List(cmd.Issuer))

	cmd.Respond("You see: %s\n\nObvious exits: %s\n\n%s", l.Desc(), l.Exits().String(), inv)
	return true
}

// lists all known exists
func (l *Location) listExists(cmd *command.Command) bool {
	cmd.Respond("Visible exits:\n%s", l.Exits().String())
	return true
}
