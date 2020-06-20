package bmailLib

import (
	"encoding/json"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-protocol/bmp"
	"github.com/BASChain/go-bmail-protocol/bmp/client"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var bmClient *client.BMailClient = nil

type MailCallBack interface {
	Process(typ int, msg string)
}

func fullFillRcpt(rcpts []*bmp.Recipient, pinCode []byte) error {

	for _, uir := range rcpts {
		aesKey, err := activeWallet.AeskeyOf(uir.ToAddr.ToPubKey())
		if err != nil {
			return err
		}

		iv := bmp.NewIV()
		encodePin, err := account.EncryptWithIV(aesKey, iv.Bytes(), pinCode)
		if err != nil {
			return err
		}
		uir.AESKey = encodePin
	}
	return nil
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

func SendMailJson(mailJson string, pinCode []byte, cb MailCallBack) bool {
	if err := validate(cb); err != nil {
		return false
	}
	fmt.Println("======>Before send mail:=>", mailJson)
	jsonMail := &bmp.BMailEnvelope{}
	if err := json.Unmarshal([]byte(mailJson), jsonMail); err != nil {
		uiCallback.Error(BMErrInvalidJson, err.Error())
		return false
	}
	if err := fullFillRcpt(jsonMail.RCPTs, pinCode); err != nil {
		cb.Process(BMErrInvalidJson, err.Error())
		return false
	}
	fmt.Println(jsonMail.ToString())
	if err := bmClient.SendMail(jsonMail); err != nil {
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

	envs, err := bmClient.ReceiveEnv(timeSince1970, olderThanSince, pieceSize)
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

func EncodePin(pinCode []byte) []byte{
	iv := bmp.NewIV()
	encoded, err := account.EncryptWithIV(activeWallet.Seeds(), iv.Bytes(), pinCode)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return encoded
}

func DecodePin(pinCipher []byte) []byte {
	decoded, err := account.Decrypt(activeWallet.Seeds(), pinCipher)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return decoded
}

func DecodePinByPeer(pinCipher []byte, fromAddr string) []byte {
	aesKey, err := activeWallet.AeskeyOf(bmail.Address(fromAddr).ToPubKey())
	if err != nil {
		fmt.Println("DecodeForPeer ===AeskeyOf===>", err)
		return nil
	}
	pinCode, err := account.Decrypt(aesKey, pinCipher)
	if err != nil {
		fmt.Println("DecodeForPeer ====Decrypt==>", err)
		return nil
	}
	return pinCode
}

func EncodeByPin(data string, pinCode []byte) string {
	iv := bmp.NewIV()
	encoded, err := account.EncryptWithIV(pinCode, iv.Bytes(), []byte(data))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hexutil.Encode(encoded)
}

func DecodeByPin(data string, pinCode []byte) string {
	d, err := hexutil.Decode(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	decoded, err := account.Decrypt(pinCode, d)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(decoded)
}

func PinCode() []byte {
	return (bmp.NewIV())[:]
}
