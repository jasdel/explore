package location

import (
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/inventory"
	"jasdel/explore/util/messaging"
)

// Directional exit that is tied to an location
type Exit struct {
	Name      string
	To        Interface
	LookAhead string
}

type Exits []Exit

// Onterface for all locations
type Interface interface {
	thing.Interface
	inventory.Interface
	command.Processor
	messaging.Broadcaster

	Exits() Exits
}

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

// Specific location on the
type Location struct {
	*thing.Thing
	inventory.Inventory
	exits Exits
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
	// Give the inventory a chance to process the command first
	// By definition all actions occur in the context
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

	cmd.Respond("You see\n%s\nVisible exits:\n%s\n%s", l.Desc(), l.exits.String(), inv)
	return true
}

// lists all known exists
func (l *Location) listExists(cmd *command.Command) bool {
	cmd.Respond("Visible exits:\n%s", l.exits.String())
	return true
}
