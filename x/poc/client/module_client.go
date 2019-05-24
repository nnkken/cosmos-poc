package client

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	amino "github.com/tendermint/go-amino"

	poccmd "github.com/nnkken/cosmos-poc/x/poc/client/cli"
)

type ModuleClient struct {
	storeKey string
	cdc      *amino.Codec
}

func NewModuleClient(storeKey string, cdc *amino.Codec) ModuleClient {
	return ModuleClient{
		storeKey: storeKey,
		cdc: cdc,
	}
}

func (mc ModuleClient) GetQueryCmd() *cobra.Command {
	pocQueryCmd := &cobra.Command{
		Use:   "query",
		Short: "Querying commands for the poc",
	}

	pocQueryCmd.AddCommand(client.GetCommands(
		poccmd.GetCmdBalance(mc.storeKey, mc.cdc),
	)...)

	return pocQueryCmd
}

func (mc ModuleClient) GetTxCmd() *cobra.Command {
	pocTxCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transaction commands for the poc",
	}

	pocTxCmd.AddCommand(client.PostCommands(
		poccmd.GetCmdTransfer(mc.cdc),
	)...)

	return pocTxCmd
}