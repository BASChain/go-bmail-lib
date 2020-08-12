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

func StampDetails(stampAddr string) []byte {
	if stampWallet == nil {
		fmt.Println("please create stamp wallet first")
		return nil
	}

	fmt.Println("--1->", stampAddr, "--2->", stampWallet.Address())
	details, err := stamp_token.DetailsOfStamp(BlockChainQueryUrl,
		common.HexToAddress(stampAddr),
		stampWallet.Address())

	if err != nil {
		fmt.Println("query stamp details err:=>", err.Error())
		return nil
	}

	byts, err := json.Marshal(details)
	if err != nil {
		fmt.Println("json Marshal err:=>", err.Error())
		return nil
	}
	fmt.Println(string(byts))
	return byts
}

func WalletEthBalance(user string) int64 {
	eth, err := stamp_token.EthBalance(BlockChainQueryUrl, user)
	if err != nil {
		return 0
	}

	return eth.Int64()
}

func QueryStampListOf(domain string) []byte {
	test_data := &bmp.StampOptsAck{
		IssuerName: "NBS Team",
		HomePage:   "https://www.baschain.cn/",
		StampAddr: []string{"0xb485616F19542fD68Bff6932C6BDd601a6e4839e",
			"0x0C1c9c063952Cd43AF8b3B527A6CE1815da034B4"},
	}
	testj, _ := json.Marshal(test_data)
	return testj

	ips, _ := basResolver.DomainMX(domain)
	if len(ips) == 0 {
		fmt.Println("no service ip found:=>", domain)
		return nil
	}

	fmt.Println("service ips:=>", ips)
	conn, err := bmp.NewBMConn(ips[0])
	if err != nil {
		fmt.Println("create bmail connection err:=>", err)
		return nil
	}

	defer conn.Close()
	if err := conn.QueryStamp(); err != nil {
		fmt.Println("send stamp list query failed:=>", err)
		return nil
	}

	ack := &bmp.StampOptsAck{}
	if err := conn.ReadWithHeader(ack); err != nil {
		fmt.Println("receive stamp options err:=>", err)
		return nil
	}
	j, _ := ack.GetBytes()
	ret := string(j)
	fmt.Println(ret)
	return j
}
