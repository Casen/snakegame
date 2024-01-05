package snake

import (
	"log"

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
	log.Printf("current direction %v", s.direction)
	// don't allow changing direction to opposite
	if oppositeDir := s.OppositeDir(newDir); !oppositeDir {
		log.Printf("changing direction to %v", newDir)
		s.direction = newDir
	}

	log.Printf("Whoops, can't go backwards")
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
	log.Printf("what %v", result)
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
	hitsSnake := s.HitsSnake(point)
	//The snake cannot reverse direction
	backwardsMove := s.Head().X-s.direction.X == point.X && s.Head().Y-s.direction.Y == point.Y

	return !hitsSnake && !backwardsMove
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
