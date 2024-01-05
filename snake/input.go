package snake

import (
	"log"

	. "github.com/casen/snakegame/model"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct{}

func NewInput() *Input {
	return &Input{}
}

func (i *Input) Dir() (ebiten.Key, Vector, bool) {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		log.Printf("pressed up")
		return ebiten.KeyArrowUp, Vector{X: -1, Y: 0}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		log.Printf("pressed left")
		return ebiten.KeyArrowLeft, Vector{X: 0, Y: -1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		log.Printf("pressed left")
		return ebiten.KeyArrowRight, Vector{X: 0, Y: 1}, true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		log.Printf("pressed down")
		return ebiten.KeyArrowDown, Vector{X: 1, Y: 0}, true
	}

	return 0, Vector{X: 0, Y: 0}, false
}
