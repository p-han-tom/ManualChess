package match

import (
	"fmt"
	"math"
)

func NecromancerPrimary(playerId string, casterId string, game *Match, targets []Position) error {
	err := validatePositions(targets)
	if err != nil {
		return err
	}

	if len(targets) != 1 {
		return fmt.Errorf("invalid number of targets, expected: 1, actual: %d", len(targets))
	}

	var targetPos Position = targets[0]
	var player *Player
	var opp *Player
	if game.Player1.ID == playerId {
		player = &game.Player1
		opp = &game.Player2
	} else {
		player = &game.Player2
		opp = &game.Player1
	}

	caster := player.Units[casterId]
	var casterPos Position = caster.Pos

	if (targetPos.Row == casterPos.Row && targetPos.Col == casterPos.Col) || (targetPos.Row != casterPos.Row && targetPos.Col != casterPos.Col) {
		return fmt.Errorf("invalid target, casterPos: %d, %d, target: %d, %d", casterPos.Row, casterPos.Col, targetPos.Row, targetPos.Col)
	}

	var distanceToTarget int = targetPos.Row - casterPos.Row + targetPos.Col - casterPos.Col
	if math.Abs(float64(distanceToTarget)) > 4 {
		return fmt.Errorf("target is too far, distance must be less than or equal to 4")
	}

	var targetId string = game.Board[targetPos.Row][targetPos.Col].OccupantId
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
}

func NecromancerSecondary() {

}
