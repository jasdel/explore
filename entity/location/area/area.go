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
	start bool
}

// Processes a group of locations within a single go routine
// the Area has complete ownership of all functionality of
// any thing within any location within.  Including items
// actors, and locations
type Area struct {
	actors []actor.Interface
	locs   []location.Interface

	addLocCh    chan location.Interface
	cmdCh       chan *command.Command
	moveActorCh chan actorMove

	doneCh <-chan struct{}
}

// Creates a new instance of the area
func New(doneCh <-chan struct{}) *Area {
	return &Area{
		addLocCh:    make(chan location.Interface, 1),
		cmdCh:       make(chan *command.Command, 10),
		moveActorCh: make(chan actorMove, 10),
		doneCh:      doneCh,
	}
}

// Runs the area
//
// TODO does not check if location already added on add
// TODO do more than just send command to the location, verify ownership of location
// TODO need to check to make sure this area controls the location
//
func (a *Area) Run() {
	for {
		select {
		case loc := <-a.addLocCh:
			a.locs = append(a.locs, loc)
		case cmd := <-a.cmdCh:
			if actor, ok := cmd.Issuer.(actor.Interface); ok && actor.Locate() != nil {
				actor.Locate().Process(cmd)
			}
		case actorMove := <-a.moveActorCh:
			actorMove.actor.Relocate(actorMove.toLoc)
			actorMove.toLoc.Add(actorMove.actor)

			fmt.Println("actorMove", actorMove.toLoc.Name(), actorMove.actor.Name(), actorMove.start)

			if actorMove.start {
				actorMove.toLoc.Broadcast([]thing.Interface{actorMove.actor}, "Out of thin air %s appears in a puff of smoke looking dazed and confused.", actorMove.actor)
				messaging.Respond(actorMove.actor, "You look around dazed as a swirl of smoke fades around you.")
			} else {
				actorMove.toLoc.Broadcast([]thing.Interface{actorMove.actor}, "%s enters", actorMove.actor)
			}

			actorMove.toLoc.Process(command.New(actorMove.actor, "look"))
		case _, ok := <-a.doneCh:
			if !ok {
				return
			}
		}
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

// Moves the actor from their current location into a new location.
// TODO should this move the actor from the origin?
func (a *Area) MoveActor(actor actor.Interface, loc location.Interface, start bool) {
	a.moveActorCh <- actorMove{
		actor: actor,
		toLoc: loc,
		start: start,
	}
}
