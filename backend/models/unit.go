package models

import (
	"math/rand"
)

type MovementType int

type AttackType int

type Unit struct {
	HP        int      `json:"hp"`
	MoveRange int      `json:"moveRange"`
	IsAlive   bool     `json:"isAlive"`
	Traits    []string `json:"traits"`
}

var (
	ThreeCostPool []string = []string{"necromancer", "elder_sage", "reflection", "barbarian", "stranger"}
	TwoCostPool   []string = []string{"werewolf", "dancer", "sentinel", "demolitionist", "survivor"}
	OneCostPool   []string = []string{"hedge_knight", "trapper", "rider", "archer", "apprentice"}
)

// Generates unit pools in order of cost
func GenerateRoster() [][]string {
	var roster [][]string
	for i := 0; i < 3; i++ {
		roster = append(roster, []string{})
	}

	for i := 0; i < 5; i++ {
		roster[0] = append(roster[0], OneCostPool[rand.Intn(5)])
		roster[1] = append(roster[1], TwoCostPool[rand.Intn(5)])
		roster[2] = append(roster[2], ThreeCostPool[rand.Intn(5)])
	}

	return roster
}
