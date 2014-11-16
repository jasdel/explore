package entity

import (
	"fmt"
	"strings"
)

const (
	DirectionalExitMsgFmt  = `%[1]s leaves to the %[2]s`
	DirectionalEnterMsgFmt = `%[1]s enters from the %[2]s`
)

// Directional exit that is tied to an location
type Exit struct {
	Name      string
	Aliases   []string
	ExitMsg   string
	EnterMsg  string
	Loc       LocationInterface
	LookAhead string
}

// Processes the command determining if this exit is where the thing is going through
// Expects to be called in the same context as a location
//
// TODO need to refactor how locatables are moved between locations. This method
// 'works', but may not be safe.
//
func (e *Exit) Process(cmd *Command) bool {

	if e.Name == cmd.Verb {
		e.exit(cmd.Issuer)
		return true
	}

	for _, alias := range e.Aliases {
		if cmd.Verb == alias {
			e.exit(cmd.Issuer)
			return true
		}
	}

	return false
}

// Moves the thing out of its current location and into a new location.
func (e *Exit) exit(t ThingInterface) {
	locatable, ok := t.(Locatable)
	if !ok {
		fmt.Printf("Exit.exit: DEBUG: %s is not a locatable. %#v\n", e.Name, t)
		return
	}

	origLoc := locatable.Relocate(e.Loc)
	if origLoc != nil {
		origLoc.Broadcast(t.SelfOmit(), e.ExitMsg, t.Name(), e.Name)
		origLoc.Remove(t)
	}
	fmt.Println("Exit.Process: DEBUG:", e.Name, "relocated", t.Name(), t.UniqueId(), origLoc.Name(), origLoc.UniqueId())

	e.Loc.MoveIn(t, origLoc, fmt.Sprintf(e.EnterMsg, t.Name(), e.Name))
}

type Exits []Exit

// Prints the known exits to strings
func (e Exits) String() string {
	buf := make([]string, 0, len(e))

	for _, exit := range e {
		buf = append(buf, exit.Name)
	}
	output := strings.Join(buf, ", ")

	if output == "" {
		output = "none"
	}

	return output
}
