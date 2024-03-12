package match

import (
	"fmt"
	"math"
)

var (
	PrimaryAbilityLookupTable = map[string]func(playerId string, casterId string, match *Match, targets []Position) error{
		"necromancer": func(playerId string, casterId string, match *Match, targets []Position) error {
			err := validateTargets(targets)
			if err != nil {
				return err
			}

			if len(targets) != 1 {
				return fmt.Errorf("invalid number of targets, expected: 1, actual: %d", len(targets))
			}

			var targetPos Position = targets[0]
			var player *Player
			var opp *Player
			if match.Player1.ID == playerId {
				player = &match.Player1
				opp = &match.Player2
			} else {
				player = &match.Player2
				opp = &match.Player1
			}

			caster := player.Units[casterId]
			var casterPos Position = caster.Pos

			if (targetPos.Row == casterPos.Row && targetPos.Col == casterPos.Col) || (targetPos.Row != casterPos.Row && targetPos.Col != casterPos.Col) {
				return fmt.Errorf("invalid target, casterPos: %d, %d, target: %d, %d", casterPos.Row, casterPos.Col, targetPos.Row, targetPos.Col)
			}

			var distanceToTarget int = targetPos.Row - casterPos.Row + targetPos.Col - casterPos.Col
			if math.Abs(float64(distanceToTarget)) >= 4 {
				return fmt.Errorf("target is too far, distance must be less than or equal to 4")
			}

			var targetId string = match.Board[targetPos.Row][targetPos.Col].OccupantId
			if targetId == "" {
				return fmt.Errorf("target tile has no occupant")
			}

			target, exists := opp.Units[targetId]
			if !exists {
				return fmt.Errorf("target unit is not an enemy unit")
			}

			target.HP -= 2
			caster.HP += 1
			classData := caster.ClassData.(NecromancerData)
			classData.Souls += 1
			if classData.Souls > 3 {
				classData.Souls = 3
			}
			caster.ClassData = classData
			player.Units[casterId] = caster
			opp.Units[targetId] = target

			return nil
		},
		"sentinel": func(playerId string, casterId string, match *Match, targets []Position) error {
			return nil
		},
		"hedge_knight": func(playerId string, casterId string, match *Match, targets []Position) error {
			return nil
		},
	}

	SecondaryAbilityLookupTable map[string]func() = map[string]func(){}
)

func validateTargets(targets []Position) error {
	for _, target := range targets {
		if target.Row >= BoardHeight || target.Row < 0 || target.Col >= BoardWidth || target.Col < 0 {
			return fmt.Errorf("invalid target row: %d, col %d", target.Row, target.Col)
		}
	}
	return nil
}
