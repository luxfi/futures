// Package envconfig registers futures providers based on environment variables.
//
// Set the following env vars to enable providers:
//
//	APEX_FUTURES_API_KEY + APEX_FUTURES_API_SECRET  -> enables Apex Futures FCM
//	APEX_FUTURES_SANDBOX=true                       -> use Apex sandbox
//	IBKR_FUTURES_ENABLED=true                       -> enables IBKR futures
//	IBKR_FUTURES_PAPER=true                         -> use IBKR paper trading
package envconfig

import (
	"os"

	"github.com/luxfi/futures/pkg/provider"
	"github.com/luxfi/futures/pkg/provider/apex"
	"github.com/luxfi/futures/pkg/provider/ibkr"
	"github.com/rs/zerolog/log"
)

// Register inspects environment variables and registers available providers.
func Register(reg *provider.Registry) {
	registerApex(reg)
	registerIBKR(reg)
}

func registerApex(reg *provider.Registry) {
	key := os.Getenv("APEX_FUTURES_API_KEY")
	secret := os.Getenv("APEX_FUTURES_API_SECRET")
	if key == "" || secret == "" {
		return
	}
	sandbox := os.Getenv("APEX_FUTURES_SANDBOX") == "true"
	p := apex.New(key, secret, sandbox)
	reg.Register(p)
	log.Info().Str("provider", p.Name()).Bool("sandbox", sandbox).Msg("registered futures provider")
}

func registerIBKR(reg *provider.Registry) {
	if os.Getenv("IBKR_FUTURES_ENABLED") != "true" {
		return
	}
	paper := os.Getenv("IBKR_FUTURES_PAPER") == "true"
	p := ibkr.New(paper)
	reg.Register(p)
	log.Info().Str("provider", p.Name()).Bool("paper", paper).Msg("registered futures provider")
}
