package main

import (
	"log"

	"github.com/casen/snakegame/agent"
	"github.com/casen/snakegame/snake"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	// Game defaults to user input
	game := snake.NewGame()
	ai := agent.NewAgent(game)

	ai.Train()
	game.Reset()

	log.Printf("Training complete")

	player := NewGamePlayer(game, ai, true)

	ebiten.SetWindowSize(snake.ScreenWidth, snake.ScreenHeight)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(player); err != nil {
		log.Fatal(err)
	}
}
