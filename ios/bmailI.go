package bmailLib

import (
	"fmt"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/google/uuid"
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

func SendCryptMail(eid, to, sub, msg string) bool {

	if bmClient == nil {
		uiCallback.Error(BMErrClientInvalid, "mail client is invalid")
		return false
	}

	if activeWallet == nil || !activeWallet.IsOpen() {
		uiCallback.Error(BMErrWalletInvalid, "wallet is nil or locked")
		return false
	}

	toAddr, _ := basResolver.BMailBCA(to)
	if !toAddr.IsValid() {
		uiCallback.Error(BMErrNoSuchBas, "can't find receiver's block chain info")
		return false
	}

	env := &bmp.RawEnvelope{
		EnvelopeHead: bmp.EnvelopeHead{
			Eid:      uuid.MustParse(eid),
			From:     activeWallet.MailAddress(),
			FromAddr: activeWallet.Address(),
			To:       to,
		},
		EnvelopeBody: bmp.EnvelopeBody{
			Subject: sub,
			MsgBody: msg,
		},
	}

	if err := bmClient.SendP2pMail(env); err != nil {
		uiCallback.Error(BMErrSendFailed, err.Error())
		return false
	}
	return false
}
