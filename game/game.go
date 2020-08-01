package game

type Game struct{}

var (
	g *Game
)

func NewGame() *Game {
	g = &Game{}
	return g
}
