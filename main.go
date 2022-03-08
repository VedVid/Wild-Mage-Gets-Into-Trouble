/*
Copyright (c) 2018, Tomasz "VedVid" Nowakowski
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package main

import (
	blt "bearlibterminal"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var KeyboardLayout int
var CustomControls bool

var GlobalData GameData

var MsgBuf = []string{}
var LastTarget *Creature

func main() {
	var cells = new(Board)
	var objs = new(Objects)
	var actors = new(Creatures)
	StartGame(cells, actors, objs)
	for {
		RenderAll(*cells, *objs, *actors)
		if (*actors)[0].HPCurrent <= 0 {
			DeleteSaves()
			blt.Read()
			break
		}
		if blt.HasInput() == true {
			key := ReadInput()
			if (key == blt.TK_S && blt.Check(blt.TK_SHIFT) != 0) ||
				key == blt.TK_CLOSE {
				err := SaveGame(*cells, *actors, *objs)
				if err != nil {
					fmt.Println(err)
				}
				break
			} else if key == blt.TK_Q && blt.Check(blt.TK_SHIFT) != 0 {
				DeleteSaves()
				break
			} else {
				turnSpent := Controls(key, (*actors)[0], cells, actors, objs)
				if turnSpent == true {
					GlobalData.TurnsSpent++
					CheckMagic(*cells, *actors, objs)
					CreaturesTakeTurn(*cells, *actors, *objs)
					SpawnMonsters(actors, *cells)
					UpdateFrames(*cells, *objs, *actors)
					UpdateTiles(*cells)
				}
			}
		}
	}
	blt.Close()
}

func NewGame(b *Board, c *Creatures, o *Objects) {
	/* Function NewGame initializes game state - creates player, monsters, and game map.
	   This implementation is generic-placeholder, for testing purposes. */
	GlobalData = GameData{0, 0, 0, 0}
	player, err := NewPlayer(1, 1)
	if err != nil {
		fmt.Println(err)
	}
	*c = Creatures{player}
	*b, _, err = LoadJsonMap("smallInn.json")
	if err != nil {
		fmt.Println(err)
	}
}

func StartGame(b *Board, c *Creatures, o *Objects) {
	/* Function StartGame determines if game save is present (and valid), then
	   loads data, or initializes new game.
	   Panics if some-but-not-all save files are missing. */
	_, errBoard := os.Stat(MapPathGob)
	_, errCreatures := os.Stat(CreaturesPathGob)
	_, errObjects := os.Stat(ObjectsPathGob)
	_, errGlobalData := os.Stat(GlobalDataPathGob)
	if errBoard == nil && errCreatures == nil && errObjects == nil && errGlobalData == nil {
		LoadGame(b, c, o)
	} else if errBoard != nil && errCreatures != nil && errObjects != nil && errGlobalData != nil {
		NewGame(b, c, o)
	} else {
		txt := CorruptedSaveError(errBoard, errCreatures, errObjects, errGlobalData)
		fmt.Println("Error: save files are corrupted: " + txt)
		panic(-1)
	}
}

func init() {
	fontSizeFlag := flag.Int("fs", FontSize, "font size")
	flag.Parse()
	ActualFontSize = *fontSizeFlag
	// Get seed for random number generator
	rand.Seed(time.Now().UTC().UnixNano())
	// sin and cos tables for FOV algorithm
	InitializeFOVTables()
	// E.g. binds the process to single thread
	InitializeBLT()
	InitializeKeyboardLayouts()
	ReadOptionsControls()
	ChooseKeyboardLayout()
}
