package npc

import (
	"github.com/jasdel/explore/entity"
	"math/rand"
	"time"
)

type RndNPC struct {
	*NPC
}

func (n *RndNPC) Run() {
	for {
		<-time.After(time.Duration(rand.Intn(10)) * time.Second)
		exits := n.Locate().Exits()
		n.Locate().Command(entity.NewCommand(n, exits[rand.Intn(len(exits))].Name))
	}
}

func (n *RndNPC) Respond(format string, any ...interface{}) {
	// TODO need better parsing style string parsing for each entity is silly
}
