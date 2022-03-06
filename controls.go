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
)

const (
	/* Actions' identifiers.
	   They have to be strings due to
	   testing these values with strings
	   from options_controls.cfg file. */
	StrMoveNorth = "MOVE_NORTH"
	StrMoveWest  = "MOVE_WEST"
	StrMoveEast  = "MOVE_EAST"
	StrMoveSouth = "MOVE_SOUTH"

	StrSetWater = "SET_WATER"
	StrSetFire = "SET_FIRE"
	StrSetEarth = "SET_EARTH"
	StrNextSchool1 = "NEXT_SCHOOL_1"
	StrNextSchool2 = "NEXT_SCHOOL_2"
	StrPrevSchool1 = "PREV_SCHOOL_1"
	StrPrevSchool2 = "PREV_SCHOOL_2"

	StrSetSmall = "SET_SMALL"
	StrSetMedium = "SET_MEDIUM"
	StrSetBig = "SET_BIG"
	StrSetHuge = "SET_HUGE"
	StrNextSize1 = "NEXT_SIZE_1"
	StrNextSize2 = "NEXT_SIZE_2"
	StrPrevSize1 = "PREV_SIZE_1"
	StrPrevSize2 = "PREV_SIZE_2"

	StrTargetSpell = "TARGET_SPELL"
	StrTargetGun = "TARGET_GUN"
	StrReload = "RELOAD"
	StrLook   = "LOOK"

	StrMouseLeftClick = "LEFT_CLICK"
	StrMouseRightClick = "RIGHT_CLICK"
	StrMouseScroll = "SCROLL"
	StrMouseWheel = "WHEEL"
)

var Actions = []string{
	// List of all possible actions.
	StrMoveNorth,
	StrMoveWest,
	StrMoveEast,
	StrMoveSouth,
	StrSetWater,
	StrSetFire,
	StrSetEarth,
	StrNextSchool1,
	StrNextSchool2,
	StrPrevSchool1,
	StrPrevSchool2,
	StrSetSmall,
	StrSetMedium,
	StrSetBig,
	StrSetHuge,
	StrNextSize1,
	StrNextSize2,
	StrPrevSize1,
	StrPrevSize2,
	StrTargetSpell,
	StrTargetGun,
	StrReload,
	StrLook,
	StrMouseLeftClick,
	StrMouseRightClick,
	StrMouseScroll,
	StrMouseWheel,
}

var CommandKeys = map[int]string{
	// Mapping keyboard scancodes to Action identifiers.
	blt.TK_UP:           StrMoveNorth,
	blt.TK_RIGHT:        StrMoveEast,
	blt.TK_DOWN:         StrMoveSouth,
	blt.TK_LEFT:         StrMoveWest,
	blt.TK_F1:           StrSetWater,
	blt.TK_F2:           StrSetFire,
	blt.TK_F3:           StrSetEarth,
	blt.TK_KP_MULTIPLY:  StrNextSchool1,
	blt.TK_RBRACKET:     StrNextSchool2,
	blt.TK_KP_DIVIDE:    StrPrevSchool1,
	blt.TK_LBRACKET:     StrPrevSchool2,
	blt.TK_1:            StrSetSmall,
	blt.TK_2:            StrSetMedium,
	blt.TK_3:            StrSetBig,
	blt.TK_4:            StrSetHuge,
	blt.TK_EQUALS:       StrNextSize1,
	blt.TK_KP_PLUS:      StrNextSize2,
	blt.TK_MINUS:        StrPrevSize1,
	blt.TK_KP_MINUS:     StrPrevSize2,
	blt.TK_F:            StrTargetSpell,
	blt.TK_T:            StrTargetGun,
	blt.TK_R:            StrReload,
	blt.TK_L:            StrLook,
	blt.TK_MOUSE_LEFT:   StrMouseLeftClick,
	blt.TK_MOUSE_RIGHT:  StrMouseRightClick,
	blt.TK_MOUSE_SCROLL: StrMouseScroll,
	blt.TK_MOUSE_WHEEL:  StrMouseWheel,
}

/* Place to store customized controls scheme,
   in the same manner as CommandKeys. */
var CustomCommandKeys = map[int]string{}

