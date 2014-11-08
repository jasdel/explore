package player

import (
	"bufio"
	"fmt"
	"jasdel/explore/util/command"
	"os"
)

type StdInPlayer struct {
	*Player
}

func (p *StdInPlayer) ReadStdIn() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if l := p.Locate(); l != nil {
			l.Locator().Command(command.New(p, s.Text()))
		}
	}
}

// Implements the messaging.Responder interface for sending messages
// back to this player
func (p StdInPlayer) Respond(format string, any ...interface{}) {
	fmt.Printf(format+"\n", any...)
}
