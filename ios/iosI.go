package bmailLib

import "github.com/BASChain/go-bmail-lib/resolver"

type UICallBack interface {
	Notification(typ int, msg string)
}

var callback UICallBack
var nameRel resolver.NameResolver

func InitSystem(cb UICallBack, debug bool) {
	callback = cb
	nameRel = resolver.NewEthResolver(debug)
}
