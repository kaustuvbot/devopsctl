package doctor

import (
	"context"
	"errors"
	"testing"

	"github.com/kaustuvprajapati/devopsctl/internal/reporter"
)

// mockModule implements Module for testing
type mockModule struct {
	name     string
	results  []reporter.CheckResult
	err      error
}

func (m *mockModule) Name() string { return m.name }

func (m *mockModule) Run(ctx context.Context) ([]reporter.CheckResult, error) {
	return m.results, m.err
}

// TestRegistryRegister tests module registration
func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()

	module := &mockModule{name: "test", results: []reporter.CheckResult{}}
	err := registry.Register(module)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if registry.Len() != 1 {
		t.Errorf("Expected 1 module, got %d", registry.Len())
	}
}

// TestRegistryDuplicateRegistration tests duplicate registration
func TestRegistryDuplicateRegistration(t *testing.T) {
	registry := NewRegistry()

	module := &mockModule{name: "test", results: []reporter.CheckResult{}}
	registry.Register(module)

	err := registry.Register(module)
	if err != ErrModuleAlreadyRegistered {
		t.Errorf("Expected ErrModuleAlreadyRegistered, got %v", err)
	}
}

// TestRegistryNilModule tests nil module registration
func TestRegistryNilModule(t *testing.T) {
	registry := NewRegistry()

	err := registry.Register(nil)
	if err != ErrNilModule {
		t.Errorf("Expected ErrNilModule, got %v", err)
	}
}

// TestEngineRunAll tests running all modules
func TestEngineRunAll(t *testing.T) {
	engine := NewEngine()

	// Register modules
	engine.Register(&mockModule{
		name: "module1",
		results: []reporter.CheckResult{
			{CheckName: "check1", Severity: "LOW", ResourceID: "res1"},
		},
	})
	engine.Register(&mockModule{
		name: "module2",
		results: []reporter.CheckResult{
			{CheckName: "check2", Severity: "HIGH", ResourceID: "res2"},
		},
	})

	reports, err := engine.RunAll(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(reports) != 2 {
		t.Errorf("Expected 2 reports, got %d", len(reports))
	}
}

// TestEnginePartialFailure tests partial module failure handling
func TestEnginePartialFailure(t *testing.T) {
	engine := NewEngine()

	engine.Register(&mockModule{
		name: "success",
		results: []reporter.CheckResult{
			{CheckName: "check1", Severity: "LOW"},
		},
	})
	engine.Register(&mockModule{
		name: "failure",
		results: nil,
		err:   errors.New("module error"),
	})

	reports, err := engine.RunAll(context.Background())
	if err == nil {
		t.Error("Expected error due to module failure")
	}

	if len(reports) != 2 {
		t.Errorf("Expected 2 reports, got %d", len(reports))
	}

	// Find the failed module
	var failedReport ModuleReport
	for _, r := range reports {
		if r.Module == "failure" {
			failedReport = r
			break
		}
	}

	if failedReport.Error == "" {
		t.Error("Expected error message in failed module report")
	}
}

// TestEngineNoModules tests running with no modules
func TestEngineNoModules(t *testing.T) {
	engine := NewEngine()

	reports, err := engine.RunAll(context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(reports) != 0 {
		t.Errorf("Expected 0 reports, got %d", len(reports))
	}
}
