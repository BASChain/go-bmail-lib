package bmailLib

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmp"
)

type EnvelopeOfUI struct {
	Eid     string   `json:"eid"`
	Subject string   `json:"sub"`
	MsgBody string   `json:"msg"`
	From    string   `json:"from"`
	TOs     []string `json:"tos"`
	CCs     []string `json:"ccs"`
	BCCs    []string `json:"bccs"`
}

var bmClient *bmp.BMailClient = nil

func NewMailClient() bool {

	if basResolver == nil {
		fmt.Println("no valid bas resolver")
		return false
	}

	if activeWallet == nil {
		fmt.Println("no valid wallet")
		return false
	}

	conf := &bmp.ClientConf{
		Wallet:   activeWallet,
		Resolver: basResolver,
	}

	bc, err := bmp.NewClient(conf)
	if err != nil {
		fmt.Println(err)
		return false
	}

	bmClient = bc
	return true
}

func SendMail(mailJson string) {

}
