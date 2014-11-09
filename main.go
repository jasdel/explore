package main

import (
	"jasdel/explore/entity/actor"
	"jasdel/explore/entity/actor/player"
	"jasdel/explore/entity/location"
	"jasdel/explore/util/uid"
)

func main() {
	doneCh := make(chan struct{})

	start := location.New(<-uid.Next, "City Center", "Central location with in the city")
	south := location.New(<-uid.Next, "Southern plains", "Southern plains ")
	east := location.New(<-uid.Next, "Market District", "Busy street full of merchants selling their goods.")

	start.LinkExit(location.Exit{Name: "south", Aliases: []string{"south", "s"}, Loc: south,
		ExitMsg: location.DirectionalExitMsgFmt, EnterMsg: location.DirectionalEnterMsgFmt,
	})
	south.LinkExit(location.Exit{Name: "north", Aliases: []string{"north", "n"}, Loc: start,
		ExitMsg: location.DirectionalExitMsgFmt, EnterMsg: location.DirectionalEnterMsgFmt,
	})
	start.LinkExit(location.Exit{Name: "east", Aliases: []string{"east", "e"}, Loc: east,
		ExitMsg: location.DirectionalExitMsgFmt, EnterMsg: location.DirectionalEnterMsgFmt,
	})
	east.LinkExit(location.Exit{Name: "west", Aliases: []string{"west", "w"}, Loc: start,
		ExitMsg: location.DirectionalExitMsgFmt, EnterMsg: location.DirectionalEnterMsgFmt,
	})

	go start.Run(doneCh)
	go south.Run(doneCh)
	go east.Run(doneCh)

	p := player.StdInPlayer{Player: player.New(<-uid.Next, "You", "its you silly")}
	start.Spawn(p)

	act := actor.New(<-uid.Next, "Place Holder Actor", "Place holder actor", []string{"actor"})
	start.Spawn(act)

	go p.ReadStdIn()

	<-doneCh
}
