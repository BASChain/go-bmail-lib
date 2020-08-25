package bmailLib

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	stamp_token "github.com/realbmail/Bmail_token"
	"github.com/realbmail/go-bmail-protocol/bmp"
	"github.com/realbmail/go-stamp-walllet"
	"math/big"
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

func StampWalletIsOpen() bool {
	return stampWallet != nil && stampWallet.PriKey() != nil
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
	return byts
}

func WalletEthBalance(user string) int64 {
	eth, err := stamp_token.EthBalance(BlockChainQueryUrl, user)
	if err != nil {
		return 0
	}

	return eth.Int64()
}

func StampReceipt(domain, sAddr string) []byte {

	if stampWallet == nil || nil == stampWallet.PriKey() {
		fmt.Println("need stamp wallet to valid receipt:=>")
		return nil
	}
	userAddr := stampWallet.Address().String()
	test_data := &bmp.StampTXData{
		UserAddr:  userAddr,
		StampAddr: sAddr,
		Credit:    12,
	}

	jd, _ := json.Marshal(test_data)
	return jd

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
	rsyn := &bmp.StampReceiptSyn{
		StampAddr: sAddr,
		UserAddr:  userAddr,
	}

	if err := conn.SendWithHeader(rsyn); err != nil {
		fmt.Println("send stamp receipt query failed:=>", err)
		return nil
	}

	ack := &bmp.StampTX{}
	if err := conn.ReadWithHeader(ack); err != nil {
		fmt.Println("receive stamp receipt err:=>", err)
		return nil
	}
	if ack.UserAddr != userAddr || ack.StampAddr != sAddr {
		fmt.Println("stamp receipt is not for me")
		return nil
	}

	if false == stamp.VerifyJsonSig(stampWallet.Address(), ack.Sig, ack.StampTXData) {
		fmt.Println("stamp receipt signature validation failed")
		return nil
	}

	j, _ := json.Marshal(ack.StampTXData)
	return j
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
	if err := conn.QueryStampOpts(); err != nil {
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

func ActiveStamp(amount int64, tokenAddr string) bool {
	if amount <= 0 {
		fmt.Println("too small amount")
		return false
	}
	if stampWallet == nil {
		fmt.Println("stamp wallet is empty")
		return false
	}
	priKey := stampWallet.PriKey()
	if priKey == nil {
		fmt.Println("stamp wallet isn't open")
		return false
	}
	tx, err := stamp_token.Active(big.NewInt(amount), BlockChainQueryUrl, tokenAddr, priKey)
	if err != nil {
		fmt.Println("active stamp failed:=>", err)
	}
	fmt.Println(tx.Hash().String())

	return true
}
