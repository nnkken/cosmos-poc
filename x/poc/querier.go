package poc

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryBalance = "balance"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown poc query endpoint")
		}
	}
}

func queryBalance(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, sdkErr sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, sdk.ErrInvalidAddress("Invalid Bech32 address")
	}
	coins := keeper.GetBalance(ctx, addr)
	bs, err := codec.MarshalJSONIndent(keeper.cdc, coins)
	if err != nil {
		panic("Caanot marshal coins to JSON")
	}
	return bs, nil
}