package location

import (
	"fmt"
	"jasdel/explore/util/command"
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
	Loc       Interface
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
			locatable, ok := cmd.Issuer.(Locatable)
			if !ok {
				return false
			}

			fmt.Println("Exit.Process: DEBUG:", e.Name, "exit by", cmd.Issuer.Name(), cmd.Issuer.UniqueId())

			origLoc := locatable.Relocate(e.Loc)
			if origLoc != nil {
				origLoc.Broadcast(cmd.Issuer.SelfOmit(), e.ExitMsg, cmd.Issuer.Name(), e.Name)
				origLoc.Remove(cmd.Issuer)
			}
			fmt.Println("Exit.Process: DEBUG:", e.Name, "relocated", cmd.Issuer.Name(), cmd.Issuer.UniqueId(), origLoc.Name(), origLoc.UniqueId())

			e.Loc.MoveIn(cmd.Issuer, origLoc, fmt.Sprintf(e.EnterMsg, cmd.Issuer.Name(), e.Name))
			return true
		}
	}
	return false
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
