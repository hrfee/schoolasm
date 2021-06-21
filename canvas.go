// +build !nocanvas

package main

import (
	"fmt"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

var CANVAS = true

// Scancodes for arrow keys. Names correspond to "name" arg in wnd.KeyUp/KeyDown.
const (
	ArrowUp    = 82
	ArrowDown  = 81
	ArrowLeft  = 80
	ArrowRight = 79
)

func newCanvas(width, height, scale int, memStart addr, mem *memory) {
	wnd, cv, err := sdlcanvas.CreateWindow(width*scale, height*scale, "schoolasm")
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %v", err))
	}
	defer wnd.Destroy()
	// addresses for key events
	size := addr(width * height)
	keyUp := memStart + size + 1
	keyDown := keyUp + 1
	keyLeft := keyDown + 1
	keyRight := keyLeft + 1
	wnd.KeyUp = func(scancode int, rn rune, name string) {
		switch scancode {
		case ArrowUp:
			mem[keyUp] = 0
		case ArrowDown:
			mem[keyDown] = 0
		case ArrowLeft:
			mem[keyLeft] = 0
		case ArrowRight:
			mem[keyRight] = 0
		}
		Println("UP", scancode, rn, name)
	}
	wnd.KeyDown = func(scancode int, rn rune, name string) {
		switch scancode {
		case ArrowUp:
			mem[keyUp] = 1
		case ArrowDown:
			mem[keyDown] = 1
		case ArrowLeft:
			mem[keyLeft] = 1
		case ArrowRight:
			mem[keyRight] = 1
		}
		Println("DOWN", scancode, rn, name)
	}
	// rows, then columns
	wnd.MainLoop(func() {
		offset := 0
		for sy := 0; sy <= height; sy++ {
			y := sy * scale
			for sx := 0; sx < width; sx++ {
				x := sx * scale
				off := addr(sx) + addr(offset)
				if off >= size {
					break
				}
				white := mem[memStart+off]
				if white == 1 {
					cv.SetFillStyle("#ffffff")
				} else {
					cv.SetFillStyle("#000000")
				}
				cv.FillRect(float64(x), float64(y), float64(scale), float64(scale))
			}
			offset += width
		}
	})
}
