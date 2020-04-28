package bmailLib

import "github.com/BASChain/go-bmail-lib/resolver"

type UICallBack interface {
	Notification(typ int, msg string)
}

var callback UICallBack
var basResolver resolver.NameResolver

func InitSystem(cb UICallBack, debug bool) {
	callback = cb
	basResolver = resolver.NewEthResolver(debug)
}

func CalculateHash(mailName string) string {
	return resolver.BMailNameHash(mailName)
}

func MailBcaByHash(mailHash string) string {
	bca, _ := basResolver.BMailBCA(mailHash)
	return string(bca)
}
