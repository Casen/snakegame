package agent

import (
	"log"
	"os"

	"github.com/casen/snakegame/model"
	"github.com/casen/snakegame/snake"
	. "gorgonia.org/gorgonia"
	"gorgonia.org/tensor"
)

type Agent struct {
	dqn  *DQN
	game *snake.Game
}

func NewAgent(game *snake.Game) *Agent {

	// var times int = 1000
	var gamma float32 = 0.95  // discount factor
	var epsilon float32 = 1.0 // exploration/exploitation bias, set to 1.0/exploration by default
	var epsilonDecayMin float32 = 0.01
	var epsilonDecay float32 = 0.995

	dqn := &DQN{
		game:        game,
		NN:          NewBrain(32),
		gamma:       gamma,
		epsilon:     epsilon,
		epsDecayMin: epsilonDecayMin,
		decay:       epsilonDecay,
	}
	dqn.init()

	return &Agent{
		dqn: dqn,
	}
}

func (a *Agent) Train() {
	a.dqn.Train()
}

func (a *Agent) BestMove() model.Vector {
	return a.dqn.BestMove()
}

func (a *Agent) Test() {
	g := NewGraph()
	xB := []float32{2, 4}
	xT := tensor.New(tensor.WithBacking(xB), tensor.WithShape(2))
	x := NewVector(g, tensor.Float32, WithName("X"), WithShape(2), WithValue(xT))

	yB := []float32{2, 4, 1, 2}
	yT := tensor.New(tensor.WithBacking(yB), tensor.WithShape(2, 2))
	y := NewMatrix(g, tensor.Float32, WithName("Y"), WithShape(2, 2), WithValue(yT))

	z, err := Mul(y, x)

	machine := NewTapeMachine(g)
	if machine.RunAll() != nil {
		panic(err)
	}

	os.WriteFile("simple_graph.dot", []byte(g.ToDot()), 0644)
	log.Printf("z: %v", z.Value())
}
