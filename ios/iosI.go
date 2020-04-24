package ios

type UICallBack interface {
}

var callback UICallBack

func InitSystem(cb UICallBack) {
	callback = cb
}
