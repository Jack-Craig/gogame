package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	inited bool
}

var (
	img *ebiten.Image
)

func (g *Game) init() {
	g.inited = true
	gdl := graphics.NewGraphicsDataLoader("res/play")
	img = gdl.GetSpriteImage(1)
}

func (g *Game) Update() error {
	if !g.inited {
		g.init()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.DrawImage(img, &ebiten.DrawImageOptions{})
	fmt.Println("Draw!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
