package bmailLib

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	stamp_token "github.com/realbmail/Bmail_token"
	"github.com/realbmail/go-stamp-walllet"
)

var stampWallet stamp.Wallet

func NewStampWallet(auth string) string {
	w, e := stamp.NewWallet(auth)
	if e != nil {
		return ""
	}
	stampWallet = w
	return w.String()
}

func StampDetails(stampAddr string) string {
	if stampWallet == nil {
		fmt.Println("please create stamp wallet first")
		return ""
	}
	details, err := stamp_token.DetailsOfStamp(BlockChainQueryUrl,
		common.HexToAddress(stampAddr),
		stampWallet.Address())

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	byts, err := json.Marshal(details)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(byts)
}
