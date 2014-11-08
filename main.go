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

	a2 := area.New(doneCh)
	go a2.Run()

	start := location.New(thing.New(<-uid.Next, "City Center", "Central location with in the city", []string{}), location.Exits{})
	south := location.New(thing.New(<-uid.Next, "Southern plains", "Southern plains ", []string{}), location.Exits{})
	east := location.New(thing.New(<-uid.Next, "Market District", "Busy street full of merchants selling their goods.", []string{}), location.Exits{})

	start.LinkExit(location.Exit{Name: "south", Aliases: []string{"south", "s"}, To: south})
	south.LinkExit(location.Exit{Name: "north", Aliases: []string{"north", "n"}, To: start})

	start.LinkExit(location.Exit{Name: "east", Aliases: []string{"east", "e"}, To: east})
	east.LinkExit(location.Exit{Name: "west", Aliases: []string{"west", "w"}, To: start})

	a.AddLoc(start)
	a.AddLoc(south)
	a2.AddLoc(east)

	act := actor.New(thing.New(<-uid.Next, "Place Holder Actor", "Place holder actor", []string{"act"}), nil)
	a.Spawn(act, start)

	p := player.StdInPlayer{Player: player.New(thing.New(<-uid.Next, "You", "its you silly", []string{}), nil)}
	go p.ReadStdIn()
	a.Spawn(p, start)

	<-doneCh
}
