package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	// "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/shivam2320/modules/demo-apps/auth-chain/x/authchain/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group authchain queries under a subcommand
	authchainQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	authchainQueryCmd.AddCommand(
		flags.GetCommands(
      // this line is used by starport scaffolding # 1
		)...,
	)

	return authchainQueryCmd
}
