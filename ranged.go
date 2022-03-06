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
	"errors"
	"fmt"
	"sort"
)

func (c *Creature) Look(b Board, o Objects, cs Creatures) {
	/* Look is method of Creature (that is supposed to be player).
	   It has to take Board, "global" Objects and Creatures as arguments,
	   because function PrintBrensenham need to call RenderAll function.
	   At first, Look creates new para-vector, with player coords as
	   starting point, and dynamic end position.
	   Then ComputeBrensenham checks what tiles are present
	   between Start and End, and adds their coords to Brensenham values.
	   Line from Brensenham is drawn, then game waits for player input,
	   that will change position of "looking" cursors.
	   Loop breaks with Escape, Space or Enter input. */
	startX, startY := c.X, c.Y
	targetX, targetY := startX, startY
	msg := ""
	i := false
	for {
		vec, err := NewBrensenham(startX, startY, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		_ = ComputeBrensenham(vec)
		_, _, _, _ = ValidateBrensenham(vec, b, cs, o)
		PrintBrensenham(vec, BrensenhamWhyInspect, BrensenhamColorNeutral, BrensenhamColorNeutral, b, o, cs)
		var monster *Creature
		var monsterDead *Creature
		var tile *Tile
		for _, vc := range cs {
			if vc.X == targetX && vc.Y == targetY {
				if vc.HPCurrent > 0 {
					monster = vc
					break
				} else {
					monsterDead = vc
				}
			}
		}
		tile = b[targetX][targetY]
		if monster != nil {
			msg = monster.Name
			if monster.FireResistance == FullAbility {
				msg = msg + " [[Fire res]]"
			} else if monster.FireResistance == NoAbility {
				msg = msg + " [[Fire vuln]]"
			}
			if monster.CanFly == FullAbility {
				msg = msg + " [[Flying]]"
			}
			if monster.CanSwim == FullAbility && monster.CanFly != FullAbility {
				msg = msg + " [[Can swim]]"
			}
		} else {
			if monsterDead != nil {
				msg = tile.Name + ", corpse."
			} else {
				msg = tile.Name
			}
		}
		PrintLookingMessage(msg, i)
		key := ReadInput()
		if key == blt.TK_ESCAPE || key == blt.TK_ENTER || key == blt.TK_SPACE {
			break
		}
		CursorMovement(&targetX, &targetY, key)
		i = true
	}
}

func PrintLookingMessage(s string, b bool) {
	/* Function PrintLookingMessage takes string (message) and bool ("is it
	   a first iteration?") as arguments.
	   It is used to provide dynamic printing looking message:
	   player do not need to confirm target to see what is it, but messages
	   will not flood message log. */
	l := len(MsgBuf)
	if s != "" {
		switch {
		case l == 0:
			AddMessage(s)
		case l >= MaxMessageBuffer:
			RemoveLastMessage()
			AddMessage(s)
		case l > 0 && l < MaxMessageBuffer:
			if b == true {
				RemoveLastMessage()
			}
			AddMessage(s)
		}
	}
}

func FormatLookingMessage(s []string, fov bool) string {
	/* FormatLookingMessage is function that takes slice of strings as argument
	   and returns string.
	   Player "see" things in his fov, and "recalls" out of his fov.
	   It is used to format Look() messages properly.
	   If slice is empty, it return empty tile message.
	   If slice contains only one item, it creates simplest message.
	   If slice is longer, it starts to format message - but it is
	   explicitly visible in function body. */
	_ = fov
	return s[0]
}

func (c *Creature) Target(b Board, o *Objects, cs Creatures, t *Creature) bool {
	/* Target is method of Creature, that takes game map, objects, and
	   creatures as arguments. Returns bool that serves as indicator if
	   action took some time or not.
	   This method is "the big one", general, for handling targeting.
	   In short, player starts targetting, line is drawn from player
	   to monster, then function waits for input (confirmation - "fire",
	   breaking the loop, or continuing).
	   Explicitly:
	   - creates list of all potential targets in fov
	    * tries to automatically last target, but
	    * if fails, it targets the nearest enemy
	   - draws line between source (receiver) and target (coords)
	    * creates new vector
	    * checks if it is valid - monsterHit should not be nil
	    * prints brensenham's line (ie so-called "vector")
	   - waits for player input
	    * if player cancels, function ends
	    * if player confirms, valley is shoot (in target, or empty space)
	    * if valley is shot in empty space, vector is extrapolated to check
	      if it will hit any target
	    * player can switch between targets as well; it targets
	      next target automatically; at first, only monsters that are
	      valid target (ie clean shot is possible), then monsters that
	      are in range and fov, but line of shot is not clear
	    * in other cases, game will try to move cursor; invalid input
	      is ignored */
	turnSpent := false
	if c.ManaCurrent <= 0 {
		AddMessage("No mana!")
		return turnSpent
	}
	var target *Creature
	targets := c.FindTargets(FOVLength, b, cs, *o)
	if LastTarget != nil && LastTarget != c &&
		IsInFOV(b, c.X, c.Y, LastTarget.X, LastTarget.Y) == true {
		target = LastTarget
	} else {
		var err error
		target, err = c.FindTarget(targets)
		if err != nil {
			fmt.Println(err)
		}
	}
	if t != nil {
		target = t
	}
	targetX, targetY := target.X, target.Y
	i := false
	for {
		cursor := SetCursorSize()
		area := [][]int{}
		for _, v := range cursor {
			ax := targetX + v[0]
			ay := targetY + v[1]
			if ax >= 0 && ay >= 0 && ax < MapSizeX && ay < MapSizeY {
				area = append(area, []int{targetX+v[0], targetY+v[1]})
			}
		}
		PrintCursor(area, b, *o, cs)
		vec, err := NewBrensenham(c.X, c.Y, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		_ = ComputeBrensenham(vec)
		_, _, monsterHit, _ := ValidateBrensenham(vec, b, targets, *o)
		if monsterHit != nil {
			msg := monsterHit.Name
			if monsterHit.FireResistance == FullAbility {
				msg = msg + " [[Fire res]]"
			} else if monsterHit.FireResistance == NoAbility {
				msg = msg + " [[Fire vuln]]"
			}
			if monsterHit.CanFly == FullAbility {
				msg = msg + " [[Flying]]"
			}
			if monsterHit.CanSwim == FullAbility && monsterHit.CanFly != FullAbility {
				msg = msg + " [[Can swim]]"
			}
			PrintLookingMessage(msg, i)
		}
		key := ReadInput()
		if key == blt.TK_ESCAPE {
			break
		}
		if AdjustSpell(key, c) == true {
			continue
		}
		if key == blt.TK_F || key == blt.TK_ENTER || key == blt.TK_MOUSE_LEFT {
			turnSpent = true
			c.ManaCurrent--
			if GlobalData.CurrentSchool == SchoolFire {
				FireArea(area, b)
			} else if GlobalData.CurrentSchool == SchoolWater {
				WaterArea(area, b)
			} else if GlobalData.CurrentSchool == SchoolEarth {
				RemoveArea(area, b)
			}
			monsterAimed := FindMonsterByXY(targetX, targetY, cs)
			if monsterAimed != nil && monsterAimed != c && monsterAimed.HPCurrent > 0 {
				LastTarget = monsterAimed
			}
			break
		} else if key == blt.TK_TAB {
			i = true
			monster := FindMonsterByXY(targetX, targetY, cs)
			if monster != nil {
				target = NextTarget(monster, targets)
			} else {
				target = NextTarget(target, targets)
			}
			targetX, targetY = target.X, target.Y
			continue // Switch target
		}
		CursorMovement(&targetX, &targetY, key)
		i = true
	}
	return turnSpent
}

func (c *Creature) Target2(b Board, o *Objects, cs Creatures, t *Creature) bool {
	/* Target is method of Creature, that takes game map, objects, and
	   creatures as arguments. Returns bool that serves as indicator if
	   action took some time or not.
	   This method is "the big one", general, for handling targeting.
	   In short, player starts targetting, line is drawn from player
	   to monster, then function waits for input (confirmation - "fire",
	   breaking the loop, or continuing).
	   Explicitly:
	   - creates list of all potential targets in fov
	    * tries to automatically last target, but
	    * if fails, it targets the nearest enemy
	   - draws line between source (receiver) and target (coords)
	    * creates new vector
	    * checks if it is valid - monsterHit should not be nil
	    * prints brensenham's line (ie so-called "vector")
	   - waits for player input
	    * if player cancels, function ends
	    * if player confirms, valley is shoot (in target, or empty space)
	    * if valley is shot in empty space, vector is extrapolated to check
	      if it will hit any target
	    * player can switch between targets as well; it targets
	      next target automatically; at first, only monsters that are
	      valid target (ie clean shot is possible), then monsters that
	      are in range and fov, but line of shot is not clear
	    * in other cases, game will try to move cursor; invalid input
	      is ignored */
	turnSpent := false
	if c.AmmoCurrent <= 0 {
		AddMessage("No more bolts in crossbow!")
		return turnSpent
	}
	var target *Creature
	targets := c.FindTargets(FOVLength, b, cs, *o)
	if LastTarget != nil && LastTarget != c &&
		IsInFOV(b, c.X, c.Y, LastTarget.X, LastTarget.Y) == true {
		target = LastTarget
	} else {
		var err error
		target, err = c.FindTarget(targets)
		if err != nil {
			fmt.Println(err)
		}
	}
	if t != nil {
		target = t
	}
	targetX, targetY := target.X, target.Y
	i := false
	for {
		vec, err := NewBrensenham(c.X, c.Y, targetX, targetY)
		if err != nil {
			fmt.Println(err)
		}
		_ = ComputeBrensenham(vec)
		valid, _, monsterHit, _ := ValidateBrensenham(vec, b, targets, *o)
		PrintBrensenham(vec, BrensenhamWhyTarget, BrensenhamColorGood, BrensenhamColorBad, b, *o, cs)
		if monsterHit != nil {
			msg := monsterHit.Name
			if monsterHit.FireResistance == FullAbility {
				msg = msg + " [[Fire res]]"
			} else if monsterHit.FireResistance == NoAbility {
				msg = msg + " [[Fire vuln]]"
			}
			if monsterHit.CanFly == FullAbility {
				msg = msg + " [[Flying]]"
			}
			if monsterHit.CanSwim == FullAbility && monsterHit.CanFly != FullAbility {
				msg = msg + " [[Can swim]]"
			}
			PrintLookingMessage(msg, i)
		}
		key := ReadInput()
		if key == blt.TK_ESCAPE {
			break
		}
		if AdjustSpell(key, c) == true {
			continue
		}
		if key == blt.TK_T || key == blt.TK_ENTER || key == blt.TK_MOUSE_LEFT {
			monsterAimed := FindMonsterByXY(targetX, targetY, cs)
			if monsterAimed != nil && monsterAimed != c && monsterAimed.HPCurrent > 0 && valid == true {
				LastTarget = monsterAimed
				c.AttackTarget(monsterAimed, o)
			} else {
				if monsterAimed == c {
					break // Do not hurt yourself.
				}
				if monsterHit != nil {
					if monsterHit.HPCurrent > 0 {
						LastTarget = monsterHit
						c.AttackTarget(monsterHit, o)
					}
				} else {
					vx, vy := FindBrensenhamDirection(vec)
					v := ExtrapolateBrensenham(vec, vx, vy)
					_, _, monsterHitIndirectly, _ := ValidateBrensenham(v, b, targets, *o)
					if monsterHitIndirectly != nil {
						c.AttackTarget(monsterHitIndirectly, o)
					}
				}
			}
			c.AmmoCurrent--
			turnSpent = true
			break
		} else if key == blt.TK_TAB {
			i = true
			monster := FindMonsterByXY(targetX, targetY, cs)
			if monster != nil {
				target = NextTarget(monster, targets)
			} else {
				target = NextTarget(target, targets)
			}
			targetX, targetY = target.X, target.Y
			continue // Switch target
		}
		CursorMovement(&targetX, &targetY, key)
		i = true
	}
	return turnSpent
}

func CursorMovement(x, y *int, key int) {
	/* CursorMovement is function that takes pointers to coords, and
	   int-based user input. It uses MoveCursor function to
	   modify original values. */
	switch key {
	case blt.TK_UP:
		MoveCursor(x, y, 0, -1)
	case blt.TK_RIGHT:
		MoveCursor(x, y, 1, 0)
	case blt.TK_DOWN:
		MoveCursor(x, y, 0, 1)
	case blt.TK_LEFT:
		MoveCursor(x, y, -1, 0)
	case blt.TK_MOUSE_MOVE:
		newX := blt.State(blt.TK_MOUSE_X)
		newY := blt.State(blt.TK_MOUSE_Y)
		if newX >= 0 && newX < MapSizeX && newY >= 0 && newY < MapSizeY {
			*x = newX
			*y = newY
		}
	}
}

func MoveCursor(x, y *int, dx, dy int) {
	/* Function MoveCursor takes pointers to coords, and
	   two other ints as direction indicators.
	   It adds direction to coordinate, checks if it is in
	   map bounds, and modifies original values accordingly.
	   This function is called by CursorMovement. */
	newX, newY := *x+dx, *y+dy
	if newX < 0 || newX >= MapSizeX {
		newX = *x
	}
	if newY < 0 || newY >= MapSizeY {
		newY = *y
	}
	*x, *y = newX, newY
}

func PrintCursor(area [][]int, b Board, o Objects, c Creatures) {
	blt.Clear()
	RenderAll(b, o, c)
	for _, v := range area {
		col := "light blue"
		if GlobalData.CurrentSchool == SchoolFire {
			col = "light red"
		} else if GlobalData.CurrentSchool == SchoolEarth {
			col = "dark orange"
		}
		blt.Layer(BoardLayer)
		blt.ClearArea(v[0], v[1], 1, 1)
		blt.Layer(LookLayer)
		blt.Print(v[0], v[1], "[color="+col+"]X")
	}
	blt.Refresh()
}

func SetCursorSize() [][]int {
	cursor := [][]int{  // CurrentSize = SizeSmall
		[]int{0, -1},
		[]int{-1, 0}, []int{0, 0}, []int{1, 0},
		[]int{0, 1},
	}
	if GlobalData.CurrentSize == SizeMedium {
		cursor = nil
		cursor = [][]int{
			[]int{0, -2},
			[]int{-1, -1}, []int{0, -1}, []int{1, -1},
			[]int{-2, 0}, []int{-1, 0}, []int{0, 0}, []int{1, 0}, []int{2, 0},
			[]int{-1, 1}, []int{0, 1}, []int{1, 1},
			[]int{0, 2},
		}
	} else if GlobalData.CurrentSize == SizeBig {
		cursor = nil
		cursor = [][]int{
			[]int{-2, -4}, []int{-1, -4}, []int{0, -4}, []int{1, -4}, []int{2, -4},
			[]int{-3, -3}, []int{-2, -3}, []int{-1, -3}, []int{0, -3}, []int{1, -3}, []int{2, -3}, []int{3, -3},
			[]int{-4, -2}, []int{-3, -2}, []int{-2, -2}, []int{-1, -2}, []int{0, -2}, []int{1, -2}, []int{2, -2}, []int{3, -2}, []int{4, -2},
			[]int{-4, -1}, []int{-3, -1}, []int{-2, -1}, []int{-1, -1}, []int{0, -1}, []int{1, -1}, []int{2, -1}, []int{3, -1}, []int{4, -1},
			[]int{-4, 0}, []int{-3, 0}, []int{-2, 0}, []int{-1, 0}, []int{0, 0}, []int{1, 0}, []int{2, 0}, []int{3, 0}, []int{4, 0},
			[]int{-4, 1}, []int{-3, 1}, []int{-2, 1}, []int{-1, 1}, []int{0, 1}, []int{1, 1}, []int{2, 1}, []int{3, 1}, []int{4, 1},
			[]int{-4, 2}, []int{-3, 2}, []int{-2, 2}, []int{-1, 2}, []int{0, 2}, []int{1, 2}, []int{2, 2}, []int{3, 2}, []int{4, 2},
			[]int{-3, 3}, []int{-2, 3}, []int{-1, 3}, []int{0, 3}, []int{1, 3}, []int{2, 3}, []int{3, 3},
			[]int{-2, 4}, []int{-1, 4}, []int{0, 4}, []int{1, 4}, []int{2, 4},
		}
	} else if GlobalData.CurrentSize == SizeHuge {
		cursor = nil
		cursor = [][]int{
			[]int{-3, -8}, []int{-2, -8}, []int{-1, -8}, []int{0, -8}, []int{1, -8}, []int{2, -8}, []int{3, -8},
			[]int{-5, -7}, []int{-4, -7}, []int{-3, -7}, []int{-2, -7}, []int{-1, -7}, []int{0, -7}, []int{1, -7}, []int{2, -7}, []int{3, -7}, []int{4, -7}, []int{5, -7},
			[]int{-6, -6}, []int{-5, -6}, []int{-4, -6}, []int{-3, -6}, []int{-2, -6}, []int{-1, -6}, []int{0, -6}, []int{1, -6}, []int{2, -6}, []int{3, -6}, []int{4, -6}, []int{5, -6}, []int{6,-6},
			[]int{-7, -5}, []int{-6, -5}, []int{-5, -5}, []int{-4, -5}, []int{-3, -5}, []int{-2, -5}, []int{-1, -5}, []int{0, -5}, []int{1, -5}, []int{2, -5}, []int{3, -5}, []int{4, -5}, []int{5, -5}, []int{6, -5}, []int{7, -5},
			[]int{-7, -4}, []int{-6, -4}, []int{-5, -4}, []int{-4, -4}, []int{-3, -4}, []int{-2, -4}, []int{-1, -4}, []int{0, -4}, []int{1, -4}, []int{2, -4}, []int{3, -4}, []int{4, -4}, []int{5, -4}, []int{6, -4}, []int{7, -4},
			[]int{-8, -3}, []int{-7, -3}, []int{-6, -3}, []int{-5, -3}, []int{-4, -3}, []int{-3, -3}, []int{-2, -3}, []int{-1, -3}, []int{0, -3}, []int{1, -3}, []int{2, -3}, []int{3, -3}, []int{4, -3}, []int{5, -3}, []int{6, -3}, []int{7, -3}, []int{8, -3},
			[]int{-8, -2}, []int{-7, -2}, []int{-6, -2}, []int{-5, -2}, []int{-4, -2}, []int{-3, -2}, []int{-2, -2}, []int{-1, -2}, []int{0, -2}, []int{1, -2}, []int{2, -2}, []int{3, -2}, []int{4, -2}, []int{5, -2}, []int{6, -2}, []int{7, -2}, []int{8, -2},
			[]int{-8, -1}, []int{-7, -1}, []int{-6, -1}, []int{-5, -1}, []int{-4, -1}, []int{-3, -1}, []int{-2, -1}, []int{-1, -1}, []int{0, -1}, []int{1, -1}, []int{2, -1}, []int{3, -1}, []int{4, -1}, []int{5, -1}, []int{6, -1}, []int{7, -1}, []int{8, -1},
			[]int{-8, 0}, []int{-7, 0}, []int{-6, 0}, []int{-5, 0}, []int{-4, 0}, []int{-3, 0}, []int{-2, 0}, []int{-1, 0}, []int{0, 0}, []int{1, 0}, []int{2, 0}, []int{3, 0}, []int{4, 0}, []int{5, 0}, []int{6, 0}, []int{7, 0}, []int{8, 0},
			[]int{-8, 1}, []int{-7, 1}, []int{-6, 1}, []int{-5, 1}, []int{-4, 1}, []int{-3, 1}, []int{-2, 1}, []int{-1, 1}, []int{0, 1}, []int{1, 1}, []int{2, 1}, []int{3, 1}, []int{4, 1}, []int{5, 1}, []int{6, 1}, []int{7, 1}, []int{8, 1},
			[]int{-8, 2}, []int{-7, 2}, []int{-6, 2}, []int{-5, 2}, []int{-4, 2}, []int{-3, 2}, []int{-2, 2}, []int{-1, 2}, []int{0, 2}, []int{1, 2}, []int{2, 2}, []int{3, 2}, []int{4, 2}, []int{5, 2}, []int{6, 2}, []int{7, 2}, []int{8, 2},
			[]int{-8, 3}, []int{-7, 3}, []int{-6, 3}, []int{-5, 3}, []int{-4, 3}, []int{-3, 3}, []int{-2, 3}, []int{-1, 3}, []int{0, 3}, []int{1, 3}, []int{2, 3}, []int{3, 3}, []int{4, 3}, []int{5, 3}, []int{6, 3}, []int{7, 3}, []int{8, 3},
			[]int{-7, 4}, []int{-6, 4}, []int{-5, 4}, []int{-4, 4}, []int{-3, 4}, []int{-2, 4}, []int{-1, 4}, []int{0, 4}, []int{1, 4}, []int{2, 4}, []int{3, 4}, []int{4, 4}, []int{5, 4}, []int{6, 4}, []int{7, 4},
			[]int{-7, 5}, []int{-6, 5}, []int{-5, 5}, []int{-4, 5}, []int{-3, 5}, []int{-2, 5}, []int{-1, 5}, []int{0, 5}, []int{1, 5}, []int{2, 5}, []int{3, 5}, []int{4, 5}, []int{5, 5}, []int{6, 5}, []int{7, 5},
			[]int{-6, 6}, []int{-5, 6}, []int{-4, 6}, []int{-3, 6}, []int{-2, 6}, []int{-1, 6}, []int{0, 6}, []int{1, 6}, []int{2, 6}, []int{3, 6}, []int{4, 6}, []int{5, 6}, []int{6,6},
			[]int{-5, 7}, []int{-4, 7}, []int{-3, 7}, []int{-2, 7}, []int{-1, 7}, []int{0, 7}, []int{1, 7}, []int{2, 7}, []int{3, 7}, []int{4, 7}, []int{5, 7},
			[]int{-3, 8}, []int{-2, 8}, []int{-1, 8}, []int{0, 8}, []int{1, 8}, []int{2, 8}, []int{3, 8},
		}
	}
	return cursor
}

func AdjustSpell(key int, c *Creature) bool {
	keyCorrect := true
	if key == blt.TK_F1 {
		GlobalData.CurrentSchool = SchoolWater
		AddMessage("You invoke water aura.")
		c.Colors = []string{"#73C2FB"}
	} else if key == blt.TK_F2 {
		GlobalData.CurrentSchool = SchoolFire
		AddMessage("You invoke fire aura.")
		c.Colors = []string{"#FF7F7F"}
	} else if key == blt.TK_F3 {
		GlobalData.CurrentSchool = SchoolEarth
		AddMessage("You invoke earth aura.")
		c.Colors = []string{"#D2B48C"}
	} else if key == blt.TK_KP_DIVIDE || key == blt.TK_LBRACKET {
		if GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSchool = SchoolEarth
			AddMessage("You invoke earth aura.")
			c.Colors = []string{"#D2B48C"}
		} else if GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSchool = SchoolWater
			AddMessage("You invoke water aura.")
			c.Colors = []string{"#73C2FB"}
		} else {
			GlobalData.CurrentSchool = SchoolFire
			AddMessage("You invoke fire aura.")
			c.Colors = []string{"#FF7F7F"}
		}
	} else if key == blt.TK_KP_MULTIPLY || key == blt.TK_RBRACKET || key == blt.TK_MOUSE_RIGHT {
		if GlobalData.CurrentSchool == SchoolWater {
			GlobalData.CurrentSchool = SchoolFire
			AddMessage("You invoke fire aura.")
			c.Colors = []string{"#FF7F7F"}
		} else if GlobalData.CurrentSchool == SchoolFire {
			GlobalData.CurrentSchool = SchoolEarth
			AddMessage("You invoke earth aura.")
			c.Colors = []string{"#D2B48C"}
		} else {
			GlobalData.CurrentSchool = SchoolWater
			AddMessage("You invoke water aura.")
			c.Colors = []string{"#73C2FB"}
		}
	} else if key == blt.TK_1 {
		GlobalData.CurrentSize = SizeSmall
	} else if key == blt.TK_2 {
		GlobalData.CurrentSize = SizeMedium
	} else if key == blt.TK_3 {
		GlobalData.CurrentSize = SizeBig
	} else if key == blt.TK_4 {
		GlobalData.CurrentSize = SizeHuge
	} else if key == blt.TK_KP_MINUS || key == blt.TK_MINUS || ((key == blt.TK_MOUSE_SCROLL) && blt.State(blt.TK_MOUSE_WHEEL) == 1) {
		if GlobalData.CurrentSize == SizeSmall {
			GlobalData.CurrentSize = SizeHuge
		} else {
			GlobalData.CurrentSize--
		}
	} else if key == blt.TK_KP_PLUS || key == blt.TK_EQUALS || ((key == blt.TK_MOUSE_SCROLL) && blt.State(blt.TK_MOUSE_WHEEL) == -1) {
		if GlobalData.CurrentSize == SizeHuge {
			GlobalData.CurrentSize = SizeSmall
		} else {
			GlobalData.CurrentSize++
		}
	} else {
		keyCorrect = false
	}
	return keyCorrect
}

