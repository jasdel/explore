package main

import (
	"github.com/jasdel/explore/entity"
	"github.com/jasdel/explore/entity/player"
	"github.com/jasdel/explore/region"
	"github.com/jasdel/explore/util/uid"
	// "github.com/pkg/profile"
)

//
// TODO add signal handling for os.Interrupt, os.Kill
// TODO clean shutdown, right now done ch just kills
// oll the goroutines, need to wait for them to empty
// and shutdown command processing.
//
func main() {
	// defer profile.Start().Stop()

	doneCh := make(chan struct{})

	r1 := region.NewRegion("r1")
	go r1.Run(doneCh)
	r2 := region.NewRegion("r2")
	go r2.Run(doneCh)

	start := entity.NewLocation(<-uid.Next, "City Center", "Central location with in the city", r1.CmdCh, r1.MoveCh)
	r1.AddLoc(start)

	south := entity.NewLocation(<-uid.Next, "Southern plains", "Plains as far as the eyes can see.", r2.CmdCh, r2.MoveCh)
	r2.AddLoc(south)

	east := entity.NewLocation(<-uid.Next, "Market District", "Busy street full of merchants selling their goods.", r1.CmdCh, r1.MoveCh)
	r1.AddLoc(east)

	start.LinkExit(entity.Exit{Name: "south", Aliases: []string{"south", "s"}, Loc: south,
		ExitMsg: entity.DirectionalExitMsgFmt, EnterMsg: entity.DirectionalEnterMsgFmt,
	})
	south.LinkExit(entity.Exit{Name: "north", Aliases: []string{"north", "n"}, Loc: start,
		ExitMsg: entity.DirectionalExitMsgFmt, EnterMsg: entity.DirectionalEnterMsgFmt,
	})
	start.LinkExit(entity.Exit{Name: "east", Aliases: []string{"east", "e"}, Loc: east,
		ExitMsg: entity.DirectionalExitMsgFmt, EnterMsg: entity.DirectionalEnterMsgFmt,
	})
	east.LinkExit(entity.Exit{Name: "west", Aliases: []string{"west", "w"}, Loc: start,
		ExitMsg: entity.DirectionalExitMsgFmt, EnterMsg: entity.DirectionalEnterMsgFmt,
	})

	p := &player.StdInPlayer{Player: player.NewPlayer(<-uid.Next, "You", "its you silly")}
	p.Add(entity.NewItem(<-uid.Next, "Spoon", "a old four foot long wooden spoon", []string{"spoon"}))

	act := entity.NewActor(<-uid.Next, "Place Holder Actor", "Place holder actor", []string{"actor"})

	start.Spawn(p)
	start.Spawn(act)

	go p.ReadStdIn()

	<-doneCh
}
