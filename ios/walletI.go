package bmailLib

import (
	"fmt"
	"github.com/BASChain/go-bmail-account"
)

var activeWallet bmail.Wallet

func NewWallet(auth string) string {
	w, e := bmail.NewWallet(auth)
	if e != nil {
		return ""
	}

	return w.String()
}

func LoadWalletByPath(path string) string {
	w, err := bmail.LoadWallet(path)
	if err != nil {
		fmt.Println("======>[LoadWalletByPath]: LoadWalletByData err:", err.Error())
		return ""
	}
	activeWallet = w
	fmt.Println("======>[LoadWalletByPath]: Load wallet success:", w.Address().String())
	return w.Address().String()
}

func LoadWallet(jsonStr string) bool {
	w, err := bmail.LoadWalletByData(jsonStr)
	if err != nil {
		fmt.Println("======>[LoadWallet]: LoadWalletByData err:", err.Error())
		return false
	}
	activeWallet = w
	fmt.Println("======>[LoadWallet]: Load wallet success:", w.Address().String())
	return true
}

func OpenWallet(auth string) bool {
	if activeWallet == nil {
		fmt.Println("======>[OpenWallet]: Current wallet instance is nil")
		return false
	}
	err := activeWallet.Open(auth)
	if err != nil {
		fmt.Println("======>[OpenWallet]: open action err:", err.Error(), auth)
		return false
	}
	return true
}

func WalletIsOpen() bool {
	if activeWallet == nil {
		return false
	}
	return activeWallet.IsOpen()
}

func CloseWallet() {
	if activeWallet == nil {
		return
	}
	activeWallet.Close()
	fmt.Println("======>Wallet is closing......")
}

func Address() string {
	if activeWallet == nil {
		return ""
	}
	return activeWallet.Address().String()
}

func WalletJson() string {
	if activeWallet == nil {
		return ""
	}
	return activeWallet.String()
}

func MailName() string {
	if activeWallet == nil{
		return ""
	}
	return activeWallet.MailAddress()
}