func Command(com string, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Command handles input received from Controls.
	   Most important argument passed to Command is string "com" that
	   is action identifier (action identifiers are stored as constants
	   at the top of this file). It calls player methods regarding to
	   passed command.
	   Returns true if command is valid and takes turn.
	   Otherwise, return false. */
	turnSpent := false
	switch com {
	case StrMouseLeftClick:
		turnSpent = ContextMouseLeft(p, *b, *o, *c)

	case StrMoveNorth:
		turnSpent = p.MoveOrAttack(0, -1, *b, o, *c)
	case StrMoveEast:
		turnSpent = p.MoveOrAttack(1, 0, *b, o, *c)
	case StrMoveSouth:
		turnSpent = p.MoveOrAttack(0, 1, *b, o, *c)
	case StrMoveWest:
		turnSpent = p.MoveOrAttack(-1, 0, *b, o, *c)

	case StrSetWater:
		GlobalData.CurrentSchool = SchoolWater
		AddMessage("You invoke water aura.")
		p.Colors = []string{"#73C2FB"}
	case StrSetFire:
		GlobalData.CurrentSchool = SchoolFire
		AddMessage("You invoke fire aura.")
		p.Colors = []string{"#FF7F7F"}
	case StrSetEarth:
		GlobalData.CurrentSchool = SchoolEarth
		AddMessage("You invoke earth aura.")
		p.Colors = []string{"#D2B48C"}
	case StrPrevSchool1, StrPrevSchool2:
		if GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSchool = SchoolEarth
			AddMessage("You invoke earth aura.")
			p.Colors = []string{"#D2B48C"}
		} else if GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSchool = SchoolWater
			AddMessage("You invoke water aura.")
			p.Colors = []string{"#73C2FB"}
		} else {
			GlobalData.CurrentSchool = SchoolFire
			AddMessage("You invoke fire aura.")
			p.Colors = []string{"#FF7F7F"}
		}
	case StrNextSchool1, StrNextSchool2:
		if GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSchool = SchoolFire
			AddMessage("You invoke fire aura.")
			p.Colors = []string{"#FF7F7F"}
		} else if GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSchool = SchoolEarth
			AddMessage("You invoke earth aura.")
			p.Colors = []string{"#D2B48C"}
		} else {
			GlobalData.CurrentSchool = SchoolWater
			AddMessage("You invoke water aura.")
			p.Colors = []string{"#73C2FB"}
		}
	case StrSetSmall:
		GlobalData.CurrentSize = SizeSmall
	case StrSetMedium:
		GlobalData.CurrentSize = SizeMedium
	case StrSetBig:
		GlobalData.CurrentSize = SizeBig
	case StrSetHuge:
		GlobalData.CurrentSize = SizeHuge
	case StrPrevSize1, StrPrevSize2:
		if GlobalData.CurrentSize == SizeSmall {
			GlobalData.CurrentSize = SizeHuge
		} else {
			GlobalData.CurrentSize--
		}
	case StrNextSize1, StrNextSize2:
		if GlobalData.CurrentSize == SizeHuge {
			GlobalData.CurrentSize = SizeSmall
		} else {
			GlobalData.CurrentSize++
		}

	case StrTargetSpell:
		turnSpent = p.Target(*b, o, *c, nil)
	case StrTargetGun:
		turnSpent = p.Target2(*b, o, *c, nil)
	case StrReload:
		if p.AmmoCurrent < p.AmmoMax {
			p.AmmoCurrent++
			turnSpent = true
		}
	case StrLook:
		p.Look(*b, *o, *c)
	}
	return turnSpent
}