func (c *Creature) FindTargets(length int, b Board, cs Creatures, o Objects) Creatures {
	/* FindTargets is method of Creature that takes several arguments:
	   length (that is supposed to be max range of attack), and: map, creatures,
	   objects. Returns list of creatures.
	   At first, method creates list of all monsters im c's field of view.
	   Then, this list is divided to two: first, with all "valid" targets
	   (clean (without obstacles) line between c and target) and second,
	   with all other monsters that remains in fov.
	   Both slices are sorted by distance from receiver, then merged.
	   It is necessary for autotarget feature - switching between targets
	   player will start from the nearest valid target, to the farthest valid target;
	   THEN, it will start to target "invalid" targets - again,
	   from nearest to farthest one. */
	targets := c.MonstersInFov(b, cs)
	targetable, unreachable := c.MonstersInRange(b, targets, o, length)
	sort.Slice(targetable, func(i, j int) bool {
		return targetable[i].DistanceBetweenCreatures(c) <
			targetable[j].DistanceBetweenCreatures(c)
	})
	sort.Slice(unreachable, func(i, j int) bool {
		return unreachable[i].DistanceBetweenCreatures(c) <
			unreachable[j].DistanceBetweenCreatures(c)
	})
	targets = nil
	targets = append(targets, targetable...)
	targets = append(targets, unreachable...)
	return targets
}

