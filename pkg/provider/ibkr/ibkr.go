// Package ibkr implements the Interactive Brokers futures provider.
// IBKR provides access to global futures markets including CME, Eurex,
// HKFE, SGX, and 30+ other exchanges.
package ibkr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/luxfi/futures/pkg/provider"
	"github.com/luxfi/futures/pkg/types"
)

const (
	ProdBaseURL  = "https://localhost:5000/v1/api"
	PaperBaseURL = "https://localhost:5000/v1/api"
)

// Verify interface compliance at compile time.
var _ provider.FuturesProvider = (*Provider)(nil)

// Provider implements the IBKR futures provider via the Client Portal API.
type Provider struct {
	baseURL    string
	httpClient *http.Client
}

// New creates a new IBKR futures provider.
// IBKR uses Client Portal Gateway which runs locally; paper controls the
// gateway mode, not the URL.
func New(paper bool) *Provider {
	base := ProdBaseURL
	if paper {
		base = PaperBaseURL
	}
	return &Provider{
		baseURL:    base,
		httpClient: &http.Client{},
	}
}

func (p *Provider) Name() string { return "ibkr" }

func (p *Provider) GetFuturesContracts(ctx context.Context, underlying string) ([]*types.FuturesContract, error) {
	return nil, fmt.Errorf("ibkr: GetFuturesContracts not yet implemented")
}

func (p *Provider) GetFuturesQuote(ctx context.Context, symbol string) (*types.FuturesQuote, error) {
	return nil, fmt.Errorf("ibkr: GetFuturesQuote not yet implemented")
}

func (p *Provider) CreateFuturesOrder(ctx context.Context, accountID string, req *types.CreateFuturesOrderRequest) (*types.Order, error) {
	return nil, fmt.Errorf("ibkr: CreateFuturesOrder not yet implemented")
}

func (p *Provider) GetFuturesPositions(ctx context.Context, accountID string) ([]*types.FuturesPosition, error) {
	return nil, fmt.Errorf("ibkr: GetFuturesPositions not yet implemented")
}

func (p *Provider) CloseFuturesPosition(ctx context.Context, accountID, symbol string, qty *int) (*types.Order, error) {
	return nil, fmt.Errorf("ibkr: CloseFuturesPosition not yet implemented")
}

func (p *Provider) GetFuturesMargin(ctx context.Context, accountID, symbol string) (*types.FuturesMarginRequirement, error) {
	return nil, fmt.Errorf("ibkr: GetFuturesMargin not yet implemented")
}
