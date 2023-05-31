package poa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shivam2320/modules/x/poa/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	k.CalculateValidatorVotes(ctx)
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return k.ApplyAndReturnValidatorSetUpdates(ctx)
}
