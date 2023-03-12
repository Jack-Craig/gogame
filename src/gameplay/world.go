package gameplay

import (
	"log"
	"math/rand"

	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	// Square dimensions of tile in game coordinates
	TILEWIDTH float32 = 32
	// Maximum velocity of entities
	MAXVEL float32 = 12
	// Total length of world tile 2d array
	WORLDBUFFERLEN uint32 = 120
	// Total height of world tile 2d array
	WORLDBUFFERHEIGHT uint32 = 30
	// Length of world array to generate the world for
	WORLDGENBUFFERLEN uint32  = 60
	PLAYERWORLDSTARTX float32 = TILEWIDTH
	PLAYERWORLDSTARTY float32 = TILEWIDTH * float32(WORLDBUFFERHEIGHT-20)
	TOTALTILES        uint32  = 4
)

type World struct {
	gdl           *graphics.GraphicsDataLoader
	im            *input.InputManager
	camera        *Camera
	gameObjects   []*GameObject
	entityObjects []*Entity
	playerObjects []*Player
	gravity       float32
	worldTiles    [WORLDBUFFERHEIGHT][WORLDBUFFERLEN]*Tile
	inited        bool

	level *Level
}

func NewWorld(gdl *graphics.GraphicsDataLoader, im *input.InputManager) *World {
	w := &World{}
	w.gdl = gdl
	w.im = im
	w.camera = NewCamera(w)
	for id, playerInput := range *w.im.GetPlayerInputs() {
		log.Printf("Created player: ID: %d\n", id)
		p := NewPlayer(1, PLAYERWORLDSTARTX, PLAYERWORLDSTARTY, TILEWIDTH, TILEWIDTH, w, w.gdl.GetSpriteImage(4), playerInput)
		w.gameObjects = append(w.gameObjects, &p.GameObject)
		w.entityObjects = append(w.entityObjects, &p.Entity)
		w.playerObjects = append(w.playerObjects, p)
	}
	w.gravity = .25
	w.generateLevel()
	w.inited = true
	return w
}

func (w *World) generateLevel() {
	w.level = NewLevel(w)
	w.level.initWorld()
}

func (w *World) Update() {
	w.camera.Update()
	w.level.Update()
	for _, player := range w.playerObjects {
		player.Update()
	}
	for _, entity := range w.entityObjects {
		entity.AddVel(0, w.gravity)
		entity.Update()
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	if !w.inited {
		return
	}
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

// Given an x and y in world coordinates, returns true if there is a tile there and false otherwise
func (w *World) IsWorldCollision(x, y float32) bool {
	gridX, gridY := w.worldToBuffer(x, y)
	if int(gridY) >= len(w.worldTiles) {
		return false
	}
	if int(gridX) >= len(w.worldTiles[0]) {
		return false
	}
	tile := w.worldTiles[gridY][gridX]
	return !tile.isPassable
}

func (w *World) worldToBuffer(x, y float32) (uint32, uint32) {
	bufferY := uint32(y / TILEWIDTH)
	gridX := uint32(x / TILEWIDTH)
	bufferX := gridX % WORLDBUFFERLEN
	return bufferX, bufferY
}

// Levels, Biomes, and TileChunks handle world generation
type WorldDataLoader struct {
	tiles [TOTALTILES]*Tile
}

func NewWorldDataLoader() *WorldDataLoader {
	wdl := &WorldDataLoader{}
	return wdl
}

func (wdl *WorldDataLoader) GetTile(id uint32) *Tile {
	return wdl.tiles[id]
}

// Starts with entrance, ends with exit. Collection of biomes
type Level struct {
	world  *World
	perlin *perlin.Perlin
	//curBiome *Biome
	// In array coordinates, the start and end of the visible world. Can wrap
	worldFrameStart uint32
	worldFrameEnd   uint32
	// In world coordinates, the left most coordinate
	worldXStart float32
	// In array coordinates, the end of generated space
	worldGeneratedEnd uint32
}

func NewLevel(world *World) *Level {
	/**
	worldDataFile, err := os.Open("res/world/gen.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer worldDataFile.Close()
	worldDataBytes, _ := ioutil.ReadAll(worldDataFile)
	var worldDataJSON interface{}
	json.Unmarshal(worldDataBytes, &worldDataJSON)
	*/
	return &Level{
		world:  world,
		perlin: perlin.NewPerlin(2, 2, 3, rand.Int63()),
	}
}

func (l *Level) initWorld() {
	for x := uint32(0); x < WORLDBUFFERLEN; x++ {
		rawY := l.perlin.Noise1D(float64(x) / 10)
		groundY := (WORLDBUFFERHEIGHT - uint32(rawY*10)) - 5
		log.Printf("GroundY: %f, %d, %d\n", rawY, groundY, x)
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			if x == 0 {
				t := NewTile(10, float32(x)*TILEWIDTH, float32(y)*TILEWIDTH, l.world, l.world.gdl.GetSpriteImage(6))
				t.isPassable = true
				l.world.worldTiles[y][x] = t
				continue
			}
			if y == groundY {
				t := NewTile(10, float32(x)*TILEWIDTH, float32(y)*TILEWIDTH, l.world, l.world.gdl.GetSpriteImage(3))
				t.isPassable = false
				l.world.worldTiles[y][x] = t
			} else if y > groundY {
				t := NewTile(10, float32(x)*TILEWIDTH, float32(y)*TILEWIDTH, l.world, l.world.gdl.GetSpriteImage(1))
				t.isPassable = false
				l.world.worldTiles[y][x] = t
			} else {
				t := NewTile(10, float32(x)*TILEWIDTH, float32(y)*TILEWIDTH, l.world, l.world.gdl.GetSpriteImage(6))
				t.isPassable = true
				l.world.worldTiles[y][x] = t
			}
		}
		l.worldGeneratedEnd++
		l.worldGeneratedEnd %= WORLDBUFFERLEN
	}
}

func (l *Level) Update() {
	// Update ring offets
	if l.world.camera.screenWidth > 0 && l.worldFrameEnd == 0 {
		l.worldFrameEnd = l.worldFrameStart + uint32(l.world.camera.screenWidth/TILEWIDTH)
	}
	if l.world.camera.screenWidth > 0 && !l.world.camera.IsInsideCamera(l.worldXStart+TILEWIDTH, -1) {
		// Need to advance world buffer
		l.worldXStart += TILEWIDTH
		lastTile := l.worldFrameStart
		l.worldFrameStart += 1
		l.worldFrameStart %= WORLDBUFFERLEN
		l.worldFrameEnd = l.worldFrameStart + uint32(l.world.camera.screenWidth/TILEWIDTH)
		l.worldFrameEnd %= WORLDBUFFERLEN
		// Need to change coordinates of blocks that just left the screen
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			var prevTile uint32
			if lastTile == 0 {
				prevTile = WORLDBUFFERLEN - 1
			} else {
				prevTile = lastTile - 1
			}
			l.world.worldTiles[y][lastTile].x = l.world.worldTiles[y][prevTile].x + TILEWIDTH
		}
	}
	l.checkWorldUpdate()
}

