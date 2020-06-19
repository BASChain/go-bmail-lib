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

type RcptOfUI struct {
	MailName string `json:"name"`
	MailAddr string `json:"address"`
	Typ      int8   `json:"type"`
}

func (uir *RcptOfUI) String() string {
	return fmt.Sprintf("[name:(%s) address:(%s) type:(%d)]",
		uir.MailName, uir.MailAddr, uir.Typ)
}

type EnvelopeOfUI struct {
	Eid       string      `json:"eid"`
	Subject   string      `json:"subject"`
	MsgBody   string      `json:"mailBody"`
	FromAddr  string      `json:"fromAddr"`
	FromName  string      `json:"from"`
	RCPTs     []*RcptOfUI `json:"rcpts"`
	PinCode   []byte      `json:"pin"`
	SessionID string      `json:"sessionID"`
}

var bmClient *client.BMailClient = nil

type MailCallBack interface {
	Process(typ int, msg string)
}

func fullFillRcpt(uircpts []*RcptOfUI, pinCode []byte) ([]*bmp.Recipient, error) {
	rcpts := make([]*bmp.Recipient, 0)

	for _, uir := range uircpts {
		toAddr := bmail.Address(uir.MailAddr)
		if !toAddr.IsValid() {
			return nil, fmt.Errorf("invalid peer address[%s]", toAddr)
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
			ToName:   uir.MailName,
			ToAddr:   toAddr,
			RcptType: uir.Typ,
			AESKey:   encodePin,
		}
		rcpts = append(rcpts, rcpt)
	}

	return rcpts, nil
}

func (eui *EnvelopeOfUI) Seal() (*bmp.BMailEnvelope, error) {

	rcpts, err := fullFillRcpt(eui.RCPTs, eui.PinCode)
	if err != nil {
		return nil, err
	}
	env := &bmp.BMailEnvelope{
		Eid:      eui.Eid,
		FromName: eui.FromName,
		FromAddr: bmail.Address(eui.FromAddr),
		RCPTs:    rcpts,
		Subject:  eui.Subject,
		MailBody: eui.MsgBody,
	}
	return env, nil
}

func (eui *EnvelopeOfUI) ToString() string {
	str := fmt.Sprintf(
		"\n======================EnvelopeOfUI========================="+
			"\n\tEid:\t%20s"+
			"\n\tFromName:\t%20s"+
			"\n\tFromAddr:\t%20s"+
			"\n\tPinCode:\t%20x"+
			"\n\tSessoinID:\t%20x"+
			"\n\tRCPTs:\t%20d"+
			"\n\tSubject:\t%20s"+
			"\n\tMsgBody:\t%20s"+
			"\n===========================================================",
		eui.Eid,
		eui.FromName,
		eui.FromAddr,
		eui.PinCode,
		eui.SessionID,
		len(eui.RCPTs),
		eui.Subject,
		eui.MsgBody)
	for _, uir := range eui.RCPTs {
		str += fmt.Sprintf("\n%s\n", uir.String())
	}
	return str
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
		return false
	}
	fmt.Println("======>Before send mail:=>", mailJson)
	jsonMail := &EnvelopeOfUI{}
	if err := json.Unmarshal([]byte(mailJson), jsonMail); err != nil {
		uiCallback.Error(BMErrInvalidJson, err.Error())
		return false
	}
	fmt.Println(jsonMail.ToString())

	env, err := jsonMail.Seal()
	if err != nil {
		cb.Process(BMErrClientInvalid, err.Error())
		return false
	}
	fmt.Println(env.ToString())

	if err := bmClient.SendMail(env); err != nil {
		fmt.Println("======>SendMail failed:", err.Error())
		cb.Process(BMErrSendFailed, err.Error())
		return false
	}
	cb.Process(BMErrNone, "success")
	return true
}

func BPop(timeSince1970 int64, olderThanSince bool, pieceSize int, cb MailCallBack) []byte {

	if err := validate(cb); err != nil {
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

func PinCode() []byte {
	return (bmp.NewIV())[:]
}
