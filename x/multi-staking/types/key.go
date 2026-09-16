package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "multistaking"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	BondWeightKey = []byte{0x00}

	ValidatorMultiStakingCoinKey = []byte{0x01}

	MultiStakingLockPrefix = []byte{0x02}

	MultiStakingUnlockPrefix = []byte{0x11} // key for an unbonding-delegation

	ParamsKey = []byte{0x03} // prefix for parameters for module x/multistaking
)

func KeyPrefix(key string) []byte {
	return []byte(key)
}

// GetBondWeightKeyKey returns a key for an index containing the bond coin weight
func GetBondWeightKey(tokenDenom string) []byte {
	return append(BondWeightKey, []byte(tokenDenom)...)
}

// GetValidatorMultiStakingCoinKey returns a key for an index containing the bond denom of a validator
func GetValidatorMultiStakingCoinKey(valAddr sdk.ValAddress) []byte {
	return append(ValidatorMultiStakingCoinKey, []byte(valAddr)...)
}

func MultiStakingLockID(multiStakerAddr string, valAddr string) LockID {
	return LockID{MultiStakerAddr: multiStakerAddr, ValAddr: valAddr}
}

func MultiStakingUnlockID(multiStakerAddr string, valAddr string) UnlockID {
	return UnlockID{MultiStakerAddr: multiStakerAddr, ValAddr: valAddr}
}

// DelAddrAndValAddrFromLockID decodes a complete store key, including its prefix.
func DelAddrAndValAddrFromLockID(lockID []byte) (sdk.AccAddress, sdk.ValAddress, error) {
	return addressesFromID(lockID, MultiStakingLockPrefix[0])
}

// DelAddrAndValAddrFromUnlockID decodes a complete store key, including its prefix.
func DelAddrAndValAddrFromUnlockID(unlockID []byte) (sdk.AccAddress, sdk.ValAddress, error) {
	return addressesFromID(unlockID, MultiStakingUnlockPrefix[0])
}

func addressesFromID(id []byte, prefix byte) (sdk.AccAddress, sdk.ValAddress, error) {
	if len(id) < 2 || id[0] != prefix {
		return nil, nil, fmt.Errorf("invalid staking key header")
	}
	// Convert to int before adding: a byte length can overflow at 254 or 255.
	end := 2 + int(id[1])
	if id[1] == 0 || end >= len(id) {
		return nil, nil, fmt.Errorf("invalid staking key address lengths")
	}
	if err := sdk.VerifyAddressFormat(id[2:end]); err != nil {
		return nil, nil, fmt.Errorf("invalid delegator address: %w", err)
	}
	if err := sdk.VerifyAddressFormat(id[end:]); err != nil {
		return nil, nil, fmt.Errorf("invalid validator address: %w", err)
	}
	return sdk.AccAddress(id[2:end]), sdk.ValAddress(id[end:]), nil
}

func addressLengthPrefix(length int) byte {
	if length > 1<<8-1 {
		panic("address length exceeds one-byte key encoding")
	}

	return byte(length)
}

// // GetUBDKey creates the key for an unbonding delegation by delegator and validator addr
// // VALUE: multi-staking/MultiStakingUnlock
// func GetUBDKey(multiStakerAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
// 	return append(GetUBDsKey(delAddr.Bytes()), address.MustLengthPrefix(valAddr)...)
// }

func (l LockID) ToBytes() []byte {
	multiStakerAddr, valAcc, err := AccAddrAndValAddrFromStrings(l.MultiStakerAddr, l.ValAddr)
	if err != nil {
		panic(err)
	}

	lenMultiStakerAddr := len(multiStakerAddr)

	DVPair := make([]byte, 1+lenMultiStakerAddr+len(valAcc))

	DVPair[0] = addressLengthPrefix(lenMultiStakerAddr)

	copy(DVPair[1:], multiStakerAddr[:])

	copy(DVPair[1+lenMultiStakerAddr:], valAcc[:])

	return append(MultiStakingLockPrefix, DVPair...)
}

func (l UnlockID) ToBytes() []byte {
	multiStakerAddr, valAcc, err := AccAddrAndValAddrFromStrings(l.MultiStakerAddr, l.ValAddr)
	if err != nil {
		panic(err)
	}

	lenMultiStakerAddr := len(multiStakerAddr)

	DVPair := make([]byte, 1+lenMultiStakerAddr+len(valAcc))

	DVPair[0] = addressLengthPrefix(lenMultiStakerAddr)

	copy(DVPair[1:], multiStakerAddr[:])

	copy(DVPair[1+lenMultiStakerAddr:], valAcc[:])

	return append(MultiStakingUnlockPrefix, DVPair...)
}
