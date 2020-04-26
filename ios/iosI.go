package bmailLib

type UICallBack interface {
	Notification(typ int, msg string)
}

var callback UICallBack

func InitSystem(cb UICallBack) {
	callback = cb
}
