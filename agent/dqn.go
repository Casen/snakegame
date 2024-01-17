package agent

import (
	"log"
	"math/rand"
	"time"

	. "github.com/casen/snakegame/model"
	"github.com/casen/snakegame/snake"
	"gorgonia.org/gorgonia"
)

var cardinals = [4]Vector{
	{X: 0, Y: 1},  // E
	{X: -1, Y: 0}, // N
	{X: 1, Y: 0},  // S
	{X: 0, Y: -1}, // W
}

type DQN struct {
	game *snake.Game
	NN   *Brain
	gorgonia.VM
	gorgonia.Solver
	Memories []Memory // The Q-Table - stores State/Action/Reward/NextState/NextMoves/IsDone - added to each train x times per episode

	gamma       float32
	epsilon     float32
	epsDecayMin float32
	decay       float32
	isTraining  bool
}

func (agent *DQN) init() {
	// Construct the NN Graph to prepare for matrix multiplication and backpropagation
	if _, err := agent.NN.cons(); err != nil {
		panic(err)
	}

	// Construct the VM and Solver, which will compute all the values in the NN Graph
	agent.VM = gorgonia.NewTapeMachine(agent.NN.g)
	agent.Solver = gorgonia.NewRMSPropSolver()
	agent.isTraining = false
}

func (agent *DQN) PredictQValue(gameState [11]float32) (float32, error) {
	agent.NN.Let1(gameState)
	if err := agent.VM.RunAll(); err != nil {
		log.Printf("Got an error on VM Run %v", err)
		return 0, err
	}
	agent.VM.Reset()
	retVal := agent.NN.predVal.Data().([]float32)[0]
	return retVal, nil
}

func (agent *DQN) BestMove() Vector {
	var action Vector
	moves := getPossibleActions(agent.game)

	if len(moves) < 1 && !agent.game.GameOver() {
		panic("No possible moves")
	}

	if len(moves) > 0 {
		action = agent.BestAction(moves)
	} else {
		log.Print("We reached a terminal state")
		log.Printf("Defaulting to move in current direction")
		action = agent.game.CurrentDirection()
	}

	return action
}

func (agent *DQN) Train() (err error) {
	agent.isTraining = true
	var episodes = 100
	var games = 50
	var score float32
	var gameCount int
	var totalMoves int
	var maxGameScore int = 0

	for e := 0; e < episodes; e++ {
		if e%100 == 0 && e > 99 {
			log.Printf("Episode %d, max game score %d", e, maxGameScore)
		}

		gameCount = 0
		totalMoves = 0
		for gameCount < games {

			if totalMoves > 10000 {
				if agent.game.Score() > maxGameScore {
					maxGameScore = agent.game.Score()
				}
				agent.game.Reset()
				gameCount++
				continue
			}

			state := agent.game.CurrentState()
			moves := getPossibleActions(agent.game)

			// No possible moves means the game is over now, or in the next step
			if len(moves) < 1 {
				log.Printf("No possible moves, game over: %t", agent.game.GameOver())
				if agent.game.Score() > maxGameScore {
					maxGameScore = agent.game.Score()
				}
				agent.game.Reset()
				gameCount++
				continue
			}

			// TODO use target network to predict Q values and train on separate network
			action := agent.BestAction(moves)

			reward, isDone := agent.game.EvaluateAction(action)
			score = score + reward

			agent.game.Move(action)
			totalMoves++

			if isDone {
				if agent.game.Score() > maxGameScore {
					maxGameScore = agent.game.Score()
				}
				agent.game.Reset()
				gameCount++
			}

			nextMoves := getPossibleActions(agent.game)
			futurePossibleStates := getPossibleStates(agent.game, nextMoves)
			mem := Memory{State: state, Action: action, Reward: reward, NextState: agent.game.CurrentState(), NextMovables: futurePossibleStates, isDone: isDone}
			agent.Memories = append(agent.Memories, mem)
		}

		if err := agent.Replay(32); err != nil {
			log.Printf("Got an error on replay %v", err)
			return err
		}

	}

	agent.isTraining = false
	agent.game.Reset()

	log.Printf("Training complete. Max game score %d", maxGameScore)

	return nil
}

