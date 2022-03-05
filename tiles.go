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

const(
	FireDurationMin = 16
	FireDurationMax = 20
	FireNotFlammableDurationMin = 3
	FireNotFlammableDurationMax = 6

	BarrenDurationMin = 40
	BarrenDurationMax = 60

	FloodedDurationMin = 30
	FloodedDurationMax = 40

	DampDurationMin = 25
	DampDurationMax = 35

	ChasmDurationMin = 95
	ChasmDurationMax = 100

	UnstableDurationMin = 30
	UnstableDurationMax = 35
)

var (
	FireChars = []string{"^", "^", "^", "^", "^", "^", "^", "^", "^", "^", "^", "'", "'", "'", "'", "'", ".", ".", ".", "."}
	FireNotFlammableChars = []string{"^", "^", "'", "'", ".", "."}
	FireColors = []string{"red", "red", "red", "light red", "light red", "dark red", "dark red", "lighter red", "lighter red", "darker red", "darker red"}

	FloodedChars = []string{"≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈",
							"≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈", "≈",
							"≈", "≈", "≈", "~", "~", "~", "~", "~", "~", "~",
							"~", "~", "~", "~", "~", "~", "~", "~", "~", "~"}
	FloodedColors = []string{"blue", "blue", "blue", "blue", "blue", "dark blue", "dark blue", "dark blue", "dark blue", "dark blue",
							"blue", "blue", "blue", "blue", "blue", "lighter blue", "lighter blue", "lighter blue", "blue", "blue",
							"lighter blue", "blue", "blue", "darker blue", "light blue", "blue", "blue", "light blue", "light blue", "light blue",
							"blue", "blue", "blue", "blue", "dark blue", "light blue", "light blue", "lighter blue", "lighter blue", "blue"}
)

func (t *Tile) MakeFire() {
	if t.Barren == 0 && t.Fire == 0 && t.Flooded == 0 && t.Damp == 0 && t.Chasm == 0 {
		if t.Flammable == true {
			t.Fire = RandRange(FireDurationMin, FireDurationMax)
			t.Chars = FireChars
			t.CurrentFrame = len(FireChars) - t.Fire
			t.Delay = 1
			t.Colors = []string{}
			for i := 0; i < len(FireChars); i++ {
				t.Colors = append(t.Colors, FireColors[rand.Intn(len(FireColors))])
			}
		} else {
			t.Fire = RandRange(FireNotFlammableDurationMin, FireNotFlammableDurationMax)
			t.Chars = FireNotFlammableChars
			t.CurrentFrame = len(FireNotFlammableChars) - t.Fire
			t.Delay = 1
			t.Colors = []string{}
			for i := 0; i < len(FireNotFlammableChars); i++ {
				t.Colors = append(t.Colors, FireColors[rand.Intn(len(FireColors))])
			}
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
			if t.Damp > 0 || t.Flooded > 0 {
				chances = -1
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

func (t *Tile) MakeWater() {
	if t.Fire == 0 && t.Chasm == 0 {
		t.Flooded = RandRange(FloodedDurationMin, FloodedDurationMax)
		t.Chars = FloodedChars
		t.Colors = FloodedColors
		t.CurrentFrame = len(FloodedChars) - t.Flooded
		t.Delay = 1
	} else if t.Fire > 0 && t.Chasm == 0 {
		// Maybe we could add some STEAM here?
		t.Fire = 0
		t.Flooded = RandRange(FloodedDurationMin, FloodedDurationMax)
		t.Chars = FloodedChars
		t.Colors = FloodedColors
		t.CurrentFrame = len(FloodedChars) - t.Flooded
		t.Delay = 1
	}
	if t.Barren > 0 {
		t.Barren = 0
	}
}

func (t *Tile) MakeDamp() {
	t.Damp = RandRange(DampDurationMin, DampDurationMax)
	t.Chars = []string{"."}
	t.Colors = []string{"lighter blue"}
	t.CurrentFrame = 0
	t.Delay = 0
}

func (t *Tile) MakeChasm() {
	t.Chasm = RandRange(ChasmDurationMin, ChasmDurationMax)
	t.Fire = 0
	t.Barren = 0
	t.Flooded = 0
	t.Damp = 0
	t.Chars = []string{"…"}
	t.Colors = []string{"darkest blue"}
	t.CurrentFrame = 0
	t.Delay = 0
}

func (t *Tile) MakeUnstableGround() {
	t.Unstable = RandRange(UnstableDurationMin, UnstableDurationMax)
	t.Chars = []string{"."}
	t.Colors = []string{"darker orange"}
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
