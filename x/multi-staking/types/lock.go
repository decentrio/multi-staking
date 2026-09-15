package types

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMultiStakingLock(lockID LockID, lockedCoin MultiStakingCoin) MultiStakingLock {
	return MultiStakingLock{
		LockID:     lockID,
		LockedCoin: lockedCoin,
	}
}

func (lock MultiStakingLock) Validate() error {
	if _, err := sdk.AccAddressFromBech32(lock.LockID.MultiStakerAddr); err != nil {
		return err
	}
	if _, err := sdk.ValAddressFromBech32(lock.LockID.ValAddr); err != nil {
		return err
	}
	return lock.LockedCoin.Validate()
}

func (lock MultiStakingLock) MultiStakingCoin(withAmount math.Int) MultiStakingCoin {
	return lock.LockedCoin.WithAmount(withAmount)
}

func (lock *MultiStakingLock) RemoveCoinFromMultiStakingLock(removedCoin MultiStakingCoin) error {
	lockedCoinAfter, err := lock.LockedCoin.SafeSub(removedCoin)
	if err != nil {
		return err
	}
	lock.LockedCoin = lockedCoinAfter
	return nil
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedCoin.Amount.IsZero()
}

func (multiStakingLock *MultiStakingLock) AddCoinToMultiStakingLock(addedCoin MultiStakingCoin) error {
	lockedCoinAfter, err := multiStakingLock.LockedCoin.SafeAdd(addedCoin)
	if err != nil {
		return err
	}
	multiStakingLock.LockedCoin = lockedCoinAfter
	return nil
}

func (m MultiStakingLock) GetBondWeight() math.LegacyDec {
	return m.LockedCoin.BondWeight
}

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(amount math.Int) math.Int {
	return multiStakingLock.LockedCoin.WithAmount(amount).BondValue()
}

func (fromLock *MultiStakingLock) MoveCoinToLock(toLock *MultiStakingLock, coin MultiStakingCoin) error {
	// Apply both operations to copies so an error leaves both locks unchanged.
	from := *fromLock
	if err := from.RemoveCoinFromMultiStakingLock(coin); err != nil {
		return err
	}
	if fromLock == toLock {
		return nil
	}
	to := *toLock
	if err := to.AddCoinToMultiStakingLock(coin); err != nil {
		return err
	}
	*fromLock = from
	*toLock = to
	return nil
}
