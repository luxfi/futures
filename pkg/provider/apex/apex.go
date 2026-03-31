// Package apex implements the Apex Futures FCM (Futures Commission Merchant) provider.
// Apex Futures is a division of Apex Clearing that provides futures clearing
// for CME, NYMEX, CBOT, and ICE exchanges.
package apex

import (
	"context"
	"fmt"
	"net/http"

	"github.com/luxfi/futures/pkg/provider"
	"github.com/luxfi/futures/pkg/types"
)

const (
	ProdBaseURL    = "https://api.apexfutures.com"
	SandboxBaseURL = "https://api-sandbox.apexfutures.com"
)

// Verify interface compliance at compile time.
var _ provider.FuturesProvider = (*Provider)(nil)

// Provider implements the Apex Futures FCM provider.
type Provider struct {
	baseURL    string
	apiKey     string
	apiSecret  string
	httpClient *http.Client
}

// New creates a new Apex Futures provider.
func New(apiKey, apiSecret string, sandbox bool) *Provider {
	base := ProdBaseURL
	if sandbox {
		base = SandboxBaseURL
	}
	return &Provider{
		baseURL:    base,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		httpClient: &http.Client{},
	}
}

func (p *Provider) Name() string { return "apex" }

func (p *Provider) GetFuturesContracts(ctx context.Context, underlying string) ([]*types.FuturesContract, error) {
	return nil, fmt.Errorf("apex: GetFuturesContracts not yet implemented")
}

func (p *Provider) GetFuturesQuote(ctx context.Context, symbol string) (*types.FuturesQuote, error) {
	return nil, fmt.Errorf("apex: GetFuturesQuote not yet implemented")
}

func (p *Provider) CreateFuturesOrder(ctx context.Context, accountID string, req *types.CreateFuturesOrderRequest) (*types.Order, error) {
	return nil, fmt.Errorf("apex: CreateFuturesOrder not yet implemented")
}

func (p *Provider) GetFuturesPositions(ctx context.Context, accountID string) ([]*types.FuturesPosition, error) {
	return nil, fmt.Errorf("apex: GetFuturesPositions not yet implemented")
}

func (p *Provider) CloseFuturesPosition(ctx context.Context, accountID, symbol string, qty *int) (*types.Order, error) {
	return nil, fmt.Errorf("apex: CloseFuturesPosition not yet implemented")
}

func (p *Provider) GetFuturesMargin(ctx context.Context, accountID, symbol string) (*types.FuturesMarginRequirement, error) {
	return nil, fmt.Errorf("apex: GetFuturesMargin not yet implemented")
}
