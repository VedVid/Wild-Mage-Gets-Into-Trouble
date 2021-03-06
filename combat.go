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

func (c *Creature) AttackTarget(t *Creature, o *Objects) {
	/* Method Attack handles damage rolls for combat. Receiver "c" is attacker,
	   argument "t" is target. Including o *Objects is necessary for dropping
	   loot by dead enemies.
	   Critical hit is if attack roll is the same as receiver
	   attack attribute.
	   Result of attack is displayed in combat log, but messages need more polish. */
	att := RandInt(c.Attack)
	att2 := 0
	def := t.Defense
	crit := false
	dmg := 0
	if c.AIType == PlayerAI {
		if att == c.Attack {     //critical hit!
			crit = true
			att2 = RandInt(c.Attack)
		}
		switch {
		case att < def: // Attack score if lower than target defense.
			if crit == true {
				dmg = att2
			}
		case att == def: // Attack score is equal to target defense.
			if crit == false {
				dmg = 1 // It's just a scratch...
			} else {
				dmg = att
			}
		case att > def: // Attack score is bigger than target defense.
			if crit == false {
				dmg = att
			} else {
				dmg = att + att2 // Critical attack!
			}
		}
		if GlobalData.CurrentSchool == SchoolFire && t.FireResistance == NoAbility {
			dmg *= 2
		} else if GlobalData.CurrentSchool == SchoolWater && t.CanSwim == NoAbility {
			dmg *= 2
		} else if GlobalData.CurrentSchool == SchoolEarth && t.CanFly == NoAbility {
			dmg *= 2
		}
	} else {
		if att == c.Attack {
			crit = true
			att2 = RandInt(c.Attack)
		}
		dmg = att + att2 - def
	}
	if dmg < 0 {
		dmg = 0
	}
	t.TakeDamage(dmg, o)
}

func (c *Creature) TakeDamage(dmg int, o *Objects) {
	/* Method TakeDamage has *Creature as receiver and takes damage integer
	   as argument. dmg value is deducted from Creature current HP.
	   If HPCurrent is below zero after taking damage, Creature dies.
	   o as map objects is passed with Die to handle dropping loot. */
	c.HPCurrent -= dmg
	if c.HPCurrent <= 0 {
		c.Die(o)
	}
}

func CheckMagic(b Board, c Creatures, o *Objects) {
	if GlobalData.TurnsSpent % ManaRegenDiv == 0 {
		c[0].ManaCurrent++
		if c[0].ManaCurrent > c[0].ManaMax {
			c[0].ManaCurrent = c[0].ManaMax
		}
	}
	for i := 0; i < len(c); i++ {
		monster := c[i]
		x := monster.X
		y := monster.Y
		if b[x][y].Fire > 0 {
			if i != 0 {
				if monster.FireResistance == NoAbility {
					monster.TakeDamage(999, o)
				} else if monster.FireResistance == PartialAbility {
					monster.TakeDamage(1, o)
				}
			} else {
				if GlobalData.CurrentSchool != SchoolFire {
					monster.TakeDamage(1, o)
				}
			}
		}
		if b[x][y].Flooded > 0 {
			if i != 0 {
				if monster.CanSwim == NoAbility {
					monster.TakeDamage(999, o)
				} else if monster.CanSwim == PartialAbility {
					monster.TakeDamage(1, o)
				}
			} else {
				if GlobalData.CurrentSchool != SchoolWater {
					monster.TakeDamage(1, o)
				}
			}
		}
		if b[x][y].Chasm > 0 {
			if i != 0 {
				if monster.CanFly != FullAbility {
					monster.TakeDamage(999, o)
					monster.Chars = []string{""}
				}
			} else {
				if GlobalData.CurrentSchool != SchoolEarth {
					monster.TakeDamage(999, o)
				}
			}
		}
	}
}
