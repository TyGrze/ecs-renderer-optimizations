package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	panSpeed  = 500.0
	zoomSpeed = 0.1
	zoomMin   = 0.1
	zoomMax   = 5.0
)

type Camera struct {
	Cam rl.Camera2D
}

func NewCamera() Camera {
	return Camera{
		Cam: rl.Camera2D{
			Target: rl.NewVector2(0, 0),
			Zoom:   1.0,
			Offset: rl.NewVector2(float32(rl.GetScreenWidth())/2, float32(rl.GetScreenHeight())/2),
		},
	}
}

func (c *Camera) Update(dt float32) {
	// Camera pan with WASD
	if rl.IsKeyDown(rl.KeyW) {
		c.Cam.Target.Y -= panSpeed * dt / c.Cam.Zoom
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.Cam.Target.Y += panSpeed * dt / c.Cam.Zoom
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.Cam.Target.X -= panSpeed * dt / c.Cam.Zoom
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Cam.Target.X += panSpeed * dt / c.Cam.Zoom
	}

	// Zoom mousewheel
	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		c.Cam.Zoom += wheel * zoomSpeed
		if c.Cam.Zoom < zoomMin {
			c.Cam.Zoom = zoomMin
		}
		if c.Cam.Zoom > zoomMax {
			c.Cam.Zoom = zoomMax
		}
	}
}

func (c *Camera) UpdateBounds(bounds *CameraBounds) {
	topLeft := rl.GetScreenToWorld2D(rl.NewVector2(0, 0), c.Cam)
	botRight := rl.GetScreenToWorld2D(rl.NewVector2(float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())), c.Cam)
	bounds.MinX = topLeft.X
	bounds.MinY = topLeft.Y
	bounds.MaxX = botRight.X
	bounds.MaxY = botRight.Y
}
