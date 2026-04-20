package envchain

// Chain applies a sequence of named transformation stages to a secret map.
// Each stage is a function that receives the current map and returns a new one.
// Stages are applied in order; if any stage returns an error the chain halts.

type Stage struct {
	Name string
	Fn   func(map[string]string) (map[string]string, error)
}

type Chain struct {
	stages []Stage
}

type Result struct {
	Stage  string
	Before map[string]string
	After  map[string]string
	Err    error
}

// New creates an empty Chain.
func New() *Chain {
	return &Chain{}
}

// Add appends a stage to the chain.
func (c *Chain) Add(name string, fn func(map[string]string) (map[string]string, error)) *Chain {
	c.stages = append(c.stages, Stage{Name: name, Fn: fn})
	return c
}

// Run executes all stages in order, returning per-stage results.
// Execution stops at the first error.
func (c *Chain) Run(input map[string]string) ([]Result, map[string]string, error) {
	current := copyMap(input)
	results := make([]Result, 0, len(c.stages))

	for _, stage := range c.stages {
		before := copyMap(current)
		next, err := stage.Fn(current)
		r := Result{Stage: stage.Name, Before: before, After: next, Err: err}
		results = append(results, r)
		if err != nil {
			return results, current, err
		}
		current = next
	}
	return results, current, nil
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
