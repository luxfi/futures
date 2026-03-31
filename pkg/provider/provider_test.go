package provider_test

import (
	"testing"

	"github.com/luxfi/futures/pkg/provider"
	"github.com/luxfi/futures/pkg/provider/apex"
	"github.com/luxfi/futures/pkg/provider/ibkr"
)

func TestRegistryRegisterAndGet(t *testing.T) {
	reg := provider.NewRegistry()

	a := apex.New("key", "secret", true)
	reg.Register(a)

	got, err := reg.Get("apex")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if got.Name() != "apex" {
		t.Fatalf("expected name apex, got %s", got.Name())
	}
}

func TestRegistryGetUnknown(t *testing.T) {
	reg := provider.NewRegistry()
	_, err := reg.Get("unknown")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestRegistryList(t *testing.T) {
	reg := provider.NewRegistry()
	reg.Register(apex.New("k", "s", true))
	reg.Register(ibkr.New(true))

	names := reg.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(names))
	}

	found := map[string]bool{}
	for _, n := range names {
		found[n] = true
	}
	if !found["apex"] {
		t.Error("missing apex provider")
	}
	if !found["ibkr"] {
		t.Error("missing ibkr provider")
	}
}

// Compile-time interface checks.
var (
	_ provider.FuturesProvider = (*apex.Provider)(nil)
	_ provider.FuturesProvider = (*ibkr.Provider)(nil)
)
