package entity

import (
	"github.com/jasdel/explore/util/uid"
)

type Item struct {
	*Thing
	weight float64
}

func NewItem(id uid.UID, name, desc string, aliases Aliases) *Item {
	return &Item{
		Thing: NewThing(id, name, desc, aliases),
	}
}

// Processes the command on the item.
// examples would be a player picking up, or dropping
// this item from their inventory.
func (i *Item) Process(cmd *Command) bool {
	// This specific item?
	if !i.IsAlias(cmd.Target) {
		return false
	}

	switch cmd.Verb {
	// case "junk", "trash":
	case "drop":
		return i.drop(cmd)
	case "get", "pickup":
		return i.pickup(cmd)
	case "examine", "exam", "ex":
		return i.examine(cmd)
	case "weigh":
		return i.weigh(cmd)
	}
	return false
}

// Estimating the weight of the item.
// TODO need better output for respond
func (i *Item) weigh(cmd *Command) bool {
	cmd.Respond("You weigh %s, and it feels about %.2f.", i.Name(), i.weight)
	cmd.Broadcast(cmd.Issuer.OmitSelf(), "%s tries to determine %s's weight.", cmd.Issuer.Name(), i.Name())
	return true
}

// Examining the item
func (i *Item) examine(cmd *Command) bool {
	cmd.Respond("You look closely at %s, and you see %s", i.Name(), i.Desc())
	cmd.Broadcast(cmd.Issuer.OmitSelf(), "%s examines %s.", cmd.Issuer.Name(), i.Name())
	return true
}

// Drops an item from the inventory to the location
// Requires:
// - thing be in a location
// - thing must have an inventory
// - and this item was in the inventory
func (i *Item) drop(cmd *Command) bool {
	if loc, ok := cmd.Issuer.(Locatable); ok {
		if inv, ok := cmd.Issuer.(InventoryInterface); ok {
			if inv.Contains(i) {
				inv.Remove(i)
				loc.Locate().Add(i)

				cmd.Respond("You drop %s.", i.Name())
				cmd.Broadcast(cmd.Issuer.OmitSelf(), "You see %s drop %s.", cmd.Issuer.Name(), i.Name())
				return true
			}
		}
	} else {
		cmd.Respond("You don't see anywhere to drop %s.", i.Name())
		return true
	}

	return false
}

// Pickups an item from a location and places it in the things inventory
//
// TODO: does not work with handle containers
//
// Requires:
// - thing be in a location
// - thing must have an inventory
// - and this item was in the inventory
func (i *Item) pickup(cmd *Command) bool {
	if loc, ok := cmd.Issuer.(Locatable); ok && loc.Locate() != nil {
		if inv, ok := cmd.Issuer.(InventoryInterface); ok {
			if loc.Locate().Contains(i) {
				loc.Locate().Remove(i)
				inv.Add(i)

				cmd.Respond("You pickup %s.", i.Name())
				cmd.Broadcast(cmd.Issuer.OmitSelf(), "You see %s pickup %s.", cmd.Issuer.Name(), i.Name())
				return true
			}
		}
	}

	return false
}