func (c *Creature) FindTarget(targets Creatures) (*Creature, error) {
	/* FindTarget is method of Creature that takes Creatures as arguments.
	   It returns specific Creature and error.
	   "targets" is supposed to be slice of Creature in player's fov,
	   sorted as explained in FindTargets docstring.
	   If this slice is empty, the target is set to receiver. If not,
	   it tries to target lastly targeted Creature. If it is not possible,
	   it targets first element of slice, and marks it as LastTarget.
	   This method throws an error if it can not find any target,
	   even including receiver. */
	var target *Creature
	if len(targets) == 0 {
		target = c
	} else {
		if LastTarget != nil && CreatureIsInSlice(LastTarget, targets) {
			target = LastTarget
		} else {
			target = targets[0]
			LastTarget = target
		}
	}
	var err error
	if target == nil {
		txt := TargetNilError(c, targets)
		err = errors.New("Could not find target, even the 'self' one." + txt)
	}
	return target, err
}

func NextTarget(target *Creature, targets Creatures) *Creature {
	/* Function NextTarget takes specific creature (target) and slice of creatures
	   (targets) as arguments. It tries to find the *next* target (used
	   with switching between targets, for example using Tab key).
	   At the end, it returns the next creature. */
	i, _ := FindCreatureIndex(target, targets)
	var t *Creature
	length := len(targets)
	if length > i+1 {
		t = targets[i+1]
	} else if length == 0 {
		t = target
	} else {
		t = targets[0]
	}
	return t
}

