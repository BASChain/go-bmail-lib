package bmailLib

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/google/uuid"
)

type EnvelopeOfUI struct {
	Eid      string   `json:"eid"`
	Subject  string   `json:"sub"`
	MsgBody  string   `json:"msg"`
	From     string   `json:"from"`
	TOs      []string `json:"tos"`
	CCs      []string `json:"ccs"`
	BCCs     []string `json:"bccs"`
	MailType int8     `json:"mType"`
}

var bmClient *bmp.BMailClient = nil


type MailSendCallBack interface {
	Success(iv[]byte, enSub, enMsg string)
	Error(typ int, msg string)
}

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

func SendMailJson(mailJson string) bool {

	fmt.Println(mailJson)
	jsonMail := &EnvelopeOfUI{}
	if err := json.Unmarshal([]byte(mailJson), jsonMail); err != nil {
		uiCallback.Error(BMErrInvalidJson, "mail json data is invalid")
		return false
	}
	fmt.Println(jsonMail)
	return true
}

func GetAddrByName(to string) string{
	toAddr, _ := basResolver.BMailBCA(to)
	if !toAddr.IsValid() {
		uiCallback.Error(BMErrNoSuchBas, "can't find receiver's block chain info")
		return ""
	}

	return toAddr.String()
}

func Encode(iv, key []byte, data []byte) []byte{
	encoded, err := account.EncryptWithIV(key, iv, data)
	if err != nil {
		return nil
	}
	return encoded
}

func SendCryptMail(eid, to, sub, msg string, cb MailSendCallBack) bool {

	if bmClient == nil {
		cb.Error(BMErrClientInvalid, "mail client is invalid")
		return false
	}

	if activeWallet == nil || !activeWallet.IsOpen() {
		cb.Error(BMErrWalletInvalid, "wallet is nil or locked")
		return false
	}

	toAddr, _ := basResolver.BMailBCA(to)
	if !toAddr.IsValid() {
		cb.Error(BMErrNoSuchBas, "can't find receiver's block chain info")
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
		cb.Error(BMErrSendFailed, err.Error())
		return false
	}

	return false
}
