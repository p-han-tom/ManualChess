package match

type PlayerColour int

const (
	Red PlayerColour = iota
	Blue
)

type Player struct {
	ID     string          `json:"id"`
	Colour PlayerColour    `json:"colour"`
	Units  map[string]Unit `json:"units"`
	Gold   int             `json:"gold"`
}
