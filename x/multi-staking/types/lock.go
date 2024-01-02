package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMultiStakingLock(lockID *LockID, lockedCoin MultiStakingCoin) MultiStakingLock {
	return MultiStakingLock{
		LockID:     lockID,
		LockedCoin: lockedCoin,
	}
}

func (lock MultiStakingLock) MultiStakingCoin(withAmount sdk.Int) MultiStakingCoin {
	return lock.LockedCoin.WithAmount(withAmount)
}

func (lock *MultiStakingLock) RemoveCoinFromMultiStakingLock(removedCoin MultiStakingCoin) error {
	lockedCoinAfter, err := lock.LockedCoin.SafeSub(removedCoin)
	lock.LockedCoin = lockedCoinAfter
	return err
}

func (lock MultiStakingLock) IsEmpty() bool {
	return lock.LockedCoin.Amount.IsZero()
}

func (multiStakingLock *MultiStakingLock) AddCoinToMultiStakingLock(addedCoin MultiStakingCoin) error {
	lockedCoinAfter, err := multiStakingLock.LockedCoin.SafeAdd(addedCoin)
	multiStakingLock.LockedCoin = lockedCoinAfter
	return err
}

func (m MultiStakingLock) GetBondWeight() sdk.Dec {
	return m.LockedCoin.BondWeight
}

func (multiStakingLock MultiStakingLock) LockedAmountToBondAmount(amount math.Int) sdk.Int {
	return multiStakingLock.LockedCoin.WithAmount(amount).BondAmount()
}

func (fromLock *MultiStakingLock) MoveCoinToLock(toLock *MultiStakingLock, coin MultiStakingCoin) error {
	// remove coin from lock on source val
	err := fromLock.RemoveCoinFromMultiStakingLock(coin)
	if err != nil {
		return err
	}

	// add coin to destination lock
	err = toLock.AddCoinToMultiStakingLock(coin)
	if err != nil {
		return err
	}
	return nil
}
