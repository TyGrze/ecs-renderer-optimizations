package main

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	panSpeed  = 500.0
	zoomSpeed = 0.1
	zoomMin   = 0.1
	zoomMax   = 5.0
)

type Camera struct {
	Cam  rl.Camera3D
	Zoom float32
}

func NewCamera() Camera {
	sw := float32(rl.GetScreenWidth())
	sh := float32(rl.GetScreenHeight())
	return Camera{
		Cam: rl.Camera3D{
			Position:   rl.NewVector3(sw/2, 10, sh/2),
			Target:     rl.NewVector3(sw/2, 0, sh/2),
			Up:         rl.NewVector3(0, 0, -1),
			Fovy:       sh,
			Projection: rl.CameraOrthographic,
		},
		Zoom: 1.0,
	}
}

func (c *Camera) Update(dt float32) {
	speed := panSpeed * dt / c.Zoom

	if rl.IsKeyDown(rl.KeyW) {
		c.Cam.Position.Z -= speed
		c.Cam.Target.Z -= speed
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.Cam.Position.Z += speed
		c.Cam.Target.Z += speed
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.Cam.Position.X -= speed
		c.Cam.Target.X -= speed
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Cam.Position.X += speed
		c.Cam.Target.X += speed
	}

	wheel := rl.GetMouseWheelMove()
	if wheel != 0 {
		c.Zoom += wheel * zoomSpeed
		if c.Zoom < zoomMin {
			c.Zoom = zoomMin
		}
		if c.Zoom > zoomMax {
			c.Zoom = zoomMax
		}
	}

	c.Cam.Fovy = float32(rl.GetScreenHeight()) / c.Zoom
}

func (c *Camera) UpdateBounds(bounds *CameraBounds) {
	halfH := c.Cam.Fovy / 2
	aspect := float32(rl.GetScreenWidth()) / float32(rl.GetScreenHeight())
	halfW := halfH * aspect

	cx := c.Cam.Target.X
	cz := c.Cam.Target.Z

	bounds.MinX = cx - halfW
	bounds.MaxX = cx + halfW
	bounds.MinY = cz - halfH
	bounds.MaxY = cz + halfH
}
