package main

import (
	"log"
	"os"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

const (
	freshBit  = 0b100
	indexMask = 0b011
)

type TrippleEntityBuffer struct {
	buffers [3][]rl.Vector2
	middle  atomic.Uint32
	wrIdx   uint32
	rdIdx   uint32
}

func NewTrippleEntityBuffer(capacity int) *TrippleEntityBuffer {
	tb := &TrippleEntityBuffer{}
	for i := range tb.buffers {
		tb.buffers[i] = make([]rl.Vector2, 0, capacity)
	}
	tb.wrIdx = 0
	tb.middle.Store(1)
	tb.rdIdx = 2

	return tb
}

func (tb *TrippleEntityBuffer) WriteBuffer() *[]rl.Vector2 {
	return &tb.buffers[tb.wrIdx]
}

func (tb *TrippleEntityBuffer) ReadBuffer() []rl.Vector2 {
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

// QuadVertex matches the interleaved VBO layout: vec3 position + vec2 texcoord
type QuadVertex struct {
	Pos rl.Vector3
	Tex rl.Vector2
}

type SpriteRenderer struct {
	shaderID    uint32
	vao         uint32
	quadVBO     uint32
	ebo         uint32
	spriteVBO   uint32 // static — sprite indices, set once
	positionVBO uint32 // dynamic — positions, updated per frame
	mvpLoc      int32
	cellSizeLoc int32
	textureLoc  int32
	texture     rl.Texture2D
	entityCount int32
	tb          *TrippleEntityBuffer
}

type SpriteRendererSystem struct {
	filter        *ecs.Filter2[Position, Sprite]
	tb            *TrippleEntityBuffer
	RenderedCount int
}

func NewSpriteRendererSystem(w *ecs.World, tb *TrippleEntityBuffer) *SpriteRendererSystem {
	filter := ecs.NewFilter2[Position, Sprite](w)

	return &SpriteRendererSystem{
		filter: filter,
		tb:     tb,
	}
}

func (s *SpriteRendererSystem) Update(w *ecs.World) {
	buf := s.tb.WriteBuffer()
	*buf = (*buf)[:0]

	query := s.filter.Query()
	for query.Next() {
		pos, _ := query.Get()
		*buf = append(*buf, rl.NewVector2(pos.X, pos.Y))
	}

	s.RenderedCount = len(*buf)
	s.tb.SwapWriter()
}

// BuildSpriteIndices iterates all entities once at startup and returns sprite indices for the static VBO.
func (s *SpriteRendererSystem) BuildSpriteIndices(w *ecs.World) []float32 {
	query := s.filter.Query()
	indices := make([]float32, 0, s.RenderedCount)

	for query.Next() {
		_, sprite := query.Get()
		indices = append(indices, float32(sprite.Y*GridCols+sprite.X))
	}

	return indices
}

func NewSpriteRenderer(tb *TrippleEntityBuffer, sheet rl.Texture2D, spriteIndices []float32, entityCount int32) *SpriteRenderer {
	// Load shader from files
	vsCode, err := os.ReadFile("shaders/instancing.vs")
	if err != nil {
		log.Fatalf("failed to read vertex shader: %v", err)
	}
	fsCode, err := os.ReadFile("shaders/instancing.fs")
	if err != nil {
		log.Fatalf("failed to read fragment shader: %v", err)
	}
	shaderID := rl.LoadShaderCode(string(vsCode), string(fsCode))

	// Get uniform locations
	mvpLoc := rl.GetLocationUniform(shaderID, "mvp")
	cellSizeLoc := rl.GetLocationUniform(shaderID, "cellSize")
	textureLoc := rl.GetLocationUniform(shaderID, "texture0")

	// Set static cellSize uniform
	rl.EnableShader(shaderID)
	rl.SetUniform(cellSizeLoc, []float32{CellSize}, int32(rl.ShaderUniformFloat), 1)
	// Set texture0 sampler to slot 0
	rl.SetUniform(textureLoc, []int32{0}, int32(rl.ShaderUniformInt), 1)
	rl.DisableShader()

	// Create VAO
	vao := rl.LoadVertexArray()
	rl.EnableVertexArray(vao)

	// Quad geometry: unit quad on XZ plane, Y=0
	quadVertices := []QuadVertex{
		{Pos: rl.NewVector3(-0.5, 0, -0.5), Tex: rl.NewVector2(0, 0)},
		{Pos: rl.NewVector3(0.5, 0, -0.5), Tex: rl.NewVector2(1, 0)},
		{Pos: rl.NewVector3(-0.5, 0, 0.5), Tex: rl.NewVector2(0, 1)},
		{Pos: rl.NewVector3(0.5, 0, 0.5), Tex: rl.NewVector2(1, 1)},
	}
	quadIndices := []uint16{0, 2, 1, 1, 2, 3}

	// Quad VBO (static) — location 0: vec3 position, location 1: vec2 texcoord
	quadVBO := rl.LoadVertexBuffer(quadVertices, false)
	rl.SetVertexAttributes(quadVertices, []rl.VertexAttributesConfig{
		{Field: "Pos", Attribute: 0},
		{Field: "Tex", Attribute: 1},
	})

	// EBO (static)
	ebo := rl.LoadVertexBufferElements(quadIndices, false)

	// Sprite index VBO (static, divisor=1) — location=2, 1 float per instance
	spriteVBO := rl.LoadVertexBuffer(spriteIndices, false)
	rl.SetVertexAttribute(2, 1, rl.Float, false, 4, 0)
	rl.EnableVertexAttribute(2)
	rl.SetVertexAttributeDivisor(2, 1)

	// Position VBO (dynamic, divisor=1) — location=3, 2 floats (vec2) per instance
	positionData := make([]rl.Vector2, entityCount)
	positionVBO := rl.LoadVertexBuffer(positionData, true)
	rl.SetVertexAttribute(3, 2, rl.Float, false, 8, 0)
	rl.EnableVertexAttribute(3)
	rl.SetVertexAttributeDivisor(3, 1)

	rl.DisableVertexArray()

	return &SpriteRenderer{
		shaderID:    shaderID,
		vao:         vao,
		quadVBO:     quadVBO,
		ebo:         ebo,
		spriteVBO:   spriteVBO,
		positionVBO: positionVBO,
		mvpLoc:      mvpLoc,
		cellSizeLoc: cellSizeLoc,
		textureLoc:  textureLoc,
		texture:     sheet,
		entityCount: entityCount,
		tb:          tb,
	}
}

func (s *SpriteRenderer) Unload() {
	rl.UnloadVertexBuffer(s.quadVBO)
	rl.UnloadVertexBuffer(s.ebo)
	rl.UnloadVertexBuffer(s.spriteVBO)
	rl.UnloadVertexBuffer(s.positionVBO)
	rl.UnloadShaderProgram(s.shaderID)
}

func (s *SpriteRenderer) Render() {
	// Flush raylib's internal batch before we take over GL state
	rl.DrawRenderBatchActive()

	rl.EnableShader(s.shaderID)

	// Compute and set MVP
	mvp := rl.MatrixMultiply(rl.GetMatrixModelview(), rl.GetMatrixProjection())
	rl.SetUniformMatrix(s.mvpLoc, mvp)

	// Bind texture
	rl.ActiveTextureSlot(0)
	rl.EnableTexture(s.texture.ID)

	// Upload new position data if available
	if s.tb.SwapReader() {
		positions := s.tb.ReadBuffer()
		if len(positions) > 0 {
			rl.UpdateVertexBuffer(s.positionVBO, positions, 0)
		}
	}

	// Draw instanced
	rl.EnableVertexArray(s.vao)
	rl.DrawVertexArrayElementsInstanced(0, 6, nil, s.entityCount)
	rl.DisableVertexArray()

	rl.DisableTexture()
	rl.DisableShader()

	// Restore raylib state
	rl.DrawRenderBatchActive()
}