func (agent *DQN) Replay(batchsize int) error {
	var totalMemories int = len(agent.Memories)
	var totalScoringMoves, totalTerminalMoves int = 0, 0
	var N int
	if batchsize < len(agent.Memories) {
		N = batchsize
	} else {
		N = len(agent.Memories)
	}
	mems := make([]Memory, N)
	r := rand.New(rand.NewSource(time.Now().Unix()))

	// Select N random memories from the Q-Table
	for i := range mems {
		mems[i] = agent.Memories[r.Intn(totalMemories-i)]
	}

	for b := 0; b < batchsize; b++ {
		mem := mems[b]
		if mem.Reward == 100 {
			totalScoringMoves++
		} else if mem.Reward == -100 {
			totalTerminalMoves++
		}

		var y float32
		if mem.isDone {
			y = mem.Reward
		} else {
			var nextRewards []float32
			for _, futureState := range mem.NextMovables {
				nextReward, err := agent.PredictQValue(futureState)
				if err != nil {
					return err
				}
				nextRewards = append(nextRewards, nextReward)
			}
			reward := max(nextRewards)
			y = mem.Reward + agent.gamma*reward
		}
		// Update the NN Graph and Set input state x and the max the target value y for the action we took.
		agent.NN.Let2(mem.State, y)

		// Run the NN Graph Calcs to get the predicted target value
		if err := agent.VM.RunAll(); err != nil {
			return err
		}
		agent.VM.Reset()
		if err := agent.Solver.Step(agent.NN.model()); err != nil {
			return err
		}
		if agent.epsilon > agent.epsDecayMin {
			agent.epsilon *= agent.decay
		}
	}

	return nil
}

func (agent *DQN) BestAction(moves []Vector) (bestAction Vector) {

	// If we're not training, strip use heuristic to avoid terminal actions
	if !agent.isTraining {
		nonTerminalMoves := agent.StripTerminalActions(moves)

		if len(nonTerminalMoves) > 0 {
			moves = nonTerminalMoves
		}
	}

	if len(moves) < 1 {
		panic("bestAction called with no moves")
	}

	var bestActions []Vector = make([]Vector, 0)
	var maxActValue float32 = -100

	for _, a := range moves {
		nextState := agent.game.NextState(a)

		// If we're not training, use heuristic to gaurantee scoring moves
		if !agent.isTraining {
			reward, _ := agent.game.EvaluateAction(a)
			if reward == 100 {
				return a
			}
		}

		actionValue, err := agent.PredictQValue(nextState)

		if err != nil {
			panic(err)
		}
		if actionValue > maxActValue {
			bestAction = a
			maxActValue = actionValue
			bestActions = append(bestActions, a)
		} else if actionValue == maxActValue {
			bestActions = append(bestActions, a)
		}
	}

	if rand.Float32() < agent.epsilon && len(bestActions) > 1 {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		randomAction := bestActions[r.Intn(len(bestActions))]
		return randomAction
	}

	return bestAction
}

func (agent *DQN) StripTerminalActions(actions []Vector) []Vector {
	var retVal []Vector

	for _, a := range actions {
		reward, _ := agent.game.EvaluateAction(a)
		if reward != -100 {
			retVal = append(retVal, a)
		}
	}

	return retVal
}

func getPossibleActions(g *snake.Game) (retVal []Vector) {
	if g.GameOver() {
		return retVal
	}

	for i := range cardinals {
		if g.MoveIsValid(cardinals[i]) {
			retVal = append(retVal, cardinals[i])
		}
	}

	return retVal
}

func getPossibleStates(g *snake.Game, moves []Vector) (retVal [][11]float32) {
	if g.GameOver() {
		return retVal
	}

	for _, m := range moves {
		retVal = append(retVal, g.NextState(m))
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

func truncateMemories(memories []Memory, num int) []Memory {
	max := len(memories) - num
	if max > 0 && len(memories) > max {
		return memories[:max-1]
	}
	return memories
}
