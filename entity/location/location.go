package location

import (
	"fmt"
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/inventory"
	"jasdel/explore/util/messaging"
	"jasdel/explore/util/uid"
)

// Locateable defines the interface for something that has a location or can be
// moved from or to a location. For example a mobile.
type Locateable interface {
	// Relocates a Locateable to a new location,
	// returning the original location
	Relocate(Interface) Interface

	// Locate gets a Locateable's current location
	Locate() Interface
}

// interface for all locations
type Interface interface {
	thing.Interface
	inventory.Interface
	command.Processor
	messaging.Broadcaster

	// Exits() Exits
	// LinkExit(e Exit)

	// A Thing is moving into this location
	MoveTo(thing thing.Interface, enterMsg string)
	Command(*command.Command)
}

// Specific location on the
type Location struct {
	*thing.Thing
	inventory.Inventory
	exits Exits

	moveCh chan thingMove
	cmdCh  chan *command.Command
}

// Creates a new area
func New(id uid.UID, name, desc string) *Location {
	return &Location{
		Thing:  thing.New(id, name, desc, []string{}),
		moveCh: make(chan thingMove, 10),
		cmdCh:  make(chan *command.Command, 10),
	}
}

// Returns the list of exists
func (l *Location) Exits() Exits {
	e := make(Exits, len(l.exits))
	copy(e, l.exits)
	return e
}

// Adds a new exit to the location
//
// TODO de-dupe exits
func (l *Location) LinkExit(e Exit) {
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

type thingMove struct {
	thing    thing.Interface
	enterMsg string
	spawn    bool
}

// Adds a thing to be moved away from this location to another
func (l *Location) MoveTo(t thing.Interface, enterMsg string) {
	if _, ok := t.(Locateable); !ok {
		fmt.Println("Location.Move: DEBUG:", l.Name(), "Thing", t.Name(), t.UniqueId(), "is not a Locatable")
		return
	}

	l.moveCh <- thingMove{
		thing:    t,
		enterMsg: enterMsg,
	}
}

// Spawns the thing into the
func (l *Location) Spawn(t thing.Interface) {
	if _, ok := t.(Locateable); !ok {
		fmt.Println("Location.Spawn: DEBUG:", l.Name(), "Thing", t.Name(), t.UniqueId(), "is not a Locatable")
		return
	}

	messaging.Respond(t, "You look around dazed as a swirl of smoke fades around you.")
	l.moveCh <- thingMove{
		thing:    t,
		enterMsg: fmt.Sprintf("%s appears in a puff of smoke looking dazed and confused.", t.Name()),
		spawn:    true,
	}
}

// Sends the command to be processed by this location
func (l *Location) Command(cmd *command.Command) {
	l.cmdCh <- cmd
}

// Runs the location
func (l *Location) Run(doneCh chan struct{}) {
	for {
		select {
		case cmd := <-l.cmdCh:
			fmt.Println("Location.Run: DEBUG:", l.Name(), "command received from", cmd.Issuer.Name(), cmd.Issuer.UniqueId(), cmd.Statement)
			l.Process(cmd)

		case move := <-l.moveCh:
			l.Add(move.thing)
			l.Broadcast([]thing.Interface{move.thing}, move.enterMsg)
			locatable := move.thing.(Locateable)
			locatable.Relocate(l)
			l.Process(command.New(move.thing, "look"))

		case _, ok := <-doneCh:
			if !ok {
				return
			}
		}
	}
}

// Processes the command within the scope of the area
func (l *Location) Process(cmd *command.Command) bool {
	// Give exits first chance to run so actors can
	// leave the area
	for _, e := range l.exits {
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
	if loc, ok := cmd.Issuer.(Locateable); ok {
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
	inv := thing.StringList(l.List(cmd.Issuer))

	cmd.Respond("You see: %s\n\nObvious exits: %s\n\n%s", l.Desc(), l.exits.String(), inv)
	return true
}

// lists all known exists
func (l *Location) listExists(cmd *command.Command) bool {
	cmd.Respond("Visible exits:\n%s", l.exits.String())
	return true
}
