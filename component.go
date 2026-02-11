package main

type Position struct {
	X, Y float32
}

type Velocity struct {
  X, Y float32
}

type Sprite struct {
	X, Y int32
}

type CameraBounds struct {
	MinX, MinY, MaxX, MaxY float32
}
