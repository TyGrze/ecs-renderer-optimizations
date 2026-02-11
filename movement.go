package main

import "github.com/mlange-42/ark/ecs"

type MovementSystem struct {
  filter *ecs.Filter2[Position, Velocity]
}

func NewMovementSystem(w *ecs.World) *MovementSystem{
  filter := ecs.NewFilter2[Position, Velocity](w)

  return &MovementSystem{
    filter: filter,
  }
}

func (s *MovementSystem) Update(w *ecs.World) {
  query := s.filter.Query()

  for query.Next() {
    pos, vel := query.Get()

    pos.X += vel.X
    pos.Y += vel.Y
  }
}
