package player

import (
	"jasdel/explore/entity/actor"
	"jasdel/explore/util/uid"
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
