package app

import (
	"encoding/json"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/shivam2320/modules/demo-apps/auth-chain/x/authchain"
	authchainkeeper "github.com/shivam2320/modules/demo-apps/auth-chain/x/authchain/keeper"
	authchaintypes "github.com/shivam2320/modules/demo-apps/auth-chain/x/authchain/types"
	"github.com/shivam2320/modules/x/poa"
	poakeeper "github.com/shivam2320/modules/x/poa/keeper"
	poatypes "github.com/shivam2320/modules/x/poa/types"
	// this line is used by starport scaffolding # 1
)

const appName = "authchain"

var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.authchaincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.authchaind")
	ModuleBasics    = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		poa.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		authchain.AppModuleBasic{},
		// this line is used by starport scaffolding # 2
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName: nil,
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	poaKeeper       poakeeper.Keeper
	supplyKeeper    supply.Keeper
	paramsKeeper    params.Keeper
	authchainKeeper authchainkeeper.Keeper
	// this line is used by starport scaffolding # 3
	mm *module.Manager

	sm *module.SimulationManager
}

var _ simapp.App = (*NewApp)(nil)

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey,
		auth.StoreKey,
		poatypes.StoreKey,
		supply.StoreKey,
		params.StoreKey,
		authchaintypes.StoreKey,
		// this line is used by starport scaffolding # 5
	)

	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[poatypes.ModuleName] = app.paramsKeeper.Subspace(poakeeper.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	app.poaKeeper = poakeeper.NewKeeper(
		app.bankKeeper,
		app.cdc,
		keys[authchaintypes.StoreKey],
		app.subspaces[poatypes.ModuleName],
	)

	app.authchainKeeper = authchainkeeper.NewKeeper(
		app.bankKeeper,
		app.cdc,
		keys[authchaintypes.StoreKey],
	)

	// this line is used by starport scaffolding # 4

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.poaKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		authchain.NewAppModule(app.authchainKeeper, app.bankKeeper),
		poa.NewAppModule(app.poaKeeper, app.bankKeeper),
		// this line is used by starport scaffolding # 6
	)

	app.mm.SetOrderEndBlockers(poatypes.ModuleName)

	genutil.ModuleCdc = app.cdc

	app.mm.SetOrderInitGenesis(
		poatypes.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		authchaintypes.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		// this line is used by starport scaffolding # 7
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
