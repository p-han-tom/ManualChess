package match

import (
	"fmt"
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
	Speed                int      `json:"speed"`
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

	PrimaryAbilityLookupTable = map[string]func(playerId string, casterId string, game *Match, targets []Position) error{
		"necromancer":  NecromancerPrimary,
		"sentinel":     SentinelPrimary,
		"hedge_knight": HedgeKnightPrimary,
	}

	SecondaryAbilityLookupTable map[string]func() = map[string]func(){}
)

func validatePositions(positions []Position) error {
	for _, pos := range positions {
		if pos.Row >= BoardHeight || pos.Row < 0 || pos.Col >= BoardWidth || pos.Col < 0 {
			return fmt.Errorf("invalid position: %d, col %d", pos.Row, pos.Col)
		}
	}
	return nil
}

func UnitFactory(unitType string) Unit {
	switch unitType {
	case "necromancer":
		return Unit{
			Type:                 "necromancer",
			PrimaryAbilityName:   "Soul Sap",
			SecondaryAbilityName: "Summon Skeleton",
			HP:                   4,
			MoveRange:            2,
			Speed:                1,
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
			Speed:                0,
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
			Speed:                2,
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
			Speed:                2,
			Traits:               []string{"Human", "Knight"},
			IsAlive:              true,
		}
	}
}
