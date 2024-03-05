package models

import (
	"math/rand"
)

type Unit struct {
	Type       string   `json:"string"`
	HP         int      `json:"hp"`
	MoveRange  int      `json:"moveRange"`
	Intiative  int      `json:"intiative"`
	IsAlive    bool     `json:"isAlive"`
	IsDeployed bool     `json:"isDeployed"`
	Traits     []string `json:"traits"`
	Row        int      `json:"row"`
	Col        int      `json:"col"`
}

var (
	ThreeCostPool []string = []string{"necromancer", "elder_sage", "reflection", "barbarian", "stranger"}
	TwoCostPool   []string = []string{"werewolf", "dancer", "sentinel", "demolitionist", "survivor"}
	OneCostPool   []string = []string{"hedge_knight", "trapper", "rider", "archer", "apprentice"}

	PrimaryAbilityLookupTable map[string]func() = map[string]func(){
		"necromancer": func() {

		},
		"sentinel": func() {

		},
		"hedge_knight": func() {

		},
	}

	SecondaryAbilityLookupTable map[string]func() = map[string]func(){}
)

func newUnit(unitType string, hp int, moveRange int, initiative int, traits []string) Unit {
	return Unit{Type: unitType, HP: hp, MoveRange: moveRange, Intiative: initiative, Traits: traits, IsAlive: true}
}

func UnitFactory(unitType string) Unit {
	switch unitType {
	case "necromancer":
		return newUnit(unitType, 4, 2, 1, []string{"Undead", "Mage"})
	case "sentinel":
		return newUnit(unitType, 5, 2, 0, []string{"Human", "Knight"})
	case "hedge_knight":
		return newUnit(unitType, 4, 2, 0, []string{"Human", "Knight"})
	default:
		return newUnit(unitType, 4, 2, 0, []string{"Human", "Knight"})
	}
}

// Generates unit pools in order of cost
func GenerateRoster() [][]string {
	var roster [][]string
	for i := 0; i < 3; i++ {
		roster = append(roster, []string{})
	}

	for i := 0; i < 5; i++ {
		roster[0] = append(roster[0], OneCostPool[rand.Intn(len(OneCostPool))])
		roster[1] = append(roster[1], TwoCostPool[rand.Intn(len(TwoCostPool))])
		roster[2] = append(roster[2], ThreeCostPool[rand.Intn(len(ThreeCostPool))])
	}

	return roster
}
