package app

import (
	"encoding/json"

	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tendermint/libs/db"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	poc "github.com/nnkken/cosmos-poc/x/poc"
)

const (
	appName = "cosmos-poc"
)

func MakeCodec() *codec.Codec {
	cdc := codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	poc.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type pocApp struct {
	*bam.BaseApp

	cdc *codec.Codec

	keyMain *sdk.KVStoreKey
	keyAccount *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams *sdk.KVStoreKey
	keyPoc *sdk.KVStoreKey
	tkeyParams *sdk.TransientStoreKey

	accountKeeper auth.AccountKeeper
	bankKeeper bank.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	paramsKeeper params.Keeper
	pocKeeper poc.Keeper
}

func NewPocApp(logger log.Logger, db dbm.DB) *pocApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))
	app := &pocApp {
		BaseApp: bApp,
		cdc: cdc,

		keyMain: sdk.NewKVStoreKey("main"),
		keyAccount: sdk.NewKVStoreKey("account"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyParams: sdk.NewKVStoreKey("params"),
		keyPoc: sdk.NewKVStoreKey("poc"),
		tkeyParams: sdk.NewTransientStoreKey("transient_params"),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)
	app.pocKeeper = poc.NewKeeper(
		app.bankKeeper,
		app.keyPoc,
		app.cdc,
	)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("poc", poc.NewHandler(app.pocKeeper))
	
	app.QueryRouter().
		AddRoute("acc", auth.NewQuerier(app.accountKeeper)).
		AddRoute("poc", poc.NewQuerier(app.pocKeeper))

	app.SetInitChainer(app.initChainer)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyFeeCollection,
		app.keyParams,
		app.keyPoc,
		app.tkeyParams,
	)

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}
	return app
}

type GenesisState struct {
	AuthData auth.GenesisState `json:"auth"`
	BankData bank.GenesisState `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
}

func (app *pocApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes
	var genesisState GenesisState
	
	err := app.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}

	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)

	return abci.ResponseInitChain{}
}

func (app *pocApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	accounts := []*auth.BaseAccount{}
	app.accountKeeper.IterateAccounts(ctx, func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins: acc.GetCoins(),
		}
		accounts = append(accounts, account)
		return false
	})

	genState := GenesisState{
		Accounts: accounts,
		AuthData: auth.DefaultGenesisState(),
		BankData: bank.DefaultGenesisState(),
	}

	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	return appState, validators, err
}