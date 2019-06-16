package tokenlock

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgLockCoins - struct for locking coins in a timelock
type MsgLockCoins struct {
	Amount     sdk.Coins
	UnlockTime time.Duration
	Owner      sdk.AccAddress
}

func NewMsgLockCoins(amount sdk.Coins, unlockTime time.Duration, owner sdk.AccAddress) MsgLockCoins {
	return MsgLockCoins{
		Amount:     amount,
		UnlockTime: unlockTime,
		Owner:      owner,
	}
}

//nolint
func (msg MsgLockCoins) Route() string { return RouterKey }
func (msg MsgLockCoins) Type() string  { return "lock_coins" }
func (msg MsgLockCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Depositer)}
}

// get the bytes for the message signer to sign on
func (msg MsgLockCoins) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgLockCoins) ValidateBasic() sdk.Error {
	if msg.Depositer.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil depositer address")
	}
	return nil
}

// MsgUnlockCoins - struct for beginning an unlock procedure for locked coins
type MsgUnlockCoins struct {
	Amount     sdk.Coins
	UnlockTime time.Duration
	Owner      sdk.AccAddress
}

func NewMsgUnlockCoins(unlockTime time.Duration, amount sdk.Coins, unlocker sdk.AccAddress) MsgUnlockCoins {
	return MsgUnlockCoins{
		Amount:     amount,
		UnlockTime: unlockTime,
		Owner:      owner,
	}
}

//nolint
func (msg MsgUnlockCoins) Route() string { return RouterKey }
func (msg MsgUnlockCoins) Type() string  { return "unlock_coins" }
func (msg MsgUnlockCoins) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.Unlocker)}
}

// get the bytes for the message signer to sign on
func (msg MsgUnlockCoins) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgUnlockCoins) ValidateBasic() sdk.Error {
	if msg.Unlocker.Empty() {
		return sdk.NewError(DefaultCodespace, CodeInvalidInput, "nil depositer address")
	}
	return nil
}
