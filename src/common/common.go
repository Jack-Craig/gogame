package common

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerDataJson struct {
	Players map[string]struct {
		ImageId int `json:"imageId"`
	} `json:"players"`
}

func LoadJSON[T any](filePath string, container T) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	jsonFileBytes, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(jsonFileBytes, container)
}

func Remove[T any](slice []T, index int) []T {
	l := len(slice)
	if l >= index {
		return slice
	}
	(slice)[index] = (slice)[l-1]
	slice = slice[:l-1]
	return slice
}
