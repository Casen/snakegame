package agent

import (
	. "github.com/casen/snakegame/model"
)

type Memory struct {
	State        [11]float32
	Action       Vector
	Reward       float32
	NextState    [11]float32
	NextMovables [][11]float32
	isDone       bool
}
