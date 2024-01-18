package snake

import (
	"fmt"
	"image/color"

	"github.com/casen/snakegame/model"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
	boardRows    = 20
	boardCols    = 20
)

var (
	backgroundColor = color.RGBA{50, 100, 50, 50}
	snakeColor      = color.RGBA{0, 255, 0, 255}
	foodColor       = color.RGBA{200, 200, 50, 150}
)

type Game struct {
	board *Board
}

func NewGame() *Game {
	return &Game{
		board: NewGameBoard(boardRows, boardCols),
	}
}

func (g *Game) Update(action model.Vector) error {
	return g.board.Update(action)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)
	if g.board.gameOver {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Game Over. Score: %d", g.board.points))
	} else {
		width := ScreenHeight / boardRows

		for _, p := range g.board.snake.body {
			vector.DrawFilledRect(screen, float32(p.Y*width), float32(p.X*width), float32(width), float32(width), snakeColor, true)
		}
		vector.DrawFilledRect(screen, float32(g.board.food.Y*width), float32(g.board.food.X*width), float32(width), float32(width), foodColor, true)
		ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.board.points))
	}
}
func (g *Game) GameOver() bool {
	return g.board.GameOver()
}

func (g *Game) Reset() {
	g.board = NewGameBoard(boardRows, boardCols)
}

func (g *Game) Score() int {
	return g.board.points
}

func (g *Game) CurrentDirection() model.Vector {
	return g.board.snake.direction
}

func (g *Game) CurrentLocation() model.Point {
	return g.board.snake.Head()
}

func (g *Game) FoodLocation() model.Point {
	return g.board.food
}

func (g *Game) PrintGameBoard() {
	g.board.Print()
}

func (g *Game) MoveIsValid(dir model.Vector) bool {
	nextLocation := g.board.NextLocation(dir, 1)
	isValid := g.board.MoveIsValid(nextLocation)
	return isValid
}

func (g *Game) Move(dir model.Vector) {
	g.board.Move(dir)
}

func (g *Game) EvaluateAction(action model.Vector) (reward float32, isDone bool) {
	return g.board.EvaluateAction(action)
}

func (g *Game) CurrentState() [11]float32 {
	return g.board.CurrentState()
}

func (g *Game) NextState(action model.Vector) [11]float32 {
	return g.board.NextState(action)
}

func (g *Game) PlayerPosition() model.Point {
	return g.board.snake.Head()
}
