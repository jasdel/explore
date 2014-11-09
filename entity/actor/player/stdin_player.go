package player

import (
	"bufio"
	"fmt"
	"github.com/jasdel/explore/util/command"
	"os"
	"strings"
	"time"
)

type StdInPlayer struct {
	*Player
}

func (p *StdInPlayer) ReadStdIn() {
	s := bufio.NewScanner(os.Stdin)

	t := time.NewTicker(time.Millisecond * 100)
	defer t.Stop()

	for s.Scan() {
		cmd := strings.TrimSpace(s.Text())
		if l := p.Locate(); l != nil && cmd != "" {
			l.Command(command.New(p, cmd))
		}
		// <-t.C
	}
}

// Implements the messaging.Responder interface for sending messages
// back to this player
func (p StdInPlayer) Respond(format string, any ...interface{}) {
	fmt.Printf(format+"\n", any...)
}
