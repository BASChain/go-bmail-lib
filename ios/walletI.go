package ios

import "github.com/BASChain/go-bmail-account"

var activeWallet *bmail.Wallet

func NewWallet(auth string) string {
	w, e := bmail.NewWallet(auth)
	if e != nil {
		return ""
	}

	return w.String()
}
