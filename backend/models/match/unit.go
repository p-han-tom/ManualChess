package match

import (
	"math/rand"
)

type Position struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

type Unit struct {
	Type                 string   `json:"string"`
	PrimaryAbilityName   string   `json:"primaryAbilityName"`
	SecondaryAbilityName string   `json:"secondaryAbilityName"`
	HP                   int      `json:"hp"`
	MoveRange            int      `json:"moveRange"`
	Intiative            int      `json:"intiative"`
	ActionPoints         int      `json:"ActionPoints"`
	IsAlive              bool     `json:"isAlive"`
	IsDeployed           bool     `json:"isDeployed"`
	Traits               []string `json:"traits"`
	Pos                  Position `json:"pos"`
	ClassData            any      `json:"classData"`
}

type NecromancerData struct {
	Souls int `json:"souls"`
}

var (
	ThreeCostPool []string = []string{"necromancer", "elder_sage", "reflection", "barbarian", "stranger"}
	TwoCostPool   []string = []string{"werewolf", "dancer", "sentinel", "demolitionist", "survivor"}
	OneCostPool   []string = []string{"hedge_knight", "trapper", "rider", "archer", "apprentice"}
)

func UnitFactory(unitType string) Unit {
	switch unitType {
	case "necromancer":
		return Unit{
			Type:                 "necromancer",
			PrimaryAbilityName:   "Soul Sap",
			SecondaryAbilityName: "Summon Skeleton",
			HP:                   4,
			MoveRange:            2,
			Intiative:            1,
			Traits:               []string{"Undead", "Mage"},
			IsAlive:              true,
			ClassData:            NecromancerData{Souls: 0},
		}
	case "sentinel":
		return Unit{
			Type:                 "sentinel",
			PrimaryAbilityName:   "Forceful Shove",
			SecondaryAbilityName: "Guard",
			HP:                   6,
			MoveRange:            2,
			Intiative:            0,
			Traits:               []string{"Dwarf", "Knight"},
			IsAlive:              true,
		}
	case "hedge_knight":
		return Unit{
			Type:                 "hedge_knight",
			PrimaryAbilityName:   "Axe Slash",
			SecondaryAbilityName: "Short Rest",
			HP:                   4,
			MoveRange:            2,
			Intiative:            2,
			Traits:               []string{"Human", "Knight"},
			IsAlive:              true,
		}
	default:
		return Unit{
			Type:                 "hedge_knight",
			PrimaryAbilityName:   "Axe Slash",
			SecondaryAbilityName: "Short Rest",
			HP:                   4,
			MoveRange:            2,
			Intiative:            2,
			Traits:               []string{"Human", "Knight"},
			IsAlive:              true,
		}
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
