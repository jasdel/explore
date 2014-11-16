package entity

import (
	"strings"
)

// Processor should be implemented by anything that wants to process or react
// to commands. These commands are usually from players and actors but also
// commonly from other objects. For example a lever when pulled may issue an
// 'OPEN DOOR' command to get the door associated with the lever to open.
type Processor interface {
	Process(*Command) (handled bool)
}

type Command struct {
	Issuer    ThingInterface
	Statement []string // words in the command
	Verb      string   // first word in the statement. e.g get, examine, attack, ...
	Nouns     []string // 2nd : nth words, usually what is being
	Target    string   // alias for second word, since its normally the target
}

// Creates a new instance of a command, with the statement parsed.
func NewCommand(issuer ThingInterface, statement string) *Command {
	cmd := &Command{Issuer: issuer}
	cmd.Parse(statement)
	return cmd
}

// Parses the provided statement into the command. Replacing the
// current verb, nous, target, and statement values.
//
// TODO Allow quotes in commands to group words together
//
func (c *Command) Parse(statment string) {
	c.Statement = strings.Fields(statment)
	if l := len(c.Statement); l > 0 {
		c.Verb, c.Nouns = c.Statement[0], c.Statement[1:]
		if l > 1 {
			c.Target = c.Statement[1]
		} else {
			c.Target = ""
		}
	} else {
		c.Verb, c.Nouns, c.Target = "", []string{}, ""
	}
}

// Sends a message back directly to the issuer of the command
func (c *Command) Respond(format string, any ...interface{}) {
	Respond(c.Issuer, format, any...)
}
