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

const(
	FireDurationMin = 16
	FireDurationMax = 20
	FireDurationNotFlammableDiv = 3

	BarrenDurationMin = 40
	BarrenDurationMax = 60
)

var (
	FireChars = []string{"^", "^", "^", "^", "^", "^", "^", "^", "^", "^", "^", "'", "'", "'", "'", "'", ".", ".", ".", "."}
	FireNotFlammableChars = []string{"^", "^", "'", "'", ".", "."}
	FireColors = []string{"red", "red", "red", "light red", "light red", "dark red", "dark red", "lighter red", "lighter red", "darker red", "darker red"}
)

func FireArea(area [][]int, b Board) {
	for _, v := range area {
		x := v[0]
		y := v[1]
		t := b[x][y]
		if t.Barren == 0 {
			t.FireTile()
		}
	}
}

func (t *Tile) FireTile() {
	t.MakeFire()
}

func (t *Tile) UpdateTile() {
	if t.Fire > 0 {
		t.Fire--
		if t.Fire == 0 {
			t.MakeBarren()
		}
	} else if t.Barren > 0 {
		t.Barren--
		if t.Barren == 0 {
			switch t.Name {
			case "stone ground":
				t.MakeStoneGround()
			case "grass":
				t.MakeGrass()
			}
		}
	}
}
