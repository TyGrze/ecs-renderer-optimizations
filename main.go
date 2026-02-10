package main

import (
	"flag"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

var debugMode bool

const (
	title        = "ECS Test"
	screenWidth  = 1280
	screenHeight = 720
	entityCount  = 200000
	spawnRange   = 20000
)

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug overlay")
	flag.Parse()

  // Window setup
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(screenWidth, screenHeight, title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)

	// ECS Init
	world := ecs.NewWorld()
	mapper := ecs.NewMap2[Position, Sprite](world)
	cameraBounds := CameraBounds{}
	ecs.AddResource(world, &cameraBounds)

	sheet := GenerateSpritesheet()
	defer rl.UnloadTexture(sheet)

	spriteRenderer := NewSpriteRenderer(world, sheet)

	camera := NewCamera()

	var debugOverlay *DebugOverlay
	if debugMode {
		debugOverlay = NewDebugOverlay(world, &camera, &cameraBounds)
	}

	// Create entities
	for range entityCount {
		ranx := rand.Float32() * spawnRange
		rany := rand.Float32() * spawnRange
		_ = mapper.NewEntity(
			&Position{X: ranx, Y: rany},
			&Sprite{X: int32(ranx) % 4, Y: int32(rany) % 4},
		)
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		camera.Update(dt)
		camera.UpdateBounds(&cameraBounds)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)

		rl.BeginMode2D(camera.Cam)
		spriteRenderer.Update(world)
		rl.EndMode2D()

		if debugOverlay != nil {
			debugOverlay.Draw(spriteRenderer)
		}

		rl.DrawFPS(int32(rl.GetScreenWidth())-100, 10)
		rl.EndDrawing()

	}
}
