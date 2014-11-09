package location

import (
	"fmt"
	"jasdel/explore/entity/thing"
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
			locateable, ok := cmd.Issuer.(Locateable)
			if !ok {
				return false
			}

			if loc := locateable.Relocate(nil); loc != nil {
				loc.Broadcast([]thing.Interface{cmd.Issuer}, e.ExitMsg, cmd.Issuer.Name(), e.Name)
				loc.Remove(cmd.Issuer)
			}

			e.Loc.MoveTo(cmd.Issuer, fmt.Sprintf(e.EnterMsg, cmd.Issuer.Name(), e.Name))
			return true
		}
	}
	return false
}

type Exits []Exit

// Prints the known exits to strings
func (e Exits) String() string {
	var buf []string
	for _, exit := range e {
		buf = append(buf, exit.Name)
	}
	output := strings.Join(buf, ", ")

	if output == "" {
		output = "none"
	}

	return output
}
