package generator

import "github.com/ttacon/autumn/lib/engine"

// A Plan looks like:
//too

type Generator interface {
	CreatePlan(model []engine.ModelTarget)
}
