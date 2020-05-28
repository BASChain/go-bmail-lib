package bmailLib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/bmp/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/google/uuid"
)

type EnvelopeOfUI struct {
	Eid      string   `json:"eid"`
	Subject  string   `json:"sub"`
	MsgBody  string   `json:"msg"`
	From     string   `json:"from"`
	FromName string   `json:"fromName"`
	TOs      []string `json:"tos"`
	CCs      []string `json:"ccs"`
	BCCs     []string `json:"bccs"`
	MailType int8     `json:"mType"`
}

var bmClient *client.BMailClient = nil

type MailCallBack interface {
	Process(typ int, msg string)
}

func newClient() (*client.BMailClient, error) {

	if basResolver == nil {
		return nil, fmt.Errorf("no valid bas resolver")
	}

	conf := &client.ClientConf{
		Wallet:   activeWallet,
		Resolver: basResolver,
	}

	bc, err := client.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return bc, nil
}

func CloseClient() {
	if bmClient != nil {
		bmClient.Close()
		bmClient = nil
	}
}

func SendMailJson(mailJson string, cb MailCallBack) bool {

	if activeWallet == nil || !activeWallet.IsOpen() {
		cb.Process(BMErrWalletInvalid, "wallet is nil or locked")
		return false
	}

	if bmClient == nil {
		bc, err := newClient()
		if err != nil {
			fmt.Println(err.Error())
			cb.Process(BMErrClientInvalid, err.Error())
			return false
		}
		bc.Wallet = activeWallet
		bmClient = bc
	}
	jsonMail := &EnvelopeOfUI{}
	if err := json.Unmarshal([]byte(mailJson), jsonMail); err != nil {
		fmt.Println("mail json data is invalid", err)
		uiCallback.Error(BMErrInvalidJson, err.Error())
		return false
	}
	toAddr, _ := basResolver.BMailBCA(jsonMail.TOs[0])
	if !toAddr.IsValid() {
		fmt.Println("can't find receiver's block chain info")
		cb.Process(BMErrNoSuchBas, "can't find receiver's block chain info")
		return false
	}

	env := &bmp.RawEnvelope{
		EnvelopeHead: bmp.EnvelopeHead{
			Eid:      uuid.MustParse(jsonMail.Eid),
			From:     activeWallet.MailAddress(),
			FromAddr: activeWallet.Address(),
			To:       jsonMail.TOs[0],
			ToAddr:   toAddr,
		},
		EnvelopeBody: bmp.EnvelopeBody{
			Subject: jsonMail.Subject,
			MsgBody: jsonMail.MsgBody,
		},
	}

	if jsonMail.MailType == bmp.BMailModeP2P {
		if err := bmClient.SendP2pMail(env); err != nil {
			fmt.Println(err.Error())
			cb.Process(BMErrSendFailed, err.Error())
			return false
		}
	} else {
		if err := bmClient.SendP2sMail(env); err != nil {
			fmt.Println(err.Error())
			cb.Process(BMErrSendFailed, err.Error())
			return false
		}
	}
	cb.Process(BMErrNone, "success")
	return true
}

func BPop(timeSince1970 int64, olderThanSince bool, pieceSize int, cb MailCallBack) []byte {

	if activeWallet == nil || !activeWallet.IsOpen() {
		cb.Process(BMErrWalletInvalid, "wallet is nil or locked")
		return nil
	}
	if bmClient == nil {
		bc, err := newClient()
		if err != nil {
			fmt.Println(err.Error())
			cb.Process(BMErrClientInvalid, err.Error())
			return nil
		}
		bc.Wallet = activeWallet
		bmClient = bc
	}

	envs, err := bmClient.ReceiveEnv(timeSince1970, olderThanSince, pieceSize) //TODO:: seconds to milliseconds
	if err != nil {
		cb.Process(BMErrReceiveFailed, err.Error())
		return nil
	}

	if envs == nil || len(envs) == 0 {
		cb.Process(BMErrNone, "No update data")
		return nil
	}

	byts, err := json.Marshal(envs)
	if err != nil {
		cb.Process(BMErrMarshFailed, err.Error())
		return nil
	}
	//fmt.Println(string(byts))
	cb.Process(BMErrNone, fmt.Sprintf("New Mail got[%d]", len(envs)))
	return byts
}

func GetAddrByName(to string) string {
	toAddr, _ := basResolver.BMailBCA(to)
	if !toAddr.IsValid() {
		uiCallback.Error(BMErrNoSuchBas, "can't find receiver's block chain info")
		return ""
	}

	return toAddr.String()
}

func Encode(data string) string {
	encoded, err := account.Encrypt(activeWallet.Seeds(), []byte(data))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hexutil.Encode(encoded)
}

func Decode(data string) string {
	d, err := hexutil.Decode(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	decoded, err := account.Decrypt(activeWallet.Seeds(), d)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(decoded)
}

func DecodeForPeer(data, fromAddr string) string {
	aesKey, err := activeWallet.AeskeyOf(bmail.Address(fromAddr).ToPubKey())
	if err != nil {
		fmt.Println("DecodeForPeer ===AeskeyOf===>", err)
		return ""
	}

	bb, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		fmt.Println("DecodeForPeer ====DecodeString==>", err)
		return ""
	}
	byts, err := account.Decrypt(aesKey, bb)
	if err != nil {
		fmt.Println("DecodeForPeer ====Decrypt==>", err)
		return ""
	}
	return string(byts)
}
