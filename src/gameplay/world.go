package gameplay

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Jack-Craig/gogame/src/common"
	"github.com/Jack-Craig/gogame/src/graphics"
	"github.com/Jack-Craig/gogame/src/input"
	"github.com/aquilax/go-perlin"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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
	playerInfos   []*PlayerInfo
	projectiles   []*Projectile
	gravity       float32
	worldTiles    [WORLDBUFFERHEIGHT][WORLDBUFFERLEN]*Tile
	inited        bool
	bg            *Background

	level *Level
}

func NewWorld(gdl *graphics.GraphicsDataLoader, im *input.InputManager, spriteIds []graphics.SpriteID) *World {
	w := &World{}
	w.gdl = gdl
	w.im = im
	w.camera = NewCamera(w)
	for id, playerInput := range *w.im.GetPlayerInputs() {
		p := NewPlayer(id+1, PLAYERWORLDSTARTX, PLAYERWORLDSTARTY, TILEWIDTH, TILEWIDTH, w, w.gdl.GetSpriteImage(spriteIds[id]), playerInput)
		p.health = 100
		w.gameObjects = append(w.gameObjects, &p.GameObject)
		w.entityObjects = append(w.entityObjects, &p.Entity)
		w.playerObjects = append(w.playerObjects, p)
		playerInfo := &PlayerInfo{p}
		w.playerInfos = append(w.playerInfos, playerInfo)
	}
	w.bg = NewBackground(w)
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
		if player.shouldRemove {
			// TODO: Figure out to do when player is dead
		}
		player.Update()
	}
	for i, projectile := range w.projectiles {
		projectile.Update()
		if projectile.shouldRemove {
			common.Remove(w.projectiles, i)
			continue
		}
	}
	for i, entity := range w.entityObjects {
		if entity.shouldRemove {
			common.Remove(w.entityObjects, i)
			continue
		}
		entity.AddVel(0, w.gravity*entity.gravityMultiplier)
		entity.Update()
		entity.collidingEntities = nil
	}
	for i := 0; i < len(w.entityObjects); i++ {
		ei := w.entityObjects[i]
		for j := i + 1; j < len(w.entityObjects); j++ {
			ej := w.entityObjects[j]
			// Top left
			isCollision := ei.x >= ej.x && ei.x <= ej.x+ej.width && ei.y >= ej.y && ei.y <= ej.y+ej.height
			// Top right
			isCollision = isCollision || (ei.x+ei.width >= ej.x && ei.x+ei.width <= ej.x+ej.width && ei.y >= ej.y && ei.y <= ej.y+ej.height)
			// Bottom left
			isCollision = isCollision || (ei.x > ej.x && ei.x <= ej.x+ej.width && ei.y+ei.height >= ej.y && ei.y+ei.height <= ej.y+ej.height)
			// Bottom right
			isCollision = isCollision || (ei.x+ei.width >= ej.x && ei.x+ei.width <= ej.x+ej.width && ei.y+ei.height >= ej.y && ei.y+ei.height <= ej.y+ej.height)
			if isCollision {
				ei.collidingEntities = append(ei.collidingEntities, ej)
				ej.collidingEntities = append(ej.collidingEntities, ei)
			}
		}
	}
}
func (w *World) AddEntity(e *Entity) {
	e.w = w
	w.gameObjects = append(w.gameObjects, &e.GameObject)
	w.entityObjects = append(w.entityObjects, e)
}

func (w *World) AddProjectile(b *Projectile) {
	w.AddEntity(&b.Entity)
	w.projectiles = append(w.projectiles, b)
}