func ContextMouseLeft(p *Creature, b Board, o Objects, c Creatures) bool {
	turnSpent := false
	x := blt.State(blt.TK_MOUSE_X)
	y := blt.State(blt.TK_MOUSE_Y)
	if x >= 0 && x < MapSizeX && y >= 0 && y < MapSizeY {
		// Clicks on map
		var monster *Creature
		for _, v := range c {
			if v.X == x && v.Y == y && v.HPCurrent > 0 {
				monster = v
			}
		}
		if monster != nil {
			turnSpent = p.Target(b, &o, c, monster)
		} else {
			p.Look(b, o, c)
		}
	} else {
		if x == UIPosX+1 && y == UIPosY+5 {
			if GlobalData.CurrentSchool == SchoolWater {
				turnSpent = p.Target(b, &o, c, nil)
			} else {
				GlobalData.CurrentSchool = SchoolWater
				AddMessage("You invoke water aura.")
				p.Colors = []string{"#73C2FB"}
			}
		} else if x == UIPosX+2 && y == UIPosY+5 {
			if GlobalData.CurrentSchool == SchoolFire {
				turnSpent = p.Target(b, &o, c, nil)
			} else {
				GlobalData.CurrentSchool = SchoolFire
				AddMessage("You invoke fire aura.")
				p.Colors = []string{"#FF7F7F"}
			}
		} else if x == UIPosX+3 && y == UIPosY+5 {
			if GlobalData.CurrentSchool == SchoolEarth {
				turnSpent = p.Target(b, &o, c, nil)
			} else {
				GlobalData.CurrentSchool = SchoolEarth
				AddMessage("You invoke earth aura.")
				p.Colors = []string{"#D2B48C"}
			}
		} else if x >= UIPosX && x <= UIPosX+4 && y == UIPosY+2 {
			turnSpent = p.Target2(b, &o, c, nil)
		} else if x >= UIPosX && x <= UIPosX+4 && (y == UIPosY+3 || y == UIPosY+4) {
			turnSpent = p.Target(b, &o, c, nil)
		} else if x == UIPosX+1 && y == UIPosY+6 && GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSize = SizeSmall
		} else if x == UIPosX+1 && y == UIPosY+7 && GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSize = SizeMedium
		} else if x == UIPosX+1 && y == UIPosY+8 && GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSize = SizeBig
		} else if x == UIPosX+1 && y == UIPosY+9 && GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSize = SizeHuge
		} else if x == UIPosX+2 && y == UIPosY+6 && GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSize = SizeSmall
		} else if x == UIPosX+2 && y == UIPosY+7 && GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSize = SizeMedium
		} else if x == UIPosX+2 && y == UIPosY+8 && GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSize = SizeBig
		} else if x == UIPosX+2 && y == UIPosY+9 && GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSize = SizeHuge
		} else if x == UIPosX+3 && y == UIPosY+6 && GlobalData.CurrentSchool == SchoolEarth {
			GlobalData.CurrentSize = SizeSmall
		} else if x == UIPosX+3 && y == UIPosY+7 && GlobalData.CurrentSchool == SchoolEarth {
			GlobalData.CurrentSize = SizeMedium
		} else if x == UIPosX+3 && y == UIPosY+8 && GlobalData.CurrentSchool == SchoolEarth {
			GlobalData.CurrentSize = SizeBig
		} else if x == UIPosX+3 && y == UIPosY+9 && GlobalData.CurrentSchool == SchoolEarth {
			GlobalData.CurrentSize = SizeHuge
		}
	}
	return turnSpent
}

func Controls(k int, p *Creature, b *Board, c *Creatures, o *Objects) bool {
	/* Function Controls takes integer 'k' (that is pressed key - blt uses
	   scancodes internally) and trying to find match key-command in
	   CommandKeys.
	   Value to return is determined in Command func. */
	turnSpent := false
	var command string
	if CustomControls == false {
		command = CommandKeys[k]
	} else {
		command = CustomCommandKeys[k]
	}
	turnSpent = Command(command, p, b, c, o)
	return turnSpent
}

func ReadInput() int {
	/* Function ReadInput is replacement of default blt's Read function that
	   returns QWERTY scancode. To provide (still experimental - I don't have
	   access to non-QWERTY keyboard physically) support for different
	   keyboard layouts, there are maps (in options.go) that matches
	   non-QWERTY input with QWERTY scancodes.
	   Some keys are hardcoded - like numpad, enter, etc. These hardcoded
	   keys are tested as first place as it's much cheaper operation than
	   checking map.
	   KeyMap content depends on chosen keyboard layout. */
	key := blt.Read()
	for _, v := range HardcodedKeys {
		if key == v {
			return v
		}
	}
	var r rune
	if blt.Check(blt.TK_WCHAR) != 0 {
		r = rune(blt.State(blt.TK_WCHAR))
	}
	return KeyMap[r]
}
