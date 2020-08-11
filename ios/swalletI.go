package bmailLib

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	stamp_token "github.com/realbmail/Bmail_token"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-stamp-walllet"
)

var stampWallet stamp.Wallet

func NewStampWallet(auth string) string {
	w, e := stamp.NewWallet(auth)
	if e != nil {
		return ""
	}
	stampWallet = w
	return w.String()
}

func StampWalletAddress() string {
	if stampWallet == nil {
		return ""
	}
	return stampWallet.Address().Hex()
}

func OpenStampWallet(auth string) bool {
	if stampWallet == nil {
		return false
	}
	return stampWallet.Open(auth)
}

func StampWalletFromJson(jsonStr string) bool {
	w, e := stamp.WalletOfJson(jsonStr)
	if e != nil {
		return false
	}
	stampWallet = w
	return true
}

func StampDetails(stampAddr string) string {
	if stampWallet == nil {
		fmt.Println("please create stamp wallet first")
		return ""
	}
	details, err := stamp_token.DetailsOfStamp(BlockChainQueryUrl,
		common.HexToAddress(stampAddr),
		stampWallet.Address())

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	byts, err := json.Marshal(details)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	return string(byts)
}

func WalletEthBalance(user string) int64 {
	eth, err := stamp_token.EthBalance(BlockChainQueryUrl, user)
	if err != nil {
		return 0
	}

	return eth.Int64()
}

func QueryStampListOf(domain string) string {
	ips, _ := basResolver.DomainMX(domain)
	if len(ips) == 0 {
		return ""
	}

	conn, err := bmp.NewBMConn(ips[0])
	if err != nil {
		fmt.Println(err)
		return ""
	}

	defer conn.Close()
	if err := conn.QueryStamp(); err != nil {
		fmt.Println(err)
		return ""
	}

	ack := &bmp.StampOptsAck{}
	if err := conn.ReadWithHeader(ack); err != nil {
		fmt.Println(err)
		return ""
	}

	j, _ := json.Marshal(ack)
	return string(j)
}
