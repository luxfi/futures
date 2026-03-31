// Package types defines the domain types for the futures trading module.
package types

import "time"

// Order is a unified trade order returned by all providers.
type Order struct {
	ID             string     `json:"id"`
	Provider       string     `json:"provider"`
	ProviderID     string     `json:"provider_id"`
	AccountID      string     `json:"account_id"`
	Symbol         string     `json:"symbol"`
	Qty            string     `json:"qty"`
	Side           string     `json:"side"`     // buy, sell
	Type           string     `json:"type"`     // market, limit, stop, stop_limit
	TimeInForce    string     `json:"time_in_force"`
	LimitPrice     string     `json:"limit_price,omitempty"`
	StopPrice      string     `json:"stop_price,omitempty"`
	Status         string     `json:"status"`
	FilledQty      string     `json:"filled_qty,omitempty"`
	FilledAvgPrice string     `json:"filled_avg_price,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	FilledAt       *time.Time `json:"filled_at,omitempty"`
	CanceledAt     *time.Time `json:"canceled_at,omitempty"`
}

// FuturesContract is a tradeable futures contract.
type FuturesContract struct {
	Symbol            string  `json:"symbol"`     // e.g. ESM26
	Underlying        string  `json:"underlying"` // e.g. ES
	Name              string  `json:"name"`
	Expiration        string  `json:"expiration"`
	Exchange          string  `json:"exchange"` // CME, NYMEX, CBOT, ICE
	TickSize          float64 `json:"tick_size"`
	PointValue        float64 `json:"point_value"`
	MarginInitial     float64 `json:"margin_initial"`
	MarginMaintenance float64 `json:"margin_maintenance"`
	Last              float64 `json:"last"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Volume            int     `json:"volume"`
	OpenInterest      int     `json:"open_interest"`
	SettlementPrice   float64 `json:"settlement_price"`
	Change            float64 `json:"change"`
	ChangePercent     float64 `json:"change_percent"`
	Tradable          bool    `json:"tradable"`
}

// FuturesQuote is a real-time quote for a futures contract.
type FuturesQuote struct {
	Symbol          string  `json:"symbol"`
	Bid             float64 `json:"bid"`
	Ask             float64 `json:"ask"`
	Last            float64 `json:"last"`
	Volume          int     `json:"volume"`
	OpenInterest    int     `json:"open_interest"`
	SettlementPrice float64 `json:"settlement_price"`
}

// CreateFuturesOrderRequest is the request to place a futures order.
type CreateFuturesOrderRequest struct {
	Symbol      string `json:"symbol"`
	Side        string `json:"side"`       // buy, sell
	Qty         string `json:"qty"`
	OrderType   string `json:"order_type"` // market, limit, stop, stop_limit
	LimitPrice  string `json:"limit_price,omitempty"`
	StopPrice   string `json:"stop_price,omitempty"`
	TimeInForce string `json:"time_in_force"`
}

// FuturesPosition is a held futures position.
type FuturesPosition struct {
	Symbol        string `json:"symbol"`
	Underlying    string `json:"underlying"`
	Side          string `json:"side"` // long, short
	Qty           string `json:"qty"`
	AvgEntryPrice string `json:"avg_entry_price"`
	MarkPrice     string `json:"mark_price"`
	UnrealizedPnL string `json:"unrealized_pnl"`
	MarginUsed    string `json:"margin_used"`
	Expiration    string `json:"expiration"`
}

// FuturesMarginRequirement is the margin requirement for a futures contract.
type FuturesMarginRequirement struct {
	Symbol            string `json:"symbol"`
	InitialMargin     string `json:"initial_margin"`
	MaintenanceMargin string `json:"maintenance_margin"`
	DayTradeMargin    string `json:"day_trade_margin"`
}
