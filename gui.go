package main

import (
	"fmt"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

func newGUI(width, height, scale int, memStart addr, mem *memory) {
	wnd, cv, err := sdlcanvas.CreateWindow(width*scale, height*scale, "schoolasm")
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %v", err))
	}
	defer wnd.Destroy()
	// rows, then columns
	wnd.MainLoop(func() {
		offset := 0
		for sy := 0; sy <= height; sy++ {
			y := sy * scale
			for sx := 0; sx < width; sx++ {
				x := sx * scale
				white := mem[memStart+addr(sx)+addr(offset)]
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
