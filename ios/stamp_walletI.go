package bmailLib

import (
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
