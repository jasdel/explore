package main

import (
	"jasdel/explore/entity/actor"
	"jasdel/explore/entity/actor/player"
	"jasdel/explore/entity/location"
	"jasdel/explore/entity/location/area"
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/uid"
)

func main() {
	doneCh := make(chan struct{})
	a := area.New(doneCh)
	go a.Run()

	second := location.New(thing.New(<-uid.Next, "Second", "Second area", []string{}), nil)
	a.AddLoc(second)

	start := location.New(thing.New(<-uid.Next, "Start", "Starting area", []string{}), location.Exits{location.Exit{Name: "south", To: second}})
	a.AddLoc(start)

	act := actor.New(thing.New(<-uid.Next, "Actor", "Place holder actor", []string{"act"}), nil)
	a.MoveActor(act, start, true)

	p := player.StdInPlayer{Player: player.New(thing.New(<-uid.Next, "You", "its you silly", []string{}), nil)}
	go p.ReadStdIn()
	a.MoveActor(p, start, true)

	<-doneCh
}
