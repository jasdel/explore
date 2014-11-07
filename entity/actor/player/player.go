package player

import (
	"jasdel/explore/entity/actor"
	"jasdel/explore/entity/location"
	"jasdel/explore/entity/thing"
)

type Interface interface {
	actor.Interface
}

type Player struct {
	*actor.Actor
}

// Creates a new player instance
func New(t *thing.Thing, atLoc location.Interface) *Player {
	return &Player{
		Actor: actor.New(t, atLoc),
	}
}
