package main

import (
	// "github.com/davecheney/profile"
	"github.com/jasdel/explore/entity/actor"
	"github.com/jasdel/explore/entity/actor/player"
	"github.com/jasdel/explore/entity/location"
	"github.com/jasdel/explore/entity/location/region"
	"github.com/jasdel/explore/entity/thing"
	"github.com/jasdel/explore/util/uid"
)

//
// TODO add signal handling for os.Interrupt, os.Kill
// TODO clean shutdown, right now done ch just kills
// oll the goroutines, need to wait for them to empty
// and shutdown command processing.
//
func main() {
	// defer profile.Start(profile.CPUProfile).Stop()

	doneCh := make(chan struct{})

	r1 := region.New("r1")
	go r1.Run(doneCh)
	r2 := region.New("r2")
	go r2.Run(doneCh)

	start := location.New(<-uid.Next, "City Center", "Central location with in the city", r1.CmdCh, r1.MoveCh)
	r1.AddLoc(start)

	south := location.New(<-uid.Next, "Southern plains", "Plains as far as the eyes can see.", r2.CmdCh, r2.MoveCh)
	r2.AddLoc(south)

	east := location.New(<-uid.Next, "Market District", "Busy street full of merchants selling their goods.", r1.CmdCh, r1.MoveCh)
	r1.AddLoc(east)

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

	p := &player.StdInPlayer{Player: player.New(<-uid.Next, "You", "its you silly")}
	p.Add(thing.New(<-uid.Next, "Spoon", "a old four foot long wooden spoon", []string{"actor"}))

	act := actor.New(<-uid.Next, "Place Holder Actor", "Place holder actor", []string{"actor"})

	start.Spawn(p)
	start.Spawn(act)

	go p.ReadStdIn()

	<-doneCh
}
