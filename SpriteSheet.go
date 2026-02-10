package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	CellSize = 32
	GridCols = 4
	GridRows = 4
)

// GenerateSpritesheet creates a 4x4 grid of 32x32 colored rectangles using
// Citruszest palette colors and returns it as a GPU texture.
func GenerateSpritesheet() rl.Texture2D {
	colors := [GridRows * GridCols]rl.Color{
		// Row 0: base colors
		rl.NewColor(0x23, 0x23, 0x23, 0xFF), // black
		rl.NewColor(0xFF, 0x54, 0x54, 0xFF), // red
		rl.NewColor(0x00, 0xCC, 0x7A, 0xFF), // green
		rl.NewColor(0xFF, 0xD7, 0x00, 0xFF), // yellow
		// Row 1: more base colors
		rl.NewColor(0xFF, 0x74, 0x31, 0xFF), // orange
		rl.NewColor(0x00, 0xBF, 0xFF, 0xFF), // blue
		rl.NewColor(0x00, 0xFF, 0xFF, 0xFF), // cyan
		rl.NewColor(0xBF, 0xBF, 0xBF, 0xFF), // white
		// Row 2: bright variants
		rl.NewColor(0x76, 0x7C, 0x77, 0xFF), // bright_black
		rl.NewColor(0xFF, 0x1A, 0x75, 0xFF), // bright_red
		rl.NewColor(0x1A, 0xFF, 0xA3, 0xFF), // bright_green
		rl.NewColor(0xFF, 0xFF, 0x00, 0xFF), // bright_yellow
		// Row 3: accents
		rl.NewColor(0x9A, 0xDC, 0xFF, 0xFF), // baby_blue
		rl.NewColor(0xFF, 0xF2, 0xB3, 0xFF), // lemon_yellow
		rl.NewColor(0xB2, 0xF3, 0xAC, 0xFF), // aurora
		rl.NewColor(0xAF, 0x74, 0xEE, 0xFF), // violet
	}

	img := rl.GenImageColor(CellSize*GridCols, CellSize*GridRows, rl.Blank)
	for i, c := range colors {
		col := int32(i % GridCols)
		row := int32(i / GridCols)
		rl.ImageDrawRectangle(img, col*CellSize, row*CellSize, CellSize, CellSize, c)
	}

	texture := rl.LoadTextureFromImage(img)
	rl.UnloadImage(img)
	return texture
}
