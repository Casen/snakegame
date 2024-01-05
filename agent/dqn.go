package agent

import (
	"log"
	"math/rand"

	. "github.com/casen/snakegame/model"
	"github.com/casen/snakegame/snake"
	"gorgonia.org/gorgonia"
)

var cardinals = [4]Vector{
	Vector{X: 0, Y: 1},  // E
	Vector{X: -1, Y: 0}, // N
	Vector{X: 1, Y: 0},  // S
	Vector{X: 0, Y: -1}, // W
}

type DQN struct {
	*NN
	gorgonia.VM
	gorgonia.Solver
	Memories []Memory // The Q-Table - stores State/Action/Reward/NextState/NextMoves/IsDone - added to each train x times per episode

	gamma       float32
	epsilon     float32
	epsDecayMin float32
	decay       float32
}

func (agent *DQN) init() {
	if _, err := agent.NN.cons(); err != nil {
		panic(err)
	}
	agent.VM = gorgonia.NewTapeMachine(agent.NN.g)
	agent.Solver = gorgonia.NewRMSPropSolver()
}

func (agent *DQN) replay(batchsize int) error {
	var N int
	if batchsize < len(agent.Memories) {
		N = batchsize
	} else {
		N = len(agent.Memories)
	}
	Xs := make([]input, 0, N)
	Ys := make([]float32, 0, N)
	mems := make([]Memory, N)
	copy(mems, agent.Memories)
	rand.Shuffle(len(mems), func(i, j int) {
		mems[i], mems[j] = mems[j], mems[i]
	})

	for b := 0; b < batchsize; b++ {
		mem := mems[b]

		var y float32
		if mem.isDone {
			y = mem.Reward
		} else {
			var nextRewards []float32
			for _, next := range mem.NextMovables {
				nextReward, err := agent.predict(mem.NextState, next)
				if err != nil {
					return err
				}
				nextRewards = append(nextRewards, nextReward)
			}
			reward := max(nextRewards)
			y = mem.Reward + agent.gamma*reward
		}
		Xs = append(Xs, input{mem.State, mem.Action})
		Ys = append(Ys, y)
		if err := agent.VM.RunAll(); err != nil {
			return err
		}
		agent.VM.Reset()
		if err := agent.Solver.Step(agent.model()); err != nil {
			return err
		}
		if agent.epsilon > agent.epsDecayMin {
			agent.epsilon *= agent.decay
		}
	}
	return nil
}

func (agent *DQN) predict(player Point, action Vector) (float32, error) {
	x := input{State: player, Action: action}
	agent.Let1(x)
	if err := agent.VM.RunAll(); err != nil {
		log.Printf("Got an error on VM Run %v", err)
		return 0, err
	}
	agent.VM.Reset()
	retVal := agent.predVal.Data().([]float32)[0]
	return retVal, nil
}

func (agent *DQN) train(game *snake.Game) (err error) {
	var episodes = 20000
	var times = 1000
	var score float32

	for e := 0; e < episodes; e++ {
		for t := 0; t < times; t++ {
			if e%100 == 0 && t%999 == 1 {
				log.Printf("episode %d, %dst loop", e, t)
			}

			log.Printf("game state: %v", game.ReportState())
			moves := getPossibleActions(game)
			action := agent.bestAction(game, moves)

			reward, isDone := game.EvaluateMove(action)
			score = score + reward
			player := game.PlayerPosition()

			game.Move(action)

			nextMoves := getPossibleActions(game)
			mem := Memory{State: player, Action: action, Reward: reward, NextState: game.PlayerPosition(), NextMovables: nextMoves, isDone: isDone}
			agent.Memories = append(agent.Memories, mem)
		}
	}
	return nil
}

func (agent *DQN) bestAction(game *snake.Game, moves []Vector) (bestAction Vector) {
	var bestActions []Vector
	var maxActValue float32 = -100
	for _, a := range moves {
		playerPosition := game.PlayerPosition()
		log.Printf("player position %v", playerPosition)
		actionValue, err := agent.predict(playerPosition, a)
		if err != nil {
			panic(err)
		}
		if actionValue > maxActValue {
			maxActValue = actionValue
			bestActions = append(bestActions, a)
		} else if actionValue == maxActValue {
			bestActions = append(bestActions, a)
		}
	}
	// shuffle bestActions
	rand.Shuffle(len(bestActions), func(i, j int) {
		bestActions[i], bestActions[j] = bestActions[j], bestActions[i]
	})

	log.Printf("Best Actions %v", bestActions)

	return bestActions[0]
}

func getPossibleActions(g *snake.Game) (retVal []Vector) {
	for i := range cardinals {
		if g.MoveIsValid(cardinals[i]) {
			retVal = append(retVal, cardinals[i])
		}
	}
	return retVal
}

func max(a []float32) float32 {
	var m float32 = -999999999
	for i := range a {
		if a[i] > m {
			m = a[i]
		}
	}
	return m
}
