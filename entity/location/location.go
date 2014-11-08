package location

import (
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/inventory"
	"jasdel/explore/util/messaging"
	"sync"
)

// Directional exit that is tied to an location
type Exit struct {
	Name      string
	Aliases   []string
	ExitMsg   string
	To        Interface
	LookAhead string
}

// Processes the command determining if this exit is where the thing is going through
// Expects to be called in the same context as a location
//
// TODO need to refactor how locatables are moved between locations. This method
// 'works', but may not be safe.
//
func (e *Exit) Process(cmd *command.Command) bool {
	for _, alias := range e.Aliases {
		if cmd.Verb == alias {
			locateable, ok := cmd.Issuer.(Locateable)
			if !ok {
				return false
			}

			if loc := locateable.Locate(); loc != nil {
				loc.Broadcast([]thing.Interface{cmd.Issuer}, "%s %s", cmd.Issuer.Name(), e.ExitMsg)
				loc.Remove(cmd.Issuer)
				locateable.Relocate(nil)
			}

			e.To.Locator().Move(cmd.Issuer, e.To)
			return true
		}
	}
	return false
}

type Exits []Exit

// Prints the known exits to strings
func (e Exits) String() string {
	var output string
	for i, exit := range e {
		output += exit.Name
		if i+1 != len(e) {
			output += ", "
		}
	}

	if output == "" {
		output = "none"
	}

	return output
}

// Locateable defines the interface for something that has a location or can be
// moved from or to a location. For example a mobile.
type Locateable interface {
	Relocate(Interface) // Relocates a Locateable to a new location
	Locate() Interface  // Locate gets a Locateable's current location
}

// Locator provides an interface for moving things between locations
type Locator interface {
	Command(*command.Command)
	Spawn(thing.Interface, Interface)
	Move(thing.Interface, Interface)
}

// interface for all locations
type Interface interface {
	thing.Interface
	inventory.Interface
	command.Processor
	messaging.Broadcaster

	Exits() Exits
	Locator() Locator
	SetLocator(Locator)
	LinkExit(e Exit)
}

// Specific location on the
type Location struct {
	*thing.Thing
	inventory.Inventory
	exits   Exits
	locator Locator

	locatorMtx sync.Mutex
}

// Creates a new area
func New(t *thing.Thing, exits Exits) *Location {
	return &Location{
		Thing: t,
		exits: exits,
	}
}

// Returns the list of exists
func (l *Location) Exits() Exits {
	e := make(Exits, len(l.exits))
	copy(e, l.exits)
	return e
}

// Returns the locator for this location
func (l *Location) Locator() Locator {
	l.locatorMtx.Lock()
	defer l.locatorMtx.Unlock()

	locator := l.locator
	return locator
}

// Switches to a new locator
func (l *Location) SetLocator(locator Locator) {
	l.locatorMtx.Lock()
	defer l.locatorMtx.Unlock()

	l.locator = locator
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

// Processes the command within the scope of the area
func (l *Location) Process(cmd *command.Command) bool {
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
