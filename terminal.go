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
	"runtime"
	"strconv"

	blt "bearlibterminal"
)

const (
	// Setting BearLibTerminal window.
	WindowSizeX = 30
	WindowSizeY = 30
	MapSizeX    = 25
	MapSizeY    = 25
	UIPosX      = MapSizeX
	UIPosY      = 0
	UISizeX     = WindowSizeX - MapSizeX
	UISizeY     = WindowSizeY
	LogSizeX    = WindowSizeX
	LogSizeY    = WindowSizeY - MapSizeY
	LogPosX     = 0
	LogPosY     = MapSizeY
	GameTitle   = "Wild Mage Gets Into Trouble"
	GameVersion = "7DRL 2022"
	FontName    = "Deferral-Square.ttf"
	FontUI      = "PTMono-Regular.ttf"
	FontSize    = 18
)

var ActualFontSize = FontSize

func constrainThreads() {
	/* Constraining processor and threads is necessary,
	   because BearLibTerminal often crashes otherwise. */
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
}

func InitializeBLT() {
	/* Constraining threads and setting BearLibTerminal window. */
	constrainThreads()
	blt.Open()
	sizeX, sizeY := strconv.Itoa(WindowSizeX), strconv.Itoa(WindowSizeY)
	sizeFont := strconv.Itoa(ActualFontSize)
	window := "window: size=" + sizeX + "x" + sizeY
	gameFont := "game font: " + FontName + ", size=" + sizeFont
	blt.Set(gameFont)
	uiFont := "ui font: " + FontUI + ", size=" + sizeFont
	blt.Set(uiFont)
	blt.Set(window + ", title=' " + GameTitle + " " + GameVersion +
		"'; font: " + FontName + ", size=" + sizeFont)
	blt.Set("input.filter={keyboard, mouse+}")
	blt.Clear()
	blt.Refresh()
}
