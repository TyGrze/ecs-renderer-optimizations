package main

import (
	"ecs_test/components"
	"ecs_test/systems"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

const (
	title        = "ECS Test"
	screenWidth  = 1280
	screenHeight = 720
	entityCount  = 100000
	spawnRange   = 10000

	panSpeed  = 500.0
	zoomSpeed = 0.1
	zoomMin   = 0.1
	zoomMax   = 5.0
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowHighdpi)
	rl.InitWindow(screenWidth, screenHeight, title)
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)

	// ECS Init
	world := ecs.NewWorld()
	mapper := ecs.NewMap2[components.Position, components.Sprite](world)

	sheet := GenerateSpritesheet()
	defer rl.UnloadTexture(sheet)

	spriteRenderer := systems.NewSpriteRenderer(world, sheet)

	// Camera
	camera := rl.Camera2D{
		Target: rl.NewVector2(0, 0),
		Zoom:   1.0,
		Offset: rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2),
	}


	// Create entities
	for range entityCount {
		ranx := rand.Float64() * spawnRange
		rany := rand.Float64() * spawnRange
		_ = mapper.NewEntity(
			&components.Position{X: ranx, Y: rany},
			&components.Sprite{X: int32(ranx) % 4, Y: int32(rany) % 4},
		)
	}

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		// Camera pan with WASD
		if rl.IsKeyDown(rl.KeyW) {
			camera.Target.Y -= panSpeed * dt / camera.Zoom
		}
		if rl.IsKeyDown(rl.KeyS) {
			camera.Target.Y += panSpeed * dt / camera.Zoom
		}
		if rl.IsKeyDown(rl.KeyA) {
			camera.Target.X -= panSpeed * dt / camera.Zoom
		}
		if rl.IsKeyDown(rl.KeyD) {
			camera.Target.X += panSpeed * dt / camera.Zoom
		}
		// Zoom mousewhell
		wheel := rl.GetMouseWheelMove()
		if wheel != 0 {
			camera.Zoom += wheel * zoomSpeed
			if camera.Zoom < zoomMin {
				camera.Zoom = zoomMin
			}
			if camera.Zoom > zoomMax {
				camera.Zoom = zoomMax
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Gray)

		rl.BeginMode2D(camera)
		spriteRenderer.Update(world)
		rl.EndMode2D()

		rl.DrawFPS(int32(rl.GetScreenWidth())-100, 10)
		rl.EndDrawing()

	}
}
