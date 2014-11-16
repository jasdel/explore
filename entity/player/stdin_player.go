package player

import (
	"bufio"
	"fmt"
	"github.com/jasdel/explore/entity"
	"os"
	"strings"
)

type StdInPlayer struct {
	*Player
}

func (p *StdInPlayer) ReadStdIn() {
	s := bufio.NewScanner(os.Stdin)

	for s.Scan() {
		cmd := strings.TrimSpace(s.Text())
		if l := p.Locate(); l != nil && cmd != "" {
			l.Command(entity.NewCommand(p, cmd))
		}
	}
}

// Implements the messaging.Responder interface for sending messages
// back to this player
func (p *StdInPlayer) Respond(format string, any ...interface{}) {
	fmt.Printf(format+"\n", any...)
}
