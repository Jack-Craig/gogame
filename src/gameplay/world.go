package gameplay

import (
	"log"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TILEWIDTH float32 = 32
	MAXVEL    float32 = 12
)

type World struct {
	gdl           *graphics.GraphicsDataLoader
	im            *input.InputManager
	camera        *Camera
	gameObjects   []*GameObject
	entityObjects []*Entity
	playerObjects []*Player
	gravity       float32
	worldTiles    [][]*Tile
}

func NewWorld(gdl *graphics.GraphicsDataLoader, im *input.InputManager) *World {
	w := &World{}
	w.gdl = gdl
	w.im = im
	w.camera = NewCamera(w)
	for id, playerInput := range *w.im.GetPlayerInputs() {
		log.Printf("Created player: ID: %d\n", id)
		p := NewPlayer(1, 50, 1, TILEWIDTH, TILEWIDTH, w, w.gdl.GetSpriteImage(0), playerInput)
		w.gameObjects = append(w.gameObjects, &p.GameObject)
		w.entityObjects = append(w.entityObjects, &p.Entity)
		w.playerObjects = append(w.playerObjects, p)
	}
	w.gravity = .25
	w.generateWorld()
	return w
}

func (w *World) Update() {
	w.camera.Update()
	for _, player := range w.playerObjects {
		player.Update()
	}
	for _, entity := range w.entityObjects {
		entity.AddVel(0, w.gravity)
		entity.Update()
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	screenBounds := screen.Bounds().Max
	w.camera.screenWidth = float32(screenBounds.X)
	w.camera.screenHeight = float32(screenBounds.Y)

	for _, row := range w.worldTiles {
		for _, tile := range row {
			tile.GameObject.Draw(screen)
		}
	}
	for _, gobj := range w.gameObjects {
		gobj.Draw(screen)
	}
	w.camera.Draw(screen)
}

func (w *World) generateWorld() {
	for y := float32(0); y < 30; y++ {
		var row []*Tile
		for x := float32(0); x < 90; x++ {
			if y == 29 || y > 28 && uint32(x)&5 == 0 || y > 25 && uint32(x)&8 == 0 {
				im := w.gdl.GetSpriteImage(2)
				t := NewTile(2, x*TILEWIDTH, y*TILEWIDTH, w, im)
				row = append(row, t)
			} else {
				im := w.gdl.GetSpriteImage(3)
				t := NewTile(3, x*TILEWIDTH, y*TILEWIDTH, w, im)
				t.isPassable = true
				row = append(row, t)
			}
		}
		w.worldTiles = append(w.worldTiles, row)
	}
	log.Println("Generated World")
}

// Given an x and y in world coordinates, returns true if there is a tile there and false otherwise
func (w *World) IsWorldCollision(x, y float32) bool {
	gridX := int(x / TILEWIDTH)
	gridY := int(y / TILEWIDTH)
	if gridY < 0 || gridY >= len(w.worldTiles) {
		return false
	}
	if gridX < 0 || gridX >= len(w.worldTiles[0]) {
		return false
	}
	tile := w.worldTiles[gridY][gridX]
	return !tile.isPassable
}
