package main

import (
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/mlange-42/ark/ecs"
)

const refreshInterval = time.Second / 3

type FrameTimings struct {
	Update  float64 // CPU logic (camera, bounds) in ms
	Draw    float64 // CPU draw call submission in ms
	Present float64 // EndDrawing (GPU flush + vsync) in ms
	Total   float64 // full frame wall time in ms
}

type debugLine struct {
	text  string
	color rl.Color
}

type DebugOverlay struct {
	world   *ecs.World
	camera  *Camera
	bounds  *CameraBounds
	Timings FrameTimings

	dpi         int32
	lastRefresh time.Time
	cached      []debugLine
}

func NewDebugOverlay(world *ecs.World, camera *Camera, bounds *CameraBounds) *DebugOverlay {
	dpi := int32(rl.GetWindowScaleDPI().X)
	if dpi < 1 {
		dpi = 1
	}
	return &DebugOverlay{
		world:  world,
		camera: camera,
		bounds: bounds,
		dpi:    dpi,
	}
}

func (d *DebugOverlay) refresh(sr *SpriteRendererSystem) {
	yellow := rl.NewColor(255, 215, 0, 255)  // #FFD700
	green := rl.NewColor(0, 204, 122, 255)    // #00CC7A
	white := rl.NewColor(191, 191, 191, 255)  // #BFBFBF

	stats := d.world.Stats()
	t := &d.Timings

	d.cached = []debugLine{
		{"--- Frame Timing ---", yellow},
		{fmt.Sprintf("Total:   %6.2f ms", t.Total), green},
		{fmt.Sprintf("Update:  %6.2f ms (CPU)", t.Update), green},
		{fmt.Sprintf("Draw:    %6.2f ms (CPU)", t.Draw), green},
		{fmt.Sprintf("Present: %6.2f ms (GPU+vsync)", t.Present), green},

		{"--- Screen ---", yellow},
		{fmt.Sprintf("Screen: %dx%d", rl.GetScreenWidth(), rl.GetScreenHeight()), green},
		{fmt.Sprintf("Render: %dx%d", rl.GetRenderWidth(), rl.GetRenderHeight()), green},

		{"--- Camera ---", yellow},
		{fmt.Sprintf("Target X: %.1f", d.camera.Cam.Target.X), green},
		{fmt.Sprintf("Target Y: %.1f", d.camera.Cam.Target.Z), green},
		{fmt.Sprintf("Zoom: %.2f", d.camera.Zoom), green},
		{fmt.Sprintf("Bounds: (%.0f,%.0f)-(%.0f,%.0f)", d.bounds.MinX, d.bounds.MinY, d.bounds.MaxX, d.bounds.MaxY), green},

		{"--- Entities ---", yellow},
		{fmt.Sprintf("Total Alive: %d", stats.Entities.Used), green},
		{fmt.Sprintf("Rendered:    %d", sr.RenderedCount), green},

		{"--- ECS World ---", yellow},
		{fmt.Sprintf("Alive:     %d", stats.Entities.Used), white},
		{fmt.Sprintf("Recycled:  %d", stats.Entities.Recycled), white},
		{fmt.Sprintf("Total:     %d", stats.Entities.Total), white},
		{fmt.Sprintf("Capacity:  %d", stats.Entities.Capacity), white},
		{fmt.Sprintf("Archetypes:      %d", len(stats.Archetypes)), white},
		{fmt.Sprintf("Component Types: %d", len(stats.ComponentTypeNames)), white},
		{fmt.Sprintf("Memory Used:     %d kB", stats.MemoryUsed/1024), white},
		{fmt.Sprintf("Memory Reserved: %d kB", stats.Memory/1024), white},
		{fmt.Sprintf("Cached Filters:  %d", stats.CachedFilters), white},
		{fmt.Sprintf("Locked: %v", stats.Locked), white},
	}
}

func (d *DebugOverlay) Draw(sr *SpriteRendererSystem) {
	now := time.Now()
	if d.cached == nil || now.Sub(d.lastRefresh) >= refreshInterval {
		d.refresh(sr)
		d.lastRefresh = now
	}

	fontSize := int32(16) * d.dpi
	lineHeight := int32(20) * d.dpi
	paddingX := int32(10) * d.dpi
	paddingY := int32(10) * d.dpi
	panelWidth := int32(280) * d.dpi

	panelBg := rl.NewColor(18, 18, 18, 204) // #121212 ~80% opacity

	panelHeight := int32(len(d.cached))*lineHeight + paddingY*2

	rl.DrawRectangle(0, 0, panelWidth, panelHeight, panelBg)

	for i, l := range d.cached {
		y := paddingY + int32(i)*lineHeight
		rl.DrawText(l.text, paddingX, y, fontSize, l.color)
	}
}
