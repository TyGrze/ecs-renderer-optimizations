package main

import (
	"flag"
	"math/rand/v2"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

var (
	debugMode  bool
	stressMode bool
)

const (
	title        = "ECS Test"
	screenWidth  = 1280
	screenHeight = 720
	entityCount  = 200000
	spawnRange   = 20000
)

func main() {
	flag.BoolVar(&debugMode, "debug", false, "enable debug overlay")
	flag.BoolVar(&stressMode, "stress", false, "burn CPU cycles to test frame timings")
	flag.Parse()

  // Window setup
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(screenWidth, screenHeight, title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)

	// ECS Init
	world := ecs.NewWorld()
	mapper := ecs.NewMap3[Position, Velocity, Sprite](world)
	cameraBounds := CameraBounds{}
	ecs.AddResource(world, &cameraBounds)

	sheet := GenerateSpritesheet()
	defer rl.UnloadTexture(sheet)

  // Setup Rendere System
	spriteRenderer := NewSpriteRenderer(world, sheet)
	defer spriteRenderer.Unload()

  // Setup Movement System
  movementSystem := NewMovementSystem(world)

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
      &Velocity{X: rand.Float32()*2 - 1, Y: rand.Float32()*2 - 1},
			&Sprite{X: int32(ranx) % 4, Y: int32(rany) % 4},
		)
	}

	for !rl.WindowShouldClose() {
		frameStart := time.Now()
		dt := rl.GetFrameTime()

		// Update phase (CPU)
		updateStart := time.Now()
		camera.Update(dt)
		camera.UpdateBounds(&cameraBounds)
    
		if stressMode {
			// Busy-wait ~2ms to make CPU time visible in debug overlay
			burnUntil := time.Now().Add(2 * time.Millisecond)
			for time.Now().Before(burnUntil) {
			}
		}

		updateTime := time.Since(updateStart)

		// Draw phase (CPU draw call submission)
		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)

		drawStart := time.Now()
		rl.BeginMode3D(camera.Cam)
		spriteRenderer.Update(world)
    movementSystem.Update(world)
		rl.EndMode3D()
		drawTime := time.Since(drawStart)

		if debugOverlay != nil {
			debugOverlay.Draw(spriteRenderer)
		}

		rl.DrawFPS(int32(rl.GetScreenWidth())-100, 10)

		// Present phase (GPU flush + vsync)
		presentStart := time.Now()
		rl.EndDrawing()
		presentTime := time.Since(presentStart)

		if debugOverlay != nil {
			debugOverlay.Timings = FrameTimings{
				Update:  float64(updateTime.Microseconds()) / 1000.0,
				Draw:    float64(drawTime.Microseconds()) / 1000.0,
				Present: float64(presentTime.Microseconds()) / 1000.0,
				Total:   float64(time.Since(frameStart).Microseconds()) / 1000.0,
			}
		}
	}
}
