package poc

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgTransfer struct {
	From sdk.AccAddress
	To sdk.AccAddress
	Value sdk.Coins
}

func NewMsgTransfer(from, to sdk.AccAddress, value sdk.Coins) MsgTransfer {
	return MsgTransfer {
		From: from,
		To: to,
		Value: value,
	}
}

func (msg MsgTransfer) Route() string {
	return "poc"
}

func (msg MsgTransfer) Type() string {
	return "transfer"
}

func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress(msg.From.String())
	}
	if msg.To.Empty() {
		return sdk.ErrInvalidAddress(msg.To.String())
	}
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	bs, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bs)
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}