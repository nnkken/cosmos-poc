package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
)

func GetCmdBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command {
		Use: "balance [Becn32 address]",
		Short: "balance address",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			addrStr := args[0]
			res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/balance/%s", queryRoute, addrStr), nil)
			if err != nil {
				fmt.Printf("Cannot get balance from %s\n", addrStr)
				return nil
			}
			coins := sdk.Coins{}
			cdc.MustUnmarshalJSON(res, &coins)
			return cliCtx.PrintOutput(coins)
		},
	}
}