func (c *Creature) MonstersInRange(b Board, cs Creatures, o Objects,
	length int) (Creatures, Creatures) {
	/* MonstersInRange is method of Creature. It takes global map, Creatures
	   and Objects, and length (range indicator) as its arguments. It returns
	   two slices - one with monsters that are in range, and one with
	   monsters out of range.
	   At first, two empty slices are created, then function starts iterating
	   through Creatures from argument. It creates new vector from source (c)
	   to target, adds monster to proper slice. It also validates vector
	   (ie, won't add monster hidden behind wall) and skips all dead monsters. */
	var inRange = Creatures{}
	var outOfRange = Creatures{}
	for i, v := range cs {
		vec, err := NewBrensenham(c.X, c.Y, v.X, v.Y)
		if err != nil {
			fmt.Println(err)
		}
		if ComputeBrensenham(vec) <= length+1 { // "+1" is necessary due Brensenham values.
			valid, _, _, _ := ValidateBrensenham(vec, b, cs, o)
			if cs[i].HPCurrent <= 0 {
				continue
			}
			if valid == true {
				inRange = append(inRange, cs[i])
			} else {
				outOfRange = append(outOfRange, cs[i])
			}
		}
	}
	return inRange, outOfRange
}

func ZeroLastTarget(c *Creature) {
	/* LastTarget is global variable (will be incorporated into
	   player struct in future). Function ZeroLastTarget changes
	   last target to nil, is last target matches creature
	   passed as argument. */
	if LastTarget == c {
		LastTarget = nil
	}
}
