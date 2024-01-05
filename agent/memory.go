package agent

import (
	. "github.com/casen/snakegame/model"
)

type Memory struct {
	State        Point
	Action       Vector
	Reward       float32
	NextState    Point
	NextMovables []Vector
	isDone       bool
}
