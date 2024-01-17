package main

import (
	"log"

	"github.com/casen/snakegame/agent"
	"github.com/casen/snakegame/model"
	"github.com/casen/snakegame/snake"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 600
	ScreenHeight = 600
)

type GamePlayer struct {
	visited   []model.Point
	highScore int
	game      *snake.Game
	agent     *agent.Agent
	input     *snake.Input
	ai        bool
}

func NewGamePlayer(game *snake.Game, agent *agent.Agent, ai bool) *GamePlayer {
	if game == nil || agent == nil {
		return nil
	}

	return &GamePlayer{
		visited:   make([]model.Point, 0),
		highScore: 0,
		game:      game,
		agent:     agent,
		input:     snake.NewInput(),
		ai:        ai,
	}
}

func (gp *GamePlayer) HumanMove() error {
	_, userAction, ok := gp.input.Action()
	currDir := gp.game.CurrentDirection()
	finalAction := userAction

	if !ok {
		finalAction = currDir
	}

	return gp.game.Update(finalAction)
}

func (gp *GamePlayer) AiMove() error {

	if len(gp.visited) == 128 {
		startIdx, period, hasCycle := DetectCycles(gp.visited)

		if hasCycle {
			log.Printf("Found cycle at startIdx: %d, period: %d", startIdx, period)
			log.Printf("Game over. Score %v, High score %v. Resetting game", gp.game.Score(), gp.highScore)
			if gp.game.Score() > gp.highScore {
				gp.highScore = gp.game.Score()
			}
			gp.game.Reset()
		} else {
			log.Printf("No cycle found")
			log.Print(gp.visited)
		}
		clear(gp.visited)
	}

	if gp.game.GameOver() {
		if gp.game.Score() > gp.highScore {
			gp.highScore = gp.game.Score()
		}
		log.Printf("Game over. Score %v, High score %v. Resetting game", gp.game.Score(), gp.highScore)
		gp.game.Reset()
		clear(gp.visited)
	}

	agentAction := gp.agent.BestMove()

	// If we're not moving, we're not going to add the current location to the visited array
	if len(gp.visited) < 1 || gp.visited[len(gp.visited)-1] != gp.game.CurrentLocation() {
		gp.visited = append(gp.visited, gp.game.CurrentLocation())
	}

	return gp.game.Update(agentAction)
}

func (gp *GamePlayer) Update() error {
	if gp.ai {
		return gp.AiMove()
	} else {
		return gp.HumanMove()
	}
}

func (gp *GamePlayer) Draw(screen *ebiten.Image) {
	gp.game.Draw(screen)
}

func (gp *GamePlayer) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}
