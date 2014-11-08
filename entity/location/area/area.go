package area

import (
	"fmt"
	"jasdel/explore/entity/actor"
	"jasdel/explore/entity/location"
	"jasdel/explore/entity/thing"
	"jasdel/explore/util/command"
	"jasdel/explore/util/messaging"
)

type actorMove struct {
	actor actor.Interface
	toLoc location.Interface
	spawn bool
}

// Processes a group of locations within a single go routine
// the Area has complete ownership of all functionality of
// any thing within any location within.  Including items
// actors, and locations
type Area struct {
	name   string
	actors []actor.Interface
	locs   []location.Interface

	addLocCh chan location.Interface
	cmdCh    chan *command.Command
	moveCh   chan actorMove

	doneCh <-chan struct{}
}

// Creates a new instance of the area
func New(name string, doneCh <-chan struct{}) *Area {
	return &Area{
		name:     name,
		addLocCh: make(chan location.Interface, 1),
		cmdCh:    make(chan *command.Command, 10),
		moveCh:   make(chan actorMove, 10),
		doneCh:   doneCh,
	}
}

// Runs the area
//
// TODO does not check if location already added on add
// TODO do more than just send command to the location, verify ownership of location
// TODO need to check to make sure this area controls the location
// TODO does not support moving locations
//
func (a *Area) Run() {
	for {
		select {
		case loc := <-a.addLocCh:
			a.locs = append(a.locs, loc)
			loc.SetLocator(a)
		case cmd := <-a.cmdCh:
			if actor, ok := cmd.Issuer.(actor.Interface); ok {
				if loc := actor.Locate(); loc != nil {
					loc.Process(cmd)
				}
			}
		case move := <-a.moveCh:
			move.actor.Relocate(move.toLoc)
			move.toLoc.Add(move.actor)

			fmt.Println("DEBUG: area:", a.name, " move", move.actor.Name(), "to", move.toLoc.Name(), move.spawn)

			if move.spawn {
				move.toLoc.Broadcast([]thing.Interface{move.actor}, "Out of thin air %s appears in a puff of smoke looking dazed and confused.", move.actor)
				messaging.Respond(move.actor, "You look around dazed as a swirl of smoke fades around you.")
			} else {
				move.toLoc.Broadcast([]thing.Interface{move.actor}, "%s enters", move.actor)
			}

			move.toLoc.Process(command.New(move.actor, "look"))
		case _, ok := <-a.doneCh:
			if !ok {
				return
			}
		}
	}
}

// Moves an actor from one location to another using exits
// Implements the locator Move interface
func (a *Area) Move(t thing.Interface, loc location.Interface) {
	actor, ok := t.(actor.Interface)
	if !ok {
		return
	}

	a.moveCh <- actorMove{
		actor: actor,
		toLoc: loc,
		spawn: false,
	}
}

// Inserts an actor into a location
// implements the locator spawn interface
func (a *Area) Spawn(t thing.Interface, loc location.Interface) {
	actor, ok := t.(actor.Interface)
	if !ok {
		return
	}

	a.moveCh <- actorMove{
		actor: actor,
		toLoc: loc,
		spawn: true,
	}
}

// Adds a given location to the area
func (a *Area) AddLoc(loc location.Interface) {
	a.addLocCh <- loc
}

// Sends a command to the area
func (a *Area) Command(cmd *command.Command) {
	a.cmdCh <- cmd
}
