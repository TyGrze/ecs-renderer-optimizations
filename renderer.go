package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type SpriteRenderer struct {
	filter       *ecs.Filter2[Position, Sprite]
	cameraBounds *ecs.Resource[CameraBounds]

	sheet         rl.Texture2D
	RenderedCount int
}

func NewSpriteRenderer(w *ecs.World, sheet rl.Texture2D) *SpriteRenderer {
	filter := ecs.NewFilter2[Position, Sprite](w)
	cameraRes := ecs.NewResource[CameraBounds](w)
	return &SpriteRenderer{
		filter:       filter,
		sheet:        sheet,
		cameraBounds: &cameraRes,
	}
}

func (s *SpriteRenderer) Update(w *ecs.World) {
	s.RenderedCount = 0
	query := s.filter.Query()
	cameraBound := s.cameraBounds.Get()

	minX := cameraBound.MinX
	minY := cameraBound.MinY
	maxX := cameraBound.MaxX
	maxY := cameraBound.MaxY

	for query.Next() {
		pos, sprite := query.Get()

		// Culling
		if pos.X+32 < minX || pos.X > maxX ||
			pos.Y+32 < minY || pos.Y > maxY {
			continue
		}

		s.RenderedCount++
		rl.DrawTextureRec(s.sheet, rl.NewRectangle(float32(sprite.X*32), float32(sprite.Y*32), 32, 32),
			rl.NewVector2(float32(pos.X), float32(pos.Y)), rl.White)
	}
}
