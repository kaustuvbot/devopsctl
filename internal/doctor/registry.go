package doctor

import (
	"context"

	"github.com/kaustuvbot/devopsctl/internal/reporter"
)

// Module represents an audit/validation module that can be run by the doctor engine.
type Module interface {
	// Name returns the module identifier
	Name() string

	// Run executes the module and returns check results
	Run(ctx context.Context) ([]reporter.CheckResult, error)
}

// Registry holds registered modules.
type Registry struct {
	modules map[string]Module
}

// NewRegistry creates a new module registry.
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register adds a module to the registry.
func (r *Registry) Register(m Module) error {
	if m == nil {
		return ErrNilModule
	}
	name := m.Name()
	if name == "" {
		return ErrEmptyModuleName
	}
	if _, exists := r.modules[name]; exists {
		return ErrModuleAlreadyRegistered
	}
	r.modules[name] = m
	return nil
}

// Get returns a module by name.
func (r *Registry) Get(name string) (Module, bool) {
	m, ok := r.modules[name]
	return m, ok
}

// List returns all registered module names.
func (r *Registry) List() []string {
	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}

// Len returns the number of registered modules.
func (r *Registry) Len() int {
	return len(r.modules)
}

// Errors for registry operations
var (
	ErrNilModule           = &RegistryError{"nil module not allowed"}
	ErrEmptyModuleName     = &RegistryError{"module name cannot be empty"}
	ErrModuleAlreadyRegistered = &RegistryError{"module already registered"}
)

type RegistryError struct {
	msg string
}

func (e *RegistryError) Error() string {
	return e.msg
}
