package main

import (
	"crypto/md5"
	"fmt"
	"github.com/BASChain/go-account"
	"github.com/BASChain/go-bmail-account"
	"github.com/BASChain/go-bmail-lib/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"image/color"
	"image/png"
	"io/ioutil"
)

func main() {
	acc_stra := `{"version":1,"address":"BMFnZekMi2f6bujPmnpAu64QvbUXiGsmEMbTzPw8dSMk5n","bmail":"","cipher":"x256WHtmL522CsS2LKoASrAzWyRRSe9KcV4FzSG5q7Wn684SPJHWQEA6qWtdq7kuvJYwVYQMceQC92nTXXweVEumTq89j6WqNd3z6xMY9dX4K"}`
	acc_strb := `{"version":1,"address":"BMEmA22Y6AwU1MWXJL22SDT9yPhfEPKSaro6uVyL9f3TDB","bmail":"","cipher":"2MK4sZjuhU7XRZGJ5Ep4qggrA7nUQPS7anpFSdadTrJwvCaVLHHDqMktMUg8BWucnCDmgzrNRviiHjPs152psWwf9Y84xyRkBu66mKnqeRQtJ"}`

	w_a, e := bmail.LoadWalletByData(acc_stra)
	if e != nil {
		panic(e)
	}
	w_b, e := bmail.LoadWalletByData(acc_strb)
	if e != nil {
		panic(e)
	}
	if e := w_a.Open("BMail"); e != nil {
		panic(e)
	}
	if e := w_b.Open("BMail"); e != nil {
		panic(e)
	}

	acc_a := w_a.Address()
	acc_b := w_b.Address()
	aes_a_b, e := w_a.AeskeyOf(acc_b.ToPubKey())
	if e != nil {
		panic(e)
	}
	aes_b_a, e := w_b.AeskeyOf(acc_a.ToPubKey())
	if e != nil {
		panic(e)
	}
	fmt.Println(hexutil.Encode(aes_a_b))
	fmt.Println(hexutil.Encode(aes_b_a))

	cipher, e := account.Encrypt(aes_a_b, []byte("This is a BMail title"))
	if e != nil {
		panic(e)
	}

	fmt.Println(hexutil.Encode(cipher))
	plain, e := account.Decrypt(aes_b_a, cipher)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(plain))

	cipher, e = account.Encrypt(aes_a_b, []byte("可以放心使用BMail邮箱发送私密邮件，这是一个点对点绝对安全的保密邮箱。"))
	if e != nil {
		panic(e)
	}

	fmt.Println(hexutil.Encode(cipher))
	plain, e = account.Decrypt(aes_b_a, cipher)
	if e != nil {
		panic(e)
	}
	fmt.Println(string(plain))
}

func test2() {
	w, e := bmail.NewWallet("BMail")
	if e != nil {
		panic(e)
	}
	fmt.Println(w.String())
}

func rgb(r, g, b uint8) color.NRGBA { return color.NRGBA{r, g, b, 255} }
func md5hash(s string) []byte {
	h := md5.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}
func test1() {
	var data = md5hash("BM6MqKLq5rBJgHcR6w4p4GXuHSgBuCzxB7LVpRHWP16UTw") // hex.DecodeString("BM6MqKLq5rBJgHcR6w4p4GXuHSgBuCzxB7LVpRHWP16UTw")
	var config = utils.Sigil{
		Rows: 5,
		Foreground: []color.NRGBA{
			rgb(45, 79, 255),
			rgb(44, 172, 0),
			rgb(254, 180, 44),
			rgb(226, 121, 234),
			rgb(30, 179, 253),
			rgb(232, 77, 65),
			rgb(49, 203, 115),
			rgb(141, 69, 170),
			rgb(252, 125, 31),
		},
		Background: rgb(224, 224, 224),
	}

	img := config.Make(420, true, data)
	fil, err := ioutil.TempFile(".", "a.png")
	if err != nil {
		panic(err)
	}

	if err := png.Encode(fil, img); err != nil {
		panic(err)
	}
}
