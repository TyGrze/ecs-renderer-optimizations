package main

import (
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

const (
	freshBit  = 0b100
	indexMask = 0b011
)

type TrippleEntityBuffer struct {
	buffers [3][]rl.Matrix
	counts  [3]int
	middle  atomic.Uint32
	wrIdx   uint32
	rdIdx   uint32
}

func NewTrippleEntityBuffer(capacity int) *TrippleEntityBuffer {
	tb := &TrippleEntityBuffer{}
	for i := range tb.buffers {
		tb.buffers[i] = make([]rl.Matrix, capacity)
	}
	tb.wrIdx = 0
	tb.middle.Store(1)
	tb.rdIdx = 2

	return tb
}

func (tb *TrippleEntityBuffer) WriteBuffer() *[]rl.Matrix {
	return &tb.buffers[tb.wrIdx]
}

func (tb *TrippleEntityBuffer) ReadBuffer() []rl.Matrix {
	return tb.buffers[tb.rdIdx]
}

func (tb *TrippleEntityBuffer) SwapWriter() {
	old := tb.middle.Swap(tb.wrIdx | freshBit)
	tb.wrIdx = old & indexMask
}

func (tb *TrippleEntityBuffer) SwapReader() bool {
	for {
		mid := tb.middle.Load()

		if mid&freshBit == 0 {
			return false // No New Data
		}
		if tb.middle.CompareAndSwap(mid, tb.rdIdx) {
			tb.rdIdx = mid & indexMask
			return true
		}
		// If we get here that means the compare failed and there is new data to read
		// For loop will retry
	}
}

type SpriteRenderer struct {
	shader   rl.Shader
	mesh     rl.Mesh
	material rl.Material

	tb *TrippleEntityBuffer
}

type SpriteRendererSystem struct {
	filter       *ecs.Filter2[Position, Sprite]
	cameraBounds *ecs.Resource[CameraBounds]

	tb            *TrippleEntityBuffer
	RenderedCount int
}

func NewSpriteRendererSystem(w *ecs.World, tb *TrippleEntityBuffer) *SpriteRendererSystem {
	filter := ecs.NewFilter2[Position, Sprite](w)
	cameraRes := ecs.NewResource[CameraBounds](w)

	return &SpriteRendererSystem{
		filter:       filter,
		cameraBounds: &cameraRes,
		tb:           tb,
	}
}

func (s *SpriteRendererSystem) Update(w *ecs.World) {
	s.RenderedCount = 0
	query := s.filter.Query()
	cameraBound := s.cameraBounds.Get()

	minX := cameraBound.MinX
	minY := cameraBound.MinY
	maxX := cameraBound.MaxX
	maxY := cameraBound.MaxY

	buf := s.tb.WriteBuffer()
	*buf = (*buf)[:0] // Reset length but keeps capacity

	for query.Next() {
		pos, sprite := query.Get()

    // Simple AABB Culling
		if pos.X+CellSize < minX || pos.X > maxX ||
			pos.Y+CellSize < minY || pos.Y > maxY {
			continue
		}

		// Direct matrix construction - equivalent to MatrixMultiply(Scale, Translate)
		// but avoids 2 matrix builds + a full 4x4 multiply per entity.
		*buf = append(*buf, rl.Matrix{
			M0:  CellSize,
			M5:  float32(sprite.Y*GridCols + sprite.X), // sprite index encoded for shader
			M10: CellSize,
			M12: pos.X + CellSize/2,
			M14: pos.Y + CellSize/2,
			M15: 1,
		})
		s.RenderedCount++
	}

	// Publihs the transforms
	s.tb.counts[s.tb.wrIdx] = s.RenderedCount
	s.tb.SwapWriter()
}

func NewSpriteRenderer(tb *TrippleEntityBuffer, sheet rl.Texture2D) *SpriteRenderer {
	shader := rl.LoadShader("shaders/instancing.vs", "")
	shader.UpdateLocation(rl.ShaderLocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.ShaderLocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))

	mesh := rl.GenMeshPlane(1, 1, 1, 1)

	material := rl.LoadMaterialDefault()
	material.Shader = shader
	material.GetMap(rl.MapDiffuse).Texture = sheet

	return &SpriteRenderer{
		shader:   shader,
		mesh:     mesh,
		material: material,
		tb:       tb,
	}
}

func (s *SpriteRenderer) Unload() {
	rl.UnloadShader(s.shader)
	rl.UnloadMesh(&s.mesh)
}

func (s *SpriteRenderer) Render() {
  // If needing to optimize more in the future we can check the bool from SwapReader and if false we could re-use the VBO
  // But this would require us to create our own OpenGL instancing instead of uring rl
	s.tb.SwapReader()
	transforms := s.tb.ReadBuffer()
	count := s.tb.counts[s.tb.rdIdx]

  // NOTE: If this needs to be speed up even more I submited a merge request to add rlgl bindings to raylib-go
  // The slowest part of this RN is in rl, the only way to speed it up more is to deal with the VBO VAO managment and calls ourself
  // https://github.com/gen2brain/raylib-go/pull/537
	if count > 0 {
		rl.DrawMeshInstanced(s.mesh, s.material, transforms[:count], count)
	}
}
