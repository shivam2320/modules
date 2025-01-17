package types

import (
	"github.com/shivam2320/modules/x/poa/msg"
)

// This is a work around to allow messages to be defined in another package
// while allowing the hander to function as expected
type (
	MsgCreateValidatorPOA = msg.MsgCreateValidatorPOA
	MsgVoteValidator      = msg.MsgVoteValidator
)

var (
	NewMsgCreateValidatorPOA = msg.NewMsgCreateValidatorPOA
	NewMsgVoteValidator      = msg.NewMsgVoteValidator
)
