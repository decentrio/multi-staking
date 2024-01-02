package types

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO: make unit test for this
// this is against cosmos convention, doing this for more performance and less storage
func IntermediaryDelegator(multiStakerAddr sdk.AccAddress) sdk.AccAddress {
	return append(multiStakerAddr, 0x0)
}

func MultiStakerAddress(intermediaryDelegator sdk.AccAddress) sdk.AccAddress {
	return bytes.Clone(intermediaryDelegator[:len(intermediaryDelegator)-1])

}
