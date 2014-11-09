package entity

import (
	"github.com/jasdel/explore/entity/actor/player"
	"github.com/jasdel/explore/entity/location"
)

type World struct {
	players []player.Interface
	locs    []location.Interface
}

// func (w *World) AddPlayer(p player.Interface) {
// 	if findPlayerByName(p.Name(), w.players) == -1 {
// 		w.players = append(w.players, p)
// 	}
// }

// func (w *World) RemovePlayer(p player.Interface) {
// 	if idx := findPlayerByName(p.Name(), w.players); idx != -1 {
// 		w.players = append(w.players[:idx], w.players[idx+1:]...)
// 	}
// }

// func findPlayerByName(name string, players []player.Interface) int {
// 	for i := 0; i < len(players); i++ {
// 		if players[i].Name() == name {
// 			return i
// 		}
// 	}
// 	return -1
// }

// func (w *World) AddLoc(l location.Interface) {
// 	if findLocByName(l.Name(), w.locs) == -1 {
// 		w.locs = append(w.locs, l)
// 	}
// }

// func (w *World) RemoveLoc(l location.Interface) {
// 	if idx := findLocByName(l.Name(), w.locs); idx != -1 {
// 		w.locs = append(w.locs[:idx], w.locs[idx+1:]...)
// 	}
// }

// func findLocByName(name string, locs []location.Interface) int {
// 	for i := 0; i < len(locs); i++ {
// 		if locs[i].Name() == name {
// 			return i
// 		}
// 	}
// 	return -1
// }
