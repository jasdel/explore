package entity

import (
	"github.com/jasdel/explore/util/uid"
	"strings"
)

type ThingInterface interface {
	uid.Interface
	Name() string
	Desc() string
	Aliases() []string
	IsAlias(string) bool
	SelfOmit() []ThingInterface
}

type Thing struct {
	uid.UID
	name    string
	desc    string
	aliases []string

	selfOmit []ThingInterface
}

func NewThingNoAliases(id uid.UID, name, desc string) *Thing {
	return NewThing(id, name, desc, []string{})
}

// Returns a new Thing from the values provided
func NewThing(id uid.UID, name, desc string, aliases []string) *Thing {
	t := &Thing{
		UID:     id,
		name:    name,
		desc:    desc,
		aliases: aliases,
	}
	t.selfOmit = []ThingInterface{t}

	return t
}

// Return the thing's name
func (t *Thing) Name() string {
	return t.name
}

// Returns the thing's description
func (t *Thing) Desc() string {
	return t.desc
}

// Returns a copy of the aliases for this thing
func (t *Thing) Aliases() []string {
	a := make([]string, len(t.aliases))
	copy(a, t.aliases)
	return a
}

// Returns true if the alias matches one of the aliases
// for this thing. the alias is trimmed, and lowercased
// before comparison.
func (t *Thing) IsAlias(alias string) bool {
	aliasLower := strings.ToLower(alias)
	for i := 0; i < len(t.aliases); i++ {
		if t.aliases[i] == aliasLower {
			return true
		}
	}
	return false
}

// Returns a pre-build omit interface list so
// one doesn't need to be created for broadcasts
func (t *Thing) SelfOmit() []ThingInterface {
	return t.selfOmit
}

// Utility method to convert list of things to string
func ThingsToString(things []ThingInterface) string {
	var output string
	for _, t := range things {
		output += t.Name() + "\n"
	}
	return output
}