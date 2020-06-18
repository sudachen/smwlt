package tests

import (
	"encoding/hex"
	"fmt"
	"github.com/spacemeshos/ed25519"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/modern"
	"github.com/tyler-smith/go-bip39"
	"gotest.tools/assert"
	"strings"
	"testing"
)

const accMnemonic = "comic attend alarm duck file please wet note slow spirit elevator inch"
const accMnemoSeed = "a263df1683e5f9dfb93483c69fe1681b3a8b50165616d5b8accb9a8e1262b76434179c7586bba372e90d534bbb965010cdda61f2fb36e4aa382306f1609080ce"

type accInfo struct {
	Pub, Pri, Adr string
	No            int
}

var accounts = []accInfo{
	{"c6d93a7b4e1da3297d5c48057a20c6d7f32e5b074072cf539c1e60cbb8a3723e",
		"c6f1585d2d5ab715a2113e6f4a85de2abdbbd3aed95bb16276483581a70f6a92c6d93a7b4e1da3297d5c48057a20c6d7f32e5b074072cf539c1e60cbb8a3723e",
		"0x7a20C6D7F32e5b074072cf539C1e60Cbb8a3723e",
		0},
	{"e1272f189f67c1bafa0deedafccdb318ba94129b05f1631d5aa07ab1d6f7ee00",
		"c682bc6f38d1c671eb289104451c22dda8bdb6b63c942d4f0f0b98206636d7d9e1272f189f67c1bafa0deedafccdb318ba94129b05f1631d5aa07ab1d6f7ee00",
		"0xfcCdB318ba94129b05f1631D5Aa07aB1d6F7Ee00",
		1,
	},
	{"798ee677684d8167a6052876a4b6b9ed9b4a01c608ebdeafc55d1cac22205efb",
		"d87fb5caefee1eb22737fb1f4c49bc7ce2e94f79a6707dfca0d6de5ca48c2728798ee677684d8167a6052876a4b6b9ed9b4a01c608ebdeafc55d1cac22205efb",
		"0xa4b6B9ed9B4a01C608ebdeAFC55d1cac22205efb",
		2,
	},
}

func Test_Key(t *testing.T) {
	for _, a := range accounts {
		key, _ := hex.DecodeString(a.Pri)
		pub := hex.EncodeToString(wallet.PublicKey(key)[:])
		assert.Equal(t, pub, a.Pub)
		assert.Equal(t, len(key), ed25519.PrivateKeySize)
		assert.Equal(t, strings.ToLower(wallet.Address(key).Hex()), strings.ToLower(a.Adr))
	}
}

func Test_MnemonicToSeed(t *testing.T) {
	bs := bip39.NewSeed(accMnemonic, "")
	assert.Equal(t, accMnemoSeed, hex.EncodeToString(bs))
}

func Test_Mnemonic(t *testing.T) {
	for _, a := range accounts {
		_, key := wallet.GenPair(a.No, accMnemonic, modern.DefaultSalt)
		assert.Equal(t, len(key), ed25519.PrivateKeySize)
		pri := hex.EncodeToString(key)
		assert.Equal(t, pri, a.Pri)
		pub := hex.EncodeToString(wallet.PublicKey(key))
		assert.Equal(t, pub, a.Pub)
		assert.Equal(t, strings.ToLower(wallet.Address(key).Hex()), strings.ToLower(a.Adr))
	}
}

func Test_mmm(t *testing.T) {
	mnem := "comic attend alarm duck file please wet note slow spirit elevator inch"
	bs := bip39.NewSeed(mnem, "")
	fmt.Println(hex.EncodeToString(bs))
	fmt.Println(hex.EncodeToString([]byte(modern.DefaultSalt)))
	_, key := wallet.GenPair(3, mnem, modern.DefaultSalt)
	fmt.Println(hex.EncodeToString(key))
	fmt.Println(wallet.Address(key).Hex())
}