func (l *Level) checkWorldUpdate() {
	// Check if we should generate more shit

	for (l.worldFrameStart+WORLDGENBUFFERLEN)%WORLDBUFFERLEN != l.worldGeneratedEnd%WORLDBUFFERLEN {
		arrX := l.worldGeneratedEnd
		worldX := uint32(l.world.worldTiles[0][l.worldGeneratedEnd].x)
		rawY := l.perlin.Noise1D(float64(worldX) / float64(TILEWIDTH*10))
		groundY := (WORLDBUFFERHEIGHT - uint32(rawY*10)) - 5
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			if y == groundY {
				l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(3)
				l.world.worldTiles[y][arrX].isPassable = false
			} else if y > groundY {
				l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(1)
				l.world.worldTiles[y][arrX].isPassable = false
			} else {
				l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(6)
				l.world.worldTiles[y][arrX].isPassable = true
			}
		}
		l.worldGeneratedEnd++
		l.worldGeneratedEnd %= WORLDBUFFERLEN
	}
}

/**
// Variable length themes of tile chunks
type Biome struct {
	// Length of biome in tiles
	tileLength uint32
	// Amount of tiles traversed in biome
	curTraversed uint32
}

func (w *World) generateWorld() {
	grass := w.gdl.GetSpriteImage(2)
	sky := w.gdl.GetSpriteImage(4)
	dirt := w.gdl.GetSpriteImage(1)
	for y := float32(0); y < float32(WORLDBUFFERHEIGHT); y++ {
		var row []*Tile
		for x := float32(0); x < float32(WORLDBUFFERLEN); x++ {
			if y == float32(WORLDBUFFERHEIGHT)-6 {
				t := NewTile(2, x*TILEWIDTH, y*TILEWIDTH, w, grass)
				row = append(row, t)
			} else if y > float32(WORLDBUFFERHEIGHT)-6 {
				t := NewTile(2, x*TILEWIDTH, y*TILEWIDTH, w, dirt)
				row = append(row, t)
			} else {
				t := NewTile(3, x*TILEWIDTH, y*TILEWIDTH, w, sky)
				t.isPassable = true
				row = append(row, t)
			}
		}
		w.worldTiles = append(w.worldTiles, row)
	}
}

// Templates to stamp during world generation
type TileChunk struct {
}

*/
