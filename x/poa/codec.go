package poa

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/shivam2320/modules/x/poa/msg"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(msg.MsgCreateValidatorPOA{}, "poa/MsgCreateValidatorPOA", nil)
	cdc.RegisterConcrete(msg.MsgVoteValidator{}, "poa/MsgVoteValidator", nil)

}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
