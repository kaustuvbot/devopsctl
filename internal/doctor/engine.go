package doctor

import (
	"context"
	"fmt"

	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// ModuleReport holds results from a single module execution.
type ModuleReport struct {
	Module  string                  `json:"module"`
	Results []reporter.CheckResult  `json:"results"`
	Error   string                  `json:"error,omitempty"`
}

// Engine orchestrates all registered modules.
type Engine struct {
	registry *Registry
}

// NewEngine creates a new doctor engine.
func NewEngine() *Engine {
	return &Engine{
		registry: NewRegistry(),
	}
}

// Register adds a module to the engine.
func (e *Engine) Register(m Module) error {
	return e.registry.Register(m)
}

// RunAll executes all registered modules and returns aggregated results.
func (e *Engine) RunAll(ctx context.Context) ([]ModuleReport, error) {
	var reports []ModuleReport
	moduleNames := e.registry.List()

	for _, name := range moduleNames {
		module, ok := e.registry.Get(name)
		if !ok {
			continue
		}

		report := ModuleReport{Module: name}

		results, err := module.Run(ctx)
		if err != nil {
			report.Error = err.Error()
			// Continue running other modules despite failure
		}
		report.Results = results
		reports = append(reports, report)
	}

	// Check if any module failed
	var errs []string
	for _, r := range reports {
		if r.Error != "" {
			errs = append(errs, fmt.Sprintf("%s: %s", r.Module, r.Error))
		}
	}

	var err error
	if len(errs) > 0 {
		err = fmt.Errorf("some modules failed: %v", errs)
	}

	return reports, err
}

// Registry returns the underlying registry.
func (e *Engine) Registry() *Registry {
	return e.registry
}
