// Package provider defines the FuturesProvider interface and a registry for
// managing multiple futures clearing/brokerage backends.
package provider

import (
	"context"
	"fmt"
	"sync"

	"github.com/luxfi/futures/pkg/types"
)

// FuturesProvider is the interface every futures backend must implement.
type FuturesProvider interface {
	// Name returns the provider identifier (e.g. "apex", "ibkr").
	Name() string

	// GetFuturesContracts returns available contracts for an underlying (e.g. "ES", "CL").
	GetFuturesContracts(ctx context.Context, underlying string) ([]*types.FuturesContract, error)

	// GetFuturesQuote returns a real-time quote for a specific contract.
	GetFuturesQuote(ctx context.Context, symbol string) (*types.FuturesQuote, error)

	// CreateFuturesOrder places a futures order.
	CreateFuturesOrder(ctx context.Context, accountID string, req *types.CreateFuturesOrderRequest) (*types.Order, error)

	// GetFuturesPositions returns all futures positions for an account.
	GetFuturesPositions(ctx context.Context, accountID string) ([]*types.FuturesPosition, error)

	// CloseFuturesPosition closes a specific futures position.
	CloseFuturesPosition(ctx context.Context, accountID, symbol string, qty *int) (*types.Order, error)

	// GetFuturesMargin returns margin requirements for a contract.
	GetFuturesMargin(ctx context.Context, accountID, symbol string) (*types.FuturesMarginRequirement, error)
}

// Registry holds all registered futures providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]FuturesProvider
}

// NewRegistry creates a new empty provider registry.
func NewRegistry() *Registry {
	return &Registry{providers: make(map[string]FuturesProvider)}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p FuturesProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Name()] = p
}

// Get returns a provider by name.
func (r *Registry) Get(name string) (FuturesProvider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("futures provider %q not registered", name)
	}
	return p, nil
}

// List returns all registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for n := range r.providers {
		names = append(names, n)
	}
	return names
}
