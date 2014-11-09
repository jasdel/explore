package player

import (
	"github.com/jasdel/explore/entity/actor"
	"github.com/jasdel/explore/util/uid"
)

type Interface interface {
	actor.Interface
}

type Player struct {
	*actor.Actor
}

// Creates a new player instance
func New(id uid.UID, name, desc string) *Player {
	return &Player{
		Actor: actor.New(id, name, desc, []string{}),
	}
}
