package poc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	coinKeeper bank.Keeper
	storeKey sdk.StoreKey
	cdc *codec.Codec
}

func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey: storeKey,
		cdc: cdc,
	}
}

func (keeper Keeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return keeper.coinKeeper.GetCoins(ctx, addr)
}