package tokenlock

import (
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/cosmos/cosmos-sdk/store/prefix"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper is the model object for the package tokenlock module
type Keeper struct {
	bankKeeper bank.BaseKeeper

	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey

	// The (unexposed) keys used to access the transient stores from the Context.
	transientStoreKey sdk.StoreKey

	// The codec codec for binary encoding/decoding.
	cdc *codec.Codec

	// Reserved codespace
	codespace sdk.CodespaceType
}

// func (keeper Keeper) QueryOwnerLocks(ctx sdk.Context, owner sdk.AccAddress) amc { }

func (keeper Keeper) GetLockAmount(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration) (amount sdk.Coins) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	bz := store.Get(KeyLock(owner, unlockTime))
	if bz == nil {
		return
	}
	keeper.cdc.MustUnmarshalBinaryBare(bz, &amount)
	return
}

func (keeper Keeper) setLockAmount(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration, amount sdk.Coins) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixLocks)
	bz := keeper.cdc.MustMarshalBinaryBare(newAmount)
	store.Set(KeyLock(owner, unlockTime), bz)
}

func (keeper Keeper) LockCoins(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration, amount sdk.Coins) sdk.Error {
	_, err := keeper.bankKeeper.SubtractCoins(ctx, owner, amount)
	if err != nil {
		return err
	}
	currentAmount := keeper.GetLockAmount(ctx, owner, unlockTime)
	newAmount := currentAmount.Add(amount)
	keeper.setLockAmount(ctx, owner, unlockTime, newAmount)
}

func (keeper Keeper) BeginUnlock(ctx sdk.Context, owner sdk.AccAddress, unlockTime time.Duration, amount sdk.Coins) sdk.Error {
	currentAmount := keeper.QueryOwnerTimeLock(ctx, owner, unlockTime)
	newAmount := currentAmount.Sub(amount)
	if newAmount.IsAnyNegative() {
		return ErrInsufficientCoins()
	}
	keeper.setLockAmount(ctx, owner, unlockTime, newAmount)

	keeper.InsertUnlockQueue(
		TokenUnlock{
			Amount:     amount,
			UnlockTime: ctx.BlockHeader().Time.Add(unlockTime),
			Owner:      owner,
		},
	)
}

func (keeper Keeper) FinishUnlock(ctx sdk.Context, unlock TokenUnlock) sdk.Error {
	if unlock.UnlockTime.After(ctx.BlockHeader().Time) {
		panic("unlocked too soon")
	}

	keeper.bankKeeper.AddCoins(ctx, unlock.Owner, unlock.Amount)
	keeper.setLockAmount(ctx, owner, unlockTime, newAmount)

	keeper.InsertUnlockQueue(
		TokenUnlock{
			Amount:     amount,
			UnlockTime: ctx.BlockHeader().Time.Add(unlockTime),
			Owner:      owner,
		},
	)
}

// Returns an iterator for all the unlocks in the Unlock Queue that expire by endTime
func (keeper Keeper) UnlockQueueIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
	return store.Iterator(PrefixUnlockQueue, sdk.PrefixEndBytes(PrefixUnlockQueueTime(endTime)))
}

// Inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertUnlockQueue(ctx sdk.Context, endTime time.Time, unlock TokenUnlock) {
	store := prefix.NewStore(ctx.KVStore(keeper.storeKey), PrefixUnlockQueue)
	bz := keeper.cdc.MustMarshalBinaryBare(unlock)
	store.Set(KeyUnlock(unlock), bz)
}