func (w *World) Draw(screen *ebiten.Image) {
	if !w.inited {
		return
	}
	screen.Fill(color.RGBA{135, 205, 235, 255})
	screenBounds := screen.Bounds().Max
	w.camera.screenWidth = float32(screenBounds.X)
	w.camera.screenHeight = float32(screenBounds.Y)
	w.bg.Draw(screen)
	if w.level.worldFrameStart > (w.level.worldFrameEnd+1)%WORLDBUFFERLEN {
		for x := uint32(0); x < w.level.worldFrameEnd; x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
		for x := w.level.worldFrameEnd; x < WORLDBUFFERLEN; x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
	} else {
		for x := w.level.worldFrameStart; x <= w.level.worldFrameEnd+1; x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
	}

	for i, gobj := range w.gameObjects {
		if gobj.shouldRemove {
			w.gameObjects = common.Remove(w.gameObjects, i)
			continue
		}
		gobj.Draw(screen)
	}
	for _, playerObj := range w.playerObjects {
		playerObj.Draw(screen)
	}
	for _, playerInfo := range w.playerInfos {
		playerInfo.Draw(screen)
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

type PlayerInfo struct {
	player *Player
}

func (pi *PlayerInfo) Draw(screen *ebiten.Image) {
	renderY := int(pi.player.w.camera.screenHeight)
	renderX := int(pi.player.w.camera.screenWidth * (float32(pi.player.id) / float32(len(pi.player.w.playerObjects)+1)))

	bo := pi.player.w.gdl.GetSpriteImage(graphics.PlayerInfo)
	w, h := bo.Size()
	scale := int(TILEWIDTH*4) / w
	realWidth := w * scale
	realHeight := h * scale
	boxOp := ebiten.DrawImageOptions{}
	boxOp.GeoM.Scale(float64(scale), float64(scale))
	boxOp.GeoM.Translate(float64(renderX-realWidth/2), float64(renderY-realHeight))
	screen.DrawImage(bo, &boxOp)

	statusText := fmt.Sprintf("%0.f", pi.player.health)
	boxSize := text.BoundString(*pi.player.w.gdl.GetFontNormal(), statusText)
	width := boxSize.Max.X - boxSize.Min.X
	height := boxSize.Max.Y - boxSize.Min.Y
	f := pi.player.w.gdl.GetFontNormal()
	text.Draw(screen, pi.player.name, *pi.player.w.gdl.GetFontSmall(), renderX-width/2, renderY-height/2-24, color.White)
	text.Draw(screen, statusText, *f, renderX-width/2, renderY-height/2, color.White)

	op := ebiten.DrawImageOptions{}
	guyScale := 1.2 * float64(height) / float64(graphics.TILESIZE)
	guySize := guyScale * float64(graphics.TILESIZE)
	op.GeoM.Scale(guyScale, guyScale)
	op.GeoM.Translate(float64(renderX)-guySize-float64(width)/2-10, float64(renderY-height)-guySize/2)
	if pi.player.isDead {
		screen.DrawImage(pi.player.w.gdl.GetSpriteImage(graphics.Skull), &op)
	} else {
		screen.DrawImage(pi.player.im, &op)
	}
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
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			l.world.worldTiles[y][x] = NewTile(0, float32(x*uint32(TILEWIDTH)), float32(y*uint32(TILEWIDTH)), l.world, nil)
		}
	}
}

func (l *Level) Update() {
	// Update ring offets
	if l.world.camera.screenWidth > 0 && l.worldFrameEnd == 0 {
		l.worldFrameEnd = l.worldFrameStart + uint32(l.world.camera.screenWidth/TILEWIDTH)
		l.worldFrameEnd %= WORLDBUFFERLEN
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
	generateAmplitude := uint32(8)
	floorBase := WORLDBUFFERHEIGHT - generateAmplitude
	for (l.worldFrameStart+WORLDGENBUFFERLEN)%WORLDBUFFERLEN != l.worldGeneratedEnd%WORLDBUFFERLEN {
		arrX := l.worldGeneratedEnd
		worldX := uint32(l.world.worldTiles[0][l.worldGeneratedEnd].x)
		rawY := l.perlin.Noise1D(float64(worldX) / float64(TILEWIDTH*15))
		groundY := floorBase + uint32(rawY*float64(generateAmplitude))
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			if y == groundY {
				l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(graphics.GrassTile)
				l.world.worldTiles[y][arrX].isPassable = false
			} else if y > groundY {
				l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(graphics.DirtTile)
				l.world.worldTiles[y][arrX].isPassable = false
			} else {
				l.world.worldTiles[y][arrX].im = nil
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
	grass := w.gdl.GetTileImage(2)
	sky := w.gdl.GetTileImage(4)
	dirt := w.gdl.GetTileImage(1)
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
