package snake

import (
	"math/rand"
	"time"

	"github.com/casen/snakegame/model"
)

type Board struct {
	rows     int
	cols     int
	food     *model.Point
	snake    *Snake
	points   int
	gameOver bool
	timer    time.Time
}

func NewBoard(rows int, cols int) *Board {

	board := &Board{
		rows:  rows,
		cols:  cols,
		timer: time.Now(),
	}
	// start in top-left corner
	board.snake = NewSnake([]model.Point{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}, {X: 0, Y: 3}}, model.Vector{X: 0, Y: 1})
	board.PlaceFood()

	return board
}

func (b *Board) Update(input *Input) error {
	if b.gameOver {
		return nil
	}

	// snake goes faster when there are more points
	interval := time.Millisecond * 150
	if b.points > 10 {
		interval = time.Millisecond * 125
	} else if b.points > 20 {
		interval = time.Millisecond * 100
	}

	if _, cardinalDir, ok := input.Dir(); ok {
		b.snake.ChangeDirection(cardinalDir)
	}

	if time.Since(b.timer) >= interval {
		if err := b.MoveSnake(); err != nil {
			return err
		}

		b.timer = time.Now()
	}

	return nil
}

func (b *Board) PlaceFood() {
	var x, y int
	var point model.Point

	for {
		x = rand.Intn(b.cols)
		y = rand.Intn(b.rows)
		point = model.Point{X: x, Y: y}

		// make sure we don't put a food on a snake
		if !b.snake.HitsSnake(point) {
			break
		}
	}

	b.food = &point
}

func (b *Board) MoveSnake() error {
	// remove tail first, add 1 in front
	b.snake.Move()
	snakeHead := b.snake.Head()

	if b.OutOfBounds(snakeHead.X, snakeHead.Y) || b.snake.HeadHitsBody() {
		b.gameOver = true
		return nil
	}

	if b.snake.HeadHits(*b.food) {
		// the snake grows on the next move
		b.snake.justAte = true

		b.PlaceFood()
		b.points++
	}

	return nil
}

// Programmatically move the snake, rather than take keyboard input from player
func (b *Board) Move(dir model.Vector) {
	b.snake.ChangeDirection(dir)
	b.MoveSnake()
}

func (b *Board) OutOfBounds(x, y int) bool {
	return x > b.cols-1 || y > b.rows-1 || x < 0 || y < 0
}

func (b *Board) NextLocation(dir model.Vector) model.Point {
	currentLocation := b.snake.Head()
	nextLocation := model.Point{X: currentLocation.X + dir.X, Y: currentLocation.Y + dir.Y}
	return nextLocation
}

func (b *Board) MoveIsValid(nextLocation model.Point) bool {
	return !b.OutOfBounds(nextLocation.X, nextLocation.Y) && !b.snake.ValidMove(nextLocation)
}

func (b *Board) MoveIsTerminal(point model.Point) bool {
	return b.OutOfBounds(point.X, point.Y) || b.snake.HitsSnake(point)
}

func (b *Board) MoveIsScoring(nextLocation model.Point) bool {
	return b.food.X == nextLocation.X && b.food.Y == nextLocation.Y
}

func (b *Board) EvaluateMove(dir model.Vector) (float32, bool) {
	nextLocation := b.NextLocation(dir)
	if b.MoveIsScoring(nextLocation) {
		return 100.0, false
	} else if b.MoveIsTerminal(nextLocation) {
		return -100.0, true
	} else {
		return -1, false
	}
}

func (b *Board) DangerAhead() bool {
	nextLocation := b.NextLocation(b.snake.direction)
	return b.MoveIsTerminal(nextLocation)
}

func (b *Board) DangerRight() bool {

	var nextLocation model.Point

	// If snake is going West, its right side is North
	if b.snake.direction == (model.Vector{X: 0, Y: -1}) {
		nextLocation = b.NextLocation(model.Vector{X: -1, Y: 0})
	}
	// If snake is going East, its right side is South
	if b.snake.direction == (model.Vector{X: 0, Y: 1}) {
		nextLocation = b.NextLocation(model.Vector{X: 1, Y: 0})
	}
	// If snake is going North, its right side is East
	if b.snake.direction == (model.Vector{X: -1, Y: 0}) {
		nextLocation = b.NextLocation(model.Vector{X: 0, Y: 1})
	}
	// If snake is going South, its right side is West
	if b.snake.direction == (model.Vector{X: 1, Y: 0}) {
		nextLocation = b.NextLocation(model.Vector{X: 0, Y: -1})
	}
	return b.MoveIsTerminal(nextLocation)
}

func (b *Board) DangerLeft() bool {
	var nextLocation model.Point
	// If snake is going West, its left side is South
	if b.snake.direction == (model.Vector{X: 0, Y: -1}) {
		nextLocation = b.NextLocation(model.Vector{X: 1, Y: 0})
	}
	// If snake is going East, its left side is North
	if b.snake.direction == (model.Vector{X: 0, Y: 1}) {
		nextLocation = b.NextLocation(model.Vector{X: -1, Y: 0})
	}
	// If snake is going North, its left side is West
	if b.snake.direction == (model.Vector{X: -1, Y: 0}) {
		nextLocation = b.NextLocation(model.Vector{X: 0, Y: -1})
	}
	// If snake is going South, its left side is East
	if b.snake.direction == (model.Vector{X: 1, Y: 0}) {
		nextLocation = b.NextLocation(model.Vector{X: 0, Y: 1})
	}

	return b.MoveIsTerminal(nextLocation)
}

/*
The state is an array of 11 values, representing:
  - Danger 1 OR 2 steps ahead
  - Danger 1 OR 2 steps on the right
  - Danger 1 OR 2 steps on the left
  - Snake is moving left
  - Snake is moving right
  - Snake is moving up
  - Snake is moving down
  - The food is on the left
  - The food is on the right
  - The food is on the upper side
  - The food is on the lower side
*/
func (b *Board) ReportState() []float32 {
	out := make([]float32, 11, 11)
	out[0] = boolToFloat32(b.DangerAhead())
	out[1] = boolToFloat32(b.DangerRight())
	out[2] = boolToFloat32(b.DangerLeft())
	out[3] = boolToFloat32(b.snake.direction == model.Vector{X: 0, Y: -1})
	out[4] = boolToFloat32(b.snake.direction == model.Vector{X: 0, Y: 1})
	out[5] = boolToFloat32(b.snake.direction == model.Vector{X: -1, Y: 0})
	out[6] = boolToFloat32(b.snake.direction == model.Vector{X: 1, Y: 0})
	out[7] = boolToFloat32(b.food.X < b.snake.Head().X)
	out[8] = boolToFloat32(b.food.X > b.snake.Head().X)
	out[9] = boolToFloat32(b.food.Y < b.snake.Head().Y)
	out[10] = boolToFloat32(b.food.Y > b.snake.Head().Y)

	return out
}

func boolToFloat32(b bool) float32 {
	if b {
		return 1.0
	}
	return 0.0
}
