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
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	TILESIZE = 16
)

// Loads buffered sprite sheet into memory
// Loads sprit sheet map into memory
// Serves requests from sprite_id to spritesheet coordinates
type GraphicsDataLoader struct {
	spriteSheet *ebiten.Image
	spriteMap   map[uint32]common.Pair
	font        font.Face
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
	gdl.spriteMap = make(map[uint32]common.Pair)
	for spriteId_str, topLeft := range mapData {
		intVar, err := strconv.Atoi(spriteId_str)
		if err != nil {
			log.Fatal(err)
		}
		gdl.spriteMap[uint32(intVar)] = topLeft
	}
	// Load font
	/**
	fontFile, err := os.Open("res/font.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer fontFile.Close()
	fontBytes, _ := ioutil.ReadAll(spriteMapFile)
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
	*/
	parsedFont, _ := opentype.Parse(fonts.MPlus1pRegular_ttf)
	font, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	gdl.font = text.FaceWithLineHeight(font, 54)
	return gdl
}

func (gdl *GraphicsDataLoader) GetSpriteImage(spriteId uint32) *ebiten.Image {
	sheetLoc := gdl.spriteMap[spriteId]
	return gdl.spriteSheet.SubImage(image.Rect(
		sheetLoc.X*TILESIZE,
		sheetLoc.Y*TILESIZE,
		(sheetLoc.X+1)*TILESIZE,
		(sheetLoc.Y+1)*TILESIZE,
	)).(*ebiten.Image)
}

func (gdl *GraphicsDataLoader) GetFont() *font.Face {
	return &gdl.font
}
