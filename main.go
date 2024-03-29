package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/Jack-Craig/gogame/src/gameplay"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	inited       bool
	currentState gameplay.GameState
}

func (g *Game) init() {
	g.inited = true
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	g.currentState.Update()

	nextState := g.currentState.GetNextState()
	if nextState != nil {
		g.currentState = nextState
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	g.currentState.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	rand.Seed(time.Now().Unix())
	ebiten.SetFullscreen(false)
	ebiten.SetWindowSize(940, 720)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{currentState: gameplay.NewMenuState()}); err != nil {
		log.Fatal(err)
	}
}
