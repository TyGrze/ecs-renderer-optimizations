package systems

import "github.com/mlange-42/ark/ecs"

type ReadOnlySystem interface{
  Update(w *ecs.World)
}
