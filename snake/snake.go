package snake

import (
	"github.com/casen/snakegame/model"
)

type Coord struct {
	x, y int
}

type Snake struct {
	body      []model.Point
	direction model.Vector
	justAte   bool
}

func NewSnake(body []model.Point, direction model.Vector) *Snake {
	return &Snake{
		body:      body,
		direction: direction,
	}
}

func (s *Snake) Head() model.Point {
	return s.body[len(s.body)-1]
}

func (s *Snake) ChangeDirection(newDir model.Vector) {

	//If the new direction is the same as the current direction, do nothing
	if newDir.X == s.direction.X && newDir.Y == s.direction.Y {
		return
	}

	// don't allow changing direction to opposite
	if oppositeDir := s.OppositeDir(newDir); !oppositeDir {
		s.direction = newDir
	}
}

func (s *Snake) HeadHits(point model.Point) bool {
	h := s.Head()

	return h.X == point.X && h.Y == point.Y
}

func (s *Snake) HitsSnake(point model.Point) bool {
	for _, b := range s.body {
		if b.X == point.X && b.Y == point.Y {
			return true
		}
	}

	return false
}

// Checks if a point represents a move in the opposite direction the snake is currently moving
func (s *Snake) OppositeDir(dir model.Vector) bool {
	result := dir.X == -1*s.direction.X && dir.Y == -1*s.direction.Y
	return result
}

func (s *Snake) HeadHitsBody() bool {
	h := s.Head()
	bodyWithoutHead := s.body[:len(s.body)-1]

	for _, b := range bodyWithoutHead {
		if b.X == h.X && b.Y == h.Y {
			return true
		}
	}

	return false
}

func (s *Snake) ValidMove(point model.Point) bool {
	//The snake cannot hit itself
	//hitsSnake := s.HitsSnake(point)

	//Any backwards move would be equal to the second to last point in the snake's body
	//Let's call that point the snake's neck
	snakeNeck := s.body[len(s.body)-2]
	backwardsMove := snakeNeck.X == point.X && snakeNeck.Y == point.Y

	//isValid := !hitsSnake && !backwardsMove
	return !backwardsMove
}

func (s *Snake) Move() {
	h := s.Head()
	newHead := model.Point{X: h.X + s.direction.X, Y: h.Y + s.direction.Y}

	if s.justAte {
		s.body = append(s.body, newHead)
		s.justAte = false
	} else {
		s.body = append(s.body[1:], newHead)
	}
}

func (s *Snake) Clone() *Snake {
	bodyClone := make([]model.Point, len(s.body))
	copy(bodyClone, s.body)
	return NewSnake(bodyClone, s.direction)
}
