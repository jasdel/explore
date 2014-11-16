package npc

import (
	"github.com/jasdel/explore/entity"
	"github.com/jasdel/explore/util/uid"
)

type NPC struct {
	*entity.Actor
}

// Creates a new player instance
func NewNPC(id uid.UID, name, desc string) *NPC {
	return &NPC{
		Actor: entity.NewActor(id, name, desc, []string{}),
	}
}
