package region

import (
	"fmt"
	"github.com/jasdel/explore/entity"
	"sync"
)

type Region struct {
	name string

	locs map[entity.LocationInterface]struct{}

	CmdCh  chan *entity.Command
	MoveCh chan entity.ThingMove

	locMtx sync.Mutex
}

func NewRegion(name string) *Region {
	return &Region{
		name:   name,
		locs:   make(map[entity.LocationInterface]struct{}),
		CmdCh:  make(chan *entity.Command, 10),
		MoveCh: make(chan entity.ThingMove, 10),
	}
}

// Adds a location to the region.  The location should
// of been initialized to use this region's cmd and
// move channels.
func (r *Region) AddLoc(loc entity.LocationInterface) {
	r.locMtx.Lock()
	defer r.locMtx.Unlock()

	if _, ok := r.locs[loc]; !ok {
		r.locs[loc] = struct{}{}
	}
}

// Returns true of the location is apart of this region.
func (r *Region) hasLoc(loc entity.LocationInterface) bool {
	r.locMtx.Lock()
	defer r.locMtx.Unlock()

	if _, ok := r.locs[loc]; ok {
		return true
	}
	return false
}

// Starts processing commands and move statements
// from locations and clients
func (r *Region) Run(doneCh chan struct{}) {
	for {
		select {
		case cmd := <-r.CmdCh:
			r.command(cmd)
		case move := <-r.MoveCh:
			r.move(move)
		case _, ok := <-doneCh:
			if !ok {
				return
			}
		}
	}
}

// Processes the command, forwarding it own to the thing that
// is able to process it. Starting first with the issuer, and
// falling down into the location.
//
// TODO Should commands/moves received for other regions be forwarded?
// 		How to prevent endless loop is so?
//
func (r *Region) command(cmd *entity.Command) {
	locatable, ok := cmd.Issuer.(entity.Locatable)
	if !ok {
		fmt.Println("Region.command: ERROR:", r.name, "received command without being locatable.", cmd.Issuer.Name(), cmd.Issuer.UniqueId())
	}

	loc := locatable.Locate()
	if loc == nil {
		fmt.Println("Region.command: ERROR:", r.name, "received command with nil location", cmd.Issuer.Name(), cmd.Issuer.UniqueId())
		return
	}

	fmt.Println("Region.command: DEBUG:", r.name, "command received from", cmd.Issuer.Name(), cmd.Issuer.UniqueId(), cmd.Statement, "at", loc.Name(), loc.UniqueId())

	if !r.hasLoc(loc) {
		fmt.Println("Region.command: WARN:", r.name, "received command for location not supported by this region, dropping.")
		return
	}

	// first let the issuer process the command.
	if processor, ok := cmd.Issuer.(entity.Processor); ok && processor.Process(cmd) {
		return
	}

	// send the command along to the location to process
	if loc.Process(cmd) {
		return
	}

	fmt.Println("Region.command: DEBUG:", r.name, "command received but not processed", cmd.Verb, "by", cmd.Issuer.Name(), cmd.Issuer.UniqueId())
}

// Moves a thing which is currently not attached to any location, into a
// location that is owned by this region.
//
// TODO ensure destination location belongs to this region
//
func (r *Region) move(move entity.ThingMove) {
	fmt.Println("Region.move: DEBUG:", r.name, "move received from", move.Thing.Name(), move.Thing.UniqueId(), "to", move.ToLoc.Name(), move.ToLoc.UniqueId())

	if !r.hasLoc(move.ToLoc) {
		fmt.Println("Region.move: WARN:", r.name, "received move for location not supported by this region, dropping.")
		return
	}

	move.ToLoc.Add(move.Thing)
	move.ToLoc.Broadcast(move.Thing.OmitSelf(), move.EnterMsg)
	move.ToLoc.Process(entity.NewCommand(move.Thing, "look"))
}
