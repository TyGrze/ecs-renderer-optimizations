package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

type SpriteRenderer struct {
	filter       *ecs.Filter2[Position, Sprite]
	cameraBounds *ecs.Resource[CameraBounds]

	shader   rl.Shader
	mesh     rl.Mesh
	material rl.Material

	transforms    []rl.Matrix
	RenderedCount int
}

func NewSpriteRenderer(w *ecs.World, sheet rl.Texture2D) *SpriteRenderer {
	filter := ecs.NewFilter2[Position, Sprite](w)
	cameraRes := ecs.NewResource[CameraBounds](w)

	shader := rl.LoadShader("shaders/instancing.vs", "")
	shader.UpdateLocation(rl.ShaderLocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.ShaderLocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))

	mesh := rl.GenMeshPlane(1, 1, 1, 1)

	material := rl.LoadMaterialDefault()
	material.Shader = shader
	material.GetMap(rl.MapDiffuse).Texture = sheet

	return &SpriteRenderer{
		filter:       filter,
		cameraBounds: &cameraRes,
		shader:       shader,
		mesh:         mesh,
		material:     material,
		transforms:   make([]rl.Matrix, entityCount),
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

		if pos.X+CellSize < minX || pos.X > maxX ||
			pos.Y+CellSize < minY || pos.Y > maxY {
			continue
		}

		// Translate so the quad's top-left aligns with the entity position.
		// GenMeshPlane is centered at origin, so offset by half CellSize.
		translate := rl.MatrixTranslate(pos.X+CellSize/2, 0, pos.Y+CellSize/2)
		scale := rl.MatrixScale(CellSize, 1, CellSize)
		transform := rl.MatrixMultiply(scale, translate)

		// Encode sprite index in M5 (instanceTransform[1][1])
		transform.M5 = float32(sprite.Y*GridCols + sprite.X)

		s.transforms[s.RenderedCount] = transform
		s.RenderedCount++
	}

	if s.RenderedCount > 0 {
		rl.DrawMeshInstanced(s.mesh, s.material, s.transforms[:s.RenderedCount], s.RenderedCount)
	}
}

func (s *SpriteRenderer) Unload() {
	rl.UnloadShader(s.shader)
	rl.UnloadMesh(&s.mesh)
}
