package agent

import (
	"github.com/casen/snakegame/snake"
)

type Agent struct {
	dqn  *DQN
	game *snake.Game
}

func NewAgent() *Agent {

	// var times int = 1000
	var gamma float32 = 0.95  // discount factor
	var epsilon float32 = 1.0 // exploration/exploitation bias, set to 1.0/exploration by default
	var epsilonDecayMin float32 = 0.01
	var epsilonDecay float32 = 0.995

	dqn := &DQN{
		NN:          NewNN(32),
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

func (a *Agent) Train(game *snake.Game) {
	a.dqn.train(game)
}
