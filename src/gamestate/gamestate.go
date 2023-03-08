package gamestate

import (
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
}

func (ps *PlayState) GetNextState() GameState {
	return nil
}

func (ps *PlayState) Update() {

}

func (ps *PlayState) Draw(screen *ebiten.Image) {

}

// MENUSTATE
type MenuState struct {
	GameState
}

func (ms *MenuState) GetNextState() GameState {
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
