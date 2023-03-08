package graphics

import (
	"encoding/json"
	"image"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/Jack-Craig/gogame/src/common"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	TILESIZE = 32
)

// Loads buffered sprite sheet into memory
// Loads sprit sheet map into memory
// Serves requests from sprite_id to spritesheet coordinates
type GraphicsDataLoader struct {
	spriteSheet *ebiten.Image
	spriteMap   map[int]common.Pair
}

func NewGraphicsDataLoader(path string) *GraphicsDataLoader {
	gdl := &GraphicsDataLoader{}

	// Load spriteImage
	spriteImageFile, err := os.Open(path + "/sheet.png")
	if err != nil {
		log.Fatal(err)
	}
	defer spriteImageFile.Close()
	spriteImage, _, err := image.Decode(spriteImageFile)
	if err != nil {
		log.Fatal(err)
	}
	gdl.spriteSheet = ebiten.NewImageFromImage(spriteImage)

	// Load spriteMap
	spriteMapFile, err := os.Open(path + "/map.json")
	if err != nil {
		log.Fatal(err)
	}
	defer spriteMapFile.Close()
	spriteMapBytes, _ := ioutil.ReadAll(spriteMapFile)
	var mapData map[string]common.Pair
	json.Unmarshal(spriteMapBytes, &mapData)
	gdl.spriteMap = make(map[int]common.Pair)
	for spriteId_str, topLeft := range mapData {
		intVar, err := strconv.Atoi(spriteId_str)
		if err != nil {
			log.Fatal(err)
		}
		gdl.spriteMap[intVar] = topLeft
	}
	return gdl
}

func (gdl *GraphicsDataLoader) GetSpriteImage(spriteId int) *ebiten.Image {
	sheetLoc := gdl.spriteMap[spriteId]
	return gdl.spriteSheet.SubImage(image.Rect(
		sheetLoc.X*TILESIZE,
		sheetLoc.Y*TILESIZE,
		(sheetLoc.X+1)*TILESIZE,
		(sheetLoc.Y+1)*TILESIZE,
	)).(*ebiten.Image)
}
