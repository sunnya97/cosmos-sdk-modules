package tokenlock

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Key for getting a the next available proposalID from the store
var (
	KeyDelimiter = []byte(":")

	PrefixLocks   = []byte("locks")
	PrefixUnlockQueue = []byte("unlocks")
)

func KeyLock(owner sdk.AccAddress, unlockTime time.Duration) []byte {
	return []byte(strings.Join(string[]{owner.String(), unlockTime.String()}, KeyDelimiter))
}


// Returns the key for a proposalID in the activeProposalQueue
func PrefixUnlockQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		sdk.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func KeyUnlock(unlock TokenUnlock) []byte {
	return bytes.Join([][]byte{
		sdk.FormatTimeBytes(unlock.UnlockTime),
		unlock.Owner.String(),
	}, KeyDelimiter)
}
