package player

import (
	"github.com/jasdel/explore/entity"
	"github.com/jasdel/explore/util/uid"
)

type Player struct {
	*entity.Actor
}

// Creates a new player instance
func NewPlayer(id uid.UID, name, desc string) *Player {
	return &Player{
		Actor: entity.NewActor(id, name, desc, entity.Aliases{}),
	}
}
