package unit

type MovementType int

type AttackType int

type Unit struct {
	HP        int `json:"hp"`
	MoveRange int `json:"moveRange"`
}
