package gamestate

import (
	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type GameState interface {
	GetNextState() GameState
	Update()
	Draw(screen *ebiten.Image)
}

// PLAYSTATE
type PlayState struct {
	GameState
	inited bool
	gdl    *graphics.GraphicsDataLoader
	// World (Blocks)
	// Entities
	//	- Players
	//	- Enemies
	// Inputs
}

func (ps *PlayState) GetNextState() GameState {
	return nil
}

func (ps *PlayState) Update() {
	if !ps.inited {
		ps.init()
	}
}

func (ps *PlayState) init() {
	ps.gdl = graphics.NewGraphicsDataLoader("res/play")
	ps.inited = true
}

func (ps *PlayState) Draw(screen *ebiten.Image) {

}

// MENUSTATE
type MenuState struct {
	GameState
}

func (ms *MenuState) GetNextState() GameState {
	// TODO: Make Menu. For now, instantly play the game
	return &PlayState{}
}

func (ms *MenuState) Update() {

}

func (ms *MenuState) Draw(screen *ebiten.Image) {

}

// SHOPSTATE
type ShopState struct {
	GameState
}

func (ss *ShopState) GetNextState() GameState {
	return nil
}

func (ss *ShopState) Update() {

}

func (ss *ShopState) Draw(screen *ebiten.Image) {

}
