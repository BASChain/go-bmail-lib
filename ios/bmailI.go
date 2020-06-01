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
)

type EnvelopeOfUI struct {
	Eid      string   `json:"eid"`
	Subject  string   `json:"subject"`
	MsgBody  string   `json:"mailBody"`
	From     string   `json:"from"`
	FromName string   `json:"fromName"`
	TOs      []string `json:"tos"`
	CCs      []string `json:"ccs"`
	BCCs     []string `json:"bccs"`
	PinCode  []byte   `json:"pin"`
	PreEid	 string   `json:"preEid"`
}

var bmClient *client.BMailClient = nil

type MailCallBack interface {
	Process(typ int, msg string)
}

func fullFillRcpt(mailNames []string, typ int8, pinCode []byte) ([]*bmp.Recipient, error) {
	rcpts := make([]*bmp.Recipient, 0)

	for _, name := range mailNames {

		toAddr, _ := basResolver.BMailBCA(name)
		if !toAddr.IsValid() {
			return nil, fmt.Errorf("can't find rcpt[%s] block chain info", name)
		}

		aesKey, err := activeWallet.AeskeyOf(toAddr.ToPubKey())
		if err != nil {
			return nil, err
		}

		iv := bmp.NewIV()
		encodePin, err := account.EncryptWithIV(aesKey, iv.Bytes(), pinCode)
		if err != nil {
			return nil, err
		}

		rcpt := &bmp.Recipient{
			ToName:   name,
			ToAddr:   toAddr,
			RcptType: typ,
			AESKey:   encodePin,
		}
		rcpts = append(rcpts, rcpt)
	}

	return rcpts, nil
}

func (eui *EnvelopeOfUI) Seal() (*bmp.BMailEnvelope, error) {

	rcpts := make([]*bmp.Recipient, 0)
	if len(eui.TOs) > 0 {
		tos, err := fullFillRcpt(eui.TOs, bmp.RcpTypeTo, eui.PinCode)
		if err != nil {
			return nil, err
		}
		rcpts = append(rcpts, tos...)
	}

	if len(eui.CCs) > 0 {
		ccs, err := fullFillRcpt(eui.CCs, bmp.RcpTypeCC, eui.PinCode)
		if err != nil {
			return nil, err
		}
		rcpts = append(rcpts, ccs...)
	}

	if len(eui.BCCs) > 0 {
		bccs, err := fullFillRcpt(eui.BCCs, bmp.RcpTypeBcc, eui.PinCode)
		if err != nil {
			return nil, err
		}
		rcpts = append(rcpts, bccs...)
	}
	env := &bmp.BMailEnvelope{
		Eid:      eui.Eid,
		From:     eui.FromName,
		FromAddr: bmail.Address(eui.From),
		RCPTs:    rcpts,
		Subject:  eui.Subject,
		MailBody: eui.MsgBody,
	}
	return env, nil
}

func (eui *EnvelopeOfUI) ToString() string {
	return fmt.Sprintf("\n================EnvelopeOfUI================" +
		"\n\tEid:%40s" +
		"\n\tFromName:%40s" +
		"\n\tFrom:%40s" +
		"\n\tPinCode:%x" +
		"\n================================",
		eui.Eid,
		eui.FromName,
		eui.From,
		eui.PinCode)
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

func validate(cb MailCallBack) error {
	if activeWallet == nil || !activeWallet.IsOpen() {
		cb.Process(BMErrWalletInvalid, "wallet is nil or locked")
		return fmt.Errorf("wallet is nil or locked")
	}

	if bmClient == nil {
		bc, err := newClient()
		if err != nil {
			cb.Process(BMErrClientInvalid, err.Error())
			return err
		}
		bc.Wallet = activeWallet
		bmClient = bc
	}
	return nil
}

func SendMailJson(mailJson string, cb MailCallBack) bool {
	if err := validate(cb); err != nil {
		cb.Process(BMErrClientInvalid, err.Error())
		return false
	}
	fmt.Println("======>Before send mail:=>", mailJson)
	jsonMail := &EnvelopeOfUI{}
	if err := json.Unmarshal([]byte(mailJson), jsonMail); err != nil {
		uiCallback.Error(BMErrInvalidJson, err.Error())
		return false
	}
	fmt.Println("======EnvelopeOfUI mail:=>", jsonMail.ToString())

	env, err := jsonMail.Seal()
	if err != nil {
		cb.Process(BMErrClientInvalid, err.Error())
		return false
	}
	fmt.Println("======>BMailEnvelope:=>", env.ToString())

	if err := bmClient.SendMail(env); err != nil {
		fmt.Println(err.Error())
		cb.Process(BMErrSendFailed, err.Error())
		return false
	}
	cb.Process(BMErrNone, "success")
	return true
}

func BPop(timeSince1970 int64, olderThanSince bool, pieceSize int, cb MailCallBack) []byte {

	if err := validate(cb); err != nil {
		cb.Process(BMErrClientInvalid, err.Error())
		return nil
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

func Encode(data string, iv []byte) string {
	encoded, err := account.EncryptWithIV(activeWallet.Seeds(), iv, []byte(data))
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

func PinCode()[]byte{
	return (bmp.NewIV())[:]
}