/*
Copyright (c) 2022, Tomasz "VedVid" Nowakowski
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
	"math/rand"
)

func (t *Tile) MakeFire() {
	if t.Flammable == true && t.Barren == 0 && t.Fire == 0 {
		t.Fire = RandRange(FireDurationMin, FireDurationMax)
		t.Chars = []string{}
		t.Colors = []string{}
		for i := 0; i < t.Fire; i++ {
			t.Chars = append(t.Chars, FireChars[i])
			t.Colors = append(t.Colors, FireColors[rand.Intn(len(FireColors))])
			t.CurrentFrame = len(FireChars) - t.Fire
			t.Delay = 1
		}
	}
	if t.Flammable == false && t.Barren == 0 && t.Fire == 0 {
		t.Fire = RandRange(FireDurationMin, FireDurationMax) / FireDurationNotFlammableDiv
		t.Chars = []string{}
		t.Colors = []string{}
		for i := 0; i < t.Fire; i++ {
			t.Chars = append(t.Chars, FireNotFlammableChars[i])
			t.Colors = append(t.Colors, FireColors[rand.Intn(len(FireColors))])
			t.CurrentFrame = len(FireNotFlammableChars) - t.Fire
			t.Delay = 1
		}
	}
}

func TryFireAnotherTile(t *Tile, b Board) {
	for x := t.X-1; x <= t.X+1; x++ {
		for y := t.Y-1; y <= t.Y+1; y++ {
			if x < 0 || x >= MapSizeX || y < 0 || y >= MapSizeY {
				continue
			}
			chances := 2
			if t.Flammable == true {
				chances = 7
			}
			if rand.Intn(100) <= chances {
				b[x][y].MakeFire()
			}
		}
	}
}

func (t *Tile) MakeBarren() {
	t.Barren = RandRange(BarrenDurationMin, BarrenDurationMax)
	t.Chars = []string{"."}
	t.Colors = []string{"darker gray"}
	t.CurrentFrame = 0
	t.Delay = 0
}

func (t *Tile) MakeStoneGround() {
	t.Chars = []string{"."}
	t.Colors = []string{"gray"}
	t.CurrentFrame = 0
	t.Delay = 0
}

func (t *Tile) MakeGrass() {
	t.Chars = []string{","}
	t.Colors = []string{"dark green"}
	t.CurrentFrame = 0
	t.Delay = 0
}
