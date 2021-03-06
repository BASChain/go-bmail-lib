package bmailLib

import (
	"fmt"
	"github.com/realbmail/go-bmail-lib/utils"
	"github.com/realbmail/go-bmail-resolver"
)

type UICallBack interface {
	Notification(typ int, msg string)
	Error(typ int, msg string)
}

const (
	BlockChainQueryUrl = "https://ropsten.infura.io/v3/d64d364124684359ace20feae1f9ac20"
	BMErrNone          = iota
	BMErrClientInvalid
	BMErrWalletInvalid
	BMErrNoSuchBas
	BMErrInvalidJson
	BMErrSendFailed
	BMErrReceiveFailed
	BMErrMarshFailed
)

var uiCallback UICallBack
var basResolver resolver.NameResolver

func InitSystem(cb UICallBack, debug bool) {
	uiCallback = cb
	basResolver = resolver.NewEthResolver(debug)
}

func CalculateHash(mailName string) string {
	return resolver.GetHash(mailName).String()
}

func MailBcaByMailName(mailName string) string {
	bca, cname := basResolver.BMailBCA(mailName)
	if !bca.IsValid() {
		return ""
	}
	result := string(bca) + "," + cname
	fmt.Println("mail's block chain data:=>", result)
	return result
}

func CName(mailName string) string {
	_, cname := basResolver.BMailBCA(mailName)
	return cname
}

func MailIcon(mailName string) []byte {
	if mailName == "" {
		return nil
	}
	return utils.GenIDIcon(mailName)
}

func MailID() string {
	return utils.UUID()
}
