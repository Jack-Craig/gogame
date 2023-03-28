package gameplay

import (
	"fmt"
	"image/color"
	"log"
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
	MAXWORLDGENBUFFERLEN uint32 = 80
	// Min length of buffer to generate (trigger)
	MINWORLDGENBUFFERLEN uint32 = 40
	// Length of biome chunks
	BIOMELENGTH       uint32  = 8
	PLAYERWORLDSTARTX float32 = TILEWIDTH
	PLAYERWORLDSTARTY float32 = TILEWIDTH * float32(WORLDBUFFERHEIGHT-20)
	TOTALTILES        uint32  = 4
)

type World struct {
	Handler
	camera                                 *Camera
	gameObjects                            []*GameObject
	zombieObjects                          []*Zombie
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
	log.Printf("N Entities: %d\n", len(w.entityObjects))
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
	for i, gObj := range w.gameObjects {
		if gObj.shouldRemove {
			common.Remove(&w.gameObjects, i)
			continue
		}
	}

	for i, zombie := range w.zombieObjects {
		zombie.Update()
		if zombie.shouldRemove {
			common.Remove(&w.zombieObjects, i)
			continue
		}
	}

	for i, projectile := range w.projectiles {
		projectile.Update()
		if projectile.shouldRemove {
			common.Remove(&w.projectiles, i)
			continue
		}
	}
	for i, entity := range w.entityObjects {
		if entity.shouldRemove {
			common.Remove(&w.entityObjects, i)
			continue
		}
		entity.AddVel(0, w.gravity*entity.gravityMultiplier)
		entity.Update()
		entity.collidingEntities = nil
	}
	for i := 0; i < len(w.entityObjects); i++ {
		ei := w.entityObjects[i]
		for j := 0; j < len(w.entityObjects); j++ {
			if i == j {
				continue
			}
			ej := w.entityObjects[j]
			// Top left
			topLeft := ei.x >= ej.x && ei.x <= ej.x+ej.width && ei.y >= ej.y && ei.y <= ej.y+ej.height
			// Top right
			topRight := ei.x+ei.width >= ej.x && ei.x+ei.width <= ej.x+ej.width && ei.y >= ej.y && ei.y <= ej.y+ej.height
			// Bottom left
			bottomLeft := ei.x > ej.x && ei.x <= ej.x+ej.width && ei.y+ei.height >= ej.y && ei.y+ei.height <= ej.y+ej.height
			// Bottom right
			bottomRight := ei.x+ei.width >= ej.x && ei.x+ei.width <= ej.x+ej.width && ei.y+ei.height >= ej.y && ei.y+ei.height <= ej.y+ej.height
			isCollision := topLeft || topRight || bottomLeft || bottomRight
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

	if w.level.toBufferIndex(w.level.worldXStart) > w.level.toBufferIndex(w.level.worldXEnd+1) {
		for x := uint32(0); x < w.level.toBufferIndex(w.level.worldXStart); x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
		for x := w.level.toBufferIndex(w.level.worldXStart); x < WORLDBUFFERLEN; x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
	} else {
		for x := w.level.toBufferIndex(w.level.worldXStart); x <= w.level.toBufferIndex(w.level.worldXEnd+1); x++ {
			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				w.worldTiles[y][x].Draw(screen)
			}
		}
	}

	for _, gobj := range w.gameObjects {
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
	// Width of world in array coordinates
	worldWidth uint32
	// In array coordinates, start and end of visible world. Does not wrap
	worldXStart uint32
	worldXEnd   uint32
	// In array coordinates, the most recenty generated tile. Does not wrap
	worldXGen uint32
	// Ring array for biomes
	biomes      []Biome
	curBiomeIdx int
	biomeData   common.BiomeDataJson
}

func NewLevel(world *World, worldWidth uint32) *Level {
	l := Level{
		world:       world,
		worldWidth:  worldWidth,
		perlin:      perlin.NewPerlin(2, 2, 3, rand.Int63()),
		curBiomeIdx: 0,
		biomes:      make([]Biome, MAXWORLDGENBUFFERLEN/BIOMELENGTH+1),
	}
	common.LoadJSON("res/world/biomes.json", &l.biomeData)
	l.biomes[0].biomeType = "start"
	l.biomes[0].floorHeight = WORLDBUFFERHEIGHT / 2
	l.biomes[0].BiomeJson = l.biomeData.Biomes["start"]
	return &l
}

func (l *Level) initWorld() {
	for x := uint32(0); x < WORLDBUFFERLEN; x++ {
		for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
			l.world.worldTiles[y][x] = NewTile(0, float32(x*uint32(TILEWIDTH)), float32(y*uint32(TILEWIDTH)), l.world, nil)
		}
	}
}

func (l *Level) Update() {
	if l.worldXEnd >= l.worldWidth {
		l.world.canLeave = true
		return
	}
	offX, _ := l.world.camera.GetRenderOffset()
	l.worldXStart = uint32(-offX / TILEWIDTH)
	l.worldXEnd = l.worldXStart + uint32(l.world.camera.screenWidth/TILEWIDTH)

	//log.Printf("CurBiome: %s\n", l.biomes[int(l.worldXStart/BIOMELENGTH)%len(l.biomes)].biomeType)

	// Update world generation
	l.checkWorldUpdate()
}

func (l *Level) checkWorldUpdate() {
	// If we are MINWORLDGENBUFFERLEN away from the generated section, should generate until we are MAXWORLDGENBUFFERLEN past generated section
	if l.worldXStart+MINWORLDGENBUFFERLEN >= l.worldXGen {
		// Range of generation: l.worldXGen -> l.worldXStart + MAXWORLDGENBUFFERLEN. Typically of length MAXWORLDGENBUFFERLEN - MINWORLDGENBUFFERLEN
		for l.worldXStart+MAXWORLDGENBUFFERLEN >= l.worldXGen {
			// Check for current biome, and if we need to make a new one
			curBiome := &l.biomes[l.curBiomeIdx]
			floorBase := curBiome.floorHeight
			if curBiome.startX+BIOMELENGTH < l.worldXGen {
				// This biome has finished being generated!
				l.curBiomeIdx++
				l.curBiomeIdx %= len(l.biomes)
				// Generate new biome!
				randIdx := int(rand.Uint32()) % len(curBiome.BiomeJson.NextTo)
				newType := curBiome.BiomeJson.NextTo[randIdx]
				newCur := &l.biomes[l.curBiomeIdx]
				newCur.BiomeJson = l.biomeData.Biomes[newType]
				newCur.startX = l.worldXGen
				newCur.biomeType = newType
				curBiome = newCur
			}
			// Generate terrain
			arrX := l.toBufferIndex(l.worldXGen)
			generateAmplitude := curBiome.GenAmplitude

			rawY := l.perlin.Noise1D(float64(arrX) / (15.0 * curBiome.GenFrequency))
			groundY := floorBase + uint32(rawY*float64(generateAmplitude))
			curBiome.floorHeight = floorBase
			if curBiome.startX+BIOMELENGTH < l.worldXGen+1 {
				// This is the last square in the biome
				curBiome.floorHeight = groundY
			}

			// TOD bO: Make these actual tile objects or sm
			surfaceIm := l.world.gdl.GetSpriteImage(graphics.GrassTile)
			subsurfaceIm := l.world.gdl.GetSpriteImage(graphics.DirtTile)
			if curBiome.biomeType == "rocky" {
				surfaceIm = l.world.gdl.GetSpriteImage(graphics.RockTile)
				subsurfaceIm = l.world.gdl.GetSpriteImage(graphics.RockTile)
			}

			for y := uint32(0); y < WORLDBUFFERHEIGHT; y++ {
				tile := l.world.worldTiles[y][arrX]
				tile.x = float32(l.worldXGen) * TILEWIDTH
				if y == groundY {
					tile.im = surfaceIm
					tile.isPassable = false
				} else if y > groundY {
					tile.im = subsurfaceIm
					tile.isPassable = false
				} else {
					tile.im = nil
					tile.isPassable = true
				}
			}
			// Maybe zombie?
			if rand.Intn(10) < 1 {
				x := float32(l.worldXGen * uint32(TILEWIDTH))
				z := NewBaseZombie(x, float32(groundY)*TILEWIDTH-TILEWIDTH, l.world)
				l.world.zombieObjects = append(l.world.zombieObjects, z)
				l.world.AddEntity(&z.Entity)
			}
			l.worldXGen++
		}
	}
}

func (l *Level) toBufferIndex(x uint32) uint32 {
	return x % WORLDBUFFERLEN
}

// Variable length themes of tile chunks
type Biome struct {
	common.BiomeJson
	// string key in BiomeDataJson
	biomeType string
	// Where did this biome start, in array coords
	startX uint32
	// At the end of this biome, what is the floor height
	floorHeight uint32
}
