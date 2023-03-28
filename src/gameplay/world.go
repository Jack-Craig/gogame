package gameplay

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/Jack-Craig/gogame/src/common"
	"github.com/Jack-Craig/gogame/src/graphics"
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
	// Max length of buffer to generate (cap)
	MAXWORLDGENBUFFERLEN uint32 = 60
	// Min length of buffer to generate (trigger)
	MINWORLDGENBUFFERLEN uint32  = 30
	PLAYERWORLDSTARTX    float32 = TILEWIDTH
	PLAYERWORLDSTARTY    float32 = TILEWIDTH * float32(WORLDBUFFERHEIGHT-20)
	TOTALTILES           uint32  = 4
)

type World struct {
	Handler
	camera                                 *Camera
	gameObjects                            []*GameObject
	entityObjects                          []*Entity
	playerObjects                          []*Player
	projectiles                            []*Projectile
	gravity                                float32
	worldTiles                             [WORLDBUFFERHEIGHT][WORLDBUFFERLEN]*Tile
	inited, canLeave, allPlayersDoneOrDead bool
	bg                                     *Background

	level *Level
}

func NewWorld(handler Handler) *World {
	w := &World{Handler: handler}
	w.camera = NewCamera(w)
	for _, player := range handler.players {
		player.w = w
		player.x = PLAYERWORLDSTARTX
		player.y = PLAYERWORLDSTARTX
		player.shouldRemove = false
		w.gameObjects = append(w.gameObjects, &player.GameObject)
		w.entityObjects = append(w.entityObjects, &player.Entity)
		w.playerObjects = append(w.playerObjects, player)
	}
	w.bg = NewBackground(w)
	w.gravity = .25
	w.generateLevel()
	w.inited = true
	return w
}

func (w *World) generateLevel() {
	w.level = NewLevel(w, 100)
	w.level.initWorld()
}

func (w *World) Update() {
	w.camera.Update()
	w.level.Update()

	w.allPlayersDoneOrDead = true
	for _, player := range w.playerObjects {
		if !w.camera.IsInsideCamera(player.x, -1) {
			player.shouldRemove = true
		}
		if !player.shouldRemove {
			// Remove from entities, gameobjects, keep in players
			w.allPlayersDoneOrDead = false
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
		w.DrawPlayerInfo(playerObj, screen)
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

// Should probably be moved to player object? But is a UI elem so idk
func (w *World) DrawPlayerInfo(player *Player, screen *ebiten.Image) {
	renderY := int(w.camera.screenHeight)
	renderX := int(w.camera.screenWidth * (float32(player.id) / float32(len(w.playerObjects)+1)))

	// Render player info background thing
	bo := w.gdl.GetSpriteImage(graphics.PlayerInfo)
	width, height := bo.Size()
	scale := int(TILEWIDTH*4) / width
	realWidth := width * scale
	realHeight := height * scale
	boxOp := ebiten.DrawImageOptions{}
	boxOp.GeoM.Scale(float64(scale), float64(scale))
	boxOp.GeoM.Translate(float64(renderX-realWidth/2), float64(renderY-realHeight))
	screen.DrawImage(bo, &boxOp)

	// Render player name, health
	statusText := fmt.Sprintf("%0.f", player.health)
	boxSize := text.BoundString(*w.gdl.GetFontNormal(), statusText)
	textWidth := boxSize.Size().X
	textHeight := boxSize.Size().Y
	f := w.gdl.GetFontNormal()
	text.Draw(screen, player.name, *w.gdl.GetFontSmall(), renderX-textWidth/2, renderY-textHeight/2-24, color.White)
	text.Draw(screen, statusText, *f, renderX-textWidth/2, renderY-textHeight/2, color.White)

	// Render tiny player (or skull)
	op := ebiten.DrawImageOptions{}
	guyScale := 1.2 * float64(height) / float64(graphics.TILESIZE)
	guySize := guyScale * float64(graphics.TILESIZE)
	op.GeoM.Scale(guyScale, guyScale)
	op.GeoM.Translate(float64(renderX)-guySize-float64(textWidth)/2-10, float64(renderY-textHeight)-guySize/2)
	if player.isDead {
		screen.DrawImage(w.gdl.GetSpriteImage(graphics.Skull), &op)
	} else {
		screen.DrawImage(player.im, &op)
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
	// Width of world in array coordinates
	worldWidth uint32
	// In array coordinates, the start and end of the visible world. Wraps
	worldFrameStart uint32
	worldFrameEnd   uint32
	// In array coordinates, start and end of visible world. Does not wrap
	worldXStart uint32
	worldXEnd   uint32
	// In array coordinates, the most recenty generated tile. Does not wrap
	worldXGen uint32
}

func NewLevel(world *World, worldWidth uint32) *Level {
	return &Level{
		world:      world,
		worldWidth: worldWidth,
		perlin:     perlin.NewPerlin(2, 2, 3, rand.Int63()),
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
	// World buffer/generation stuff
	if l.worldXEnd >= l.worldWidth {
		l.world.canLeave = true
		return
	}
	if l.world.camera.screenWidth > 0 && !l.world.camera.IsInsideCamera(float32(l.worldXStart)*TILEWIDTH+TILEWIDTH, -1) {
		// Need to advance world buffer
		l.worldXStart += 1
		l.worldXEnd = l.worldXStart + uint32(l.world.camera.screenWidth/TILEWIDTH)
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
	// Update world generation
	l.checkWorldUpdate()
}

func (l *Level) checkWorldUpdate() {
	// Check if we should generate more shit
	generateAmplitude := uint32(8)
	floorBase := WORLDBUFFERHEIGHT - generateAmplitude
	// If we are MINWORLDGENBUFFERLEN away from the generated section, should generate until we are MAXWORLDGENBUFFERLEN past generated section
	if l.worldXStart+MINWORLDGENBUFFERLEN >= l.worldXGen {

		// Range of generation: l.worldXGen -> l.worldXStart + MAXWORLDGENBUFFERLEN. Typically of length MAXWORLDGENBUFFERLEN - MINWORLDGENBUFFERLEN
		for l.worldXStart+MAXWORLDGENBUFFERLEN >= l.worldXGen {
			// Generate terrain
			arrX := l.toBufferIndex(l.worldXGen)
			rawY := l.perlin.Noise1D(float64(arrX) / 15.0)
			groundY := floorBase + uint32(rawY*float64(generateAmplitude))

			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				if y == groundY {
					if int(l.worldXGen) == int(l.worldWidth) {
						l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(graphics.RockTile)
					} else {
						l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(graphics.GrassTile)
						l.world.worldTiles[y][arrX].isPassable = false
					}
				} else if y > groundY {
					l.world.worldTiles[y][arrX].im = l.world.gdl.GetSpriteImage(graphics.DirtTile)
					l.world.worldTiles[y][arrX].isPassable = false
				} else {
					l.world.worldTiles[y][arrX].im = nil
					l.world.worldTiles[y][arrX].isPassable = true
				}
			}
			l.worldXGen++
		}
	}
}

func (l *Level) toBufferIndex(x uint32) uint32 {
	return x % WORLDBUFFERLEN
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
