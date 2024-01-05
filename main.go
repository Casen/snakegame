package main

import (
	"github.com/casen/snakegame/agent"
	"github.com/casen/snakegame/snake"
)

func main() {
	ai := agent.NewAgent()
	game := snake.NewGame()

	ai.Train(game)
	/*

		ebiten.SetWindowSize(snake.ScreenWidth, snake.ScreenHeight)
		ebiten.SetWindowTitle("Snake")
		if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
		}
	*/
}
