package tokenlock

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Defines tokenlock module constants
const (
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
)

// Tokenlock stores data about a tokenlock
type TokenUnlock struct {
	Amount     sdk.Coins
	UnlockTime time.Time
	Owner      sdk.AccAddress
}
