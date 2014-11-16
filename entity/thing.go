package entity

import (
	"github.com/jasdel/explore/util/uid"
	"strings"
)

type Aliases []string

func (a Aliases) Match(name string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] == name {
			return true
		}
	}
	return false
}

type ThingInterface interface {
	uid.Interface
	Name() string
	Desc() string
	Aliases() Aliases
	IsAlias(string) bool
	OmitSelf() []ThingInterface
}

type Thing struct {
	uid.UID
	name    string
	desc    string
	aliases Aliases

	selfOmit []ThingInterface
}

func NewThingNoAliases(id uid.UID, name, desc string) *Thing {
	return NewThing(id, name, desc, Aliases{})
}

// Returns a new Thing from the values provided
func NewThing(id uid.UID, name, desc string, aliases Aliases) *Thing {
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
func (t *Thing) Aliases() Aliases {
	a := make(Aliases, len(t.aliases))
	copy(a, t.aliases)
	return a
}

// Returns true if the alias matches one of the aliases
// for this thing. the alias is trimmed, and lowercased
// before comparison.
func (t *Thing) IsAlias(alias string) bool {
	return t.aliases.Match(strings.ToLower(alias))
}

// Returns a pre-build omit interface list so
// one doesn't need to be created for broadcasts
func (t *Thing) OmitSelf() []ThingInterface {
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
