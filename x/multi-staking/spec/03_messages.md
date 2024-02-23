<!--
order: 3
-->

# Messages

In this section we describe the processing of the multi-staking messages and the corresponding updates to the state. 
All created/modified state objects specified by each message are defined within the [state](./02_state.md) section. We use staking module message types to easier migrate and apply to frontend side.

## MsgCreateValidator

A validator is created using the `MsgCreateValidator` message.
The validator must be created with an initial delegation from the operator. 
The Initial delegation token will be the `multi staking denom` of validator.
This `multi staking denom` can not be changed from validator and a validator will just have one denom.

Logic flow:

1. Lock `MultiStakingCoin` and mint `BondCoin`.

2. Create `MultiStakingLock`.

3. Use minted `BondCoin` amount as the delegation of `stakingtypes.MsgCreateValidator` and
calling `stakingkeeper.CreateValidator()`

This message is expected to fail if:

* `multi staking denom` is not supported.
* The call to `stakingkeeper.CreateValidator()` returns an error.

## MsgEditValidator

The properties expect `multi staking denom` of a validator can be updated using the
`MsgEditValidator` message.

Logic flow:

1. Just calling `stakingkeeper.EditValidator()`.

This message is expected to fail if:

* The call to `stakingkeeper.EditValidator()` returns an error.

## MsgDelegate

Within this message the delegator locked up coins in the `multi-staking` module account. 
The `multi-staking` inturns mint a calculated amount of `BondDenom` and delegate to the validator.

Logic flow:

1. Lock `MultiStakingCoin` and mint `BondCoin`.

2. Create or add `MultiStakingCoin` to `MultiStakingLock`.

3. Use minted `BondCoin` amount as the delegation of `stakingtypes.MsgDelegate` and
calling `stakingkeeper.MsgDelegate()`

This message is expected to fail if:

* `multi staking denom` is mismatched with validator.
* The call to `stakingkeeper.MsgDelegate()` returns an error.

## MsgUndelegate

The `MsgUndelegate` message allows delegators to undelegate their tokens from
validator.

Logic flow:

* Calculate amount of `BondCoin` need to be `sdk undelegated`

* Adjust `BondCoin` to fit with delegation shares

* Remove `MultiStakingCoin` from `MultiStakingLock`

* Create or add `MultiStakingCoin` to `MultiStakingUnlockEntry` at block height

* Call `stakingkeeper.Undelegate()` with the calculated amount of `BondCoin`

The rest of the unbonding logic such as sending locked coins back to user will happens at `EndBlock()`

This message is expected to fail if:

* `multi staking denom` is mismatched with validator.
* Lock record not found.
* The call to `stakingkeeper.Undelegate()` returns an error.

## MsgCancelUnbonding 

The `MsgCancelUnbonding` message allows delegators to cancel the `unbondingDelegation` entry and delegate back to a previous validator.

Logic flow:

* Calculate amount of `BondCoin` need to be `cancel undelegation`

* Decrease `MultiStakingCoin` from `MultiStakingUnlock`

* Adjust `BondCoin` to fit with cancel unbonding amount

* Redelegate cancel amount to validator

* Call `stakingkeeper.CancelUnbondingDelegation()` with the calculated amount of `sdkbond token`

This message is expected to fail if:

* `multi staking denom` is mismatched with validator.
* Unlock entry at request height not found.
* The call to `stakingkeeper.CancelUnbondingDelegation()` returns an error.

## MsgBeginRedelegate

The `MsgBeginRedelegate` message allows delegators to instantly switch validators. Once
the unbonding period has passed, the redelegation is automatically completed in
the EndBlocker.

Logic flow:

* Calculate amount of `BondCoin` need to be `redelegate`

* Decrease `MultiStakingCoin` of delegation at `fromValidator` and add that `MultiStakingCoin` to delegation at `toValidator` then calculate new `BondWeight` of this lock record.
* Call `stakingkeeper.BeginRedelegate()` with the calculated amount of `sdkbond token`

This message is expected to fail if:

* `multi staking denom` is mismatched with validator.
* Lock not found.
* The call to `stakingkeeper.BeginRedelegation()` returns an error.