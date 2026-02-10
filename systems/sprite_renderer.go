package systems

import (
	"ecs_test/components"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type SpriteRenderer struct {
	filter *ecs.Filter2[components.Position, components.Sprite]

  sheet rl.Texture2D
}

func NewSpriteRenderer(w *ecs.World, sheet rl.Texture2D) *SpriteRenderer {
	filter := ecs.NewFilter2[components.Position, components.Sprite](w)
	return &SpriteRenderer{
		filter: filter,
		sheet:  sheet,
	}
}

func (s *SpriteRenderer) Update(w *ecs.World) {
	query := s.filter.Query()

	for query.Next() {
		pos, sprite := query.Get()
		rl.DrawTextureRec(s.sheet, rl.NewRectangle(float32(sprite.X*32), float32(sprite.Y*32), 32, 32),
			rl.NewVector2(float32(pos.X), float32(pos.Y)), rl.White)
	}
}
