package tests

import (
	"encoding/hex"
	"github.com/spacemeshos/ed25519"
	"github.com/sudachen/smwlt/wallet"
	"github.com/sudachen/smwlt/wallet/modern"
	"github.com/tyler-smith/go-bip39"
	"gotest.tools/assert"
	"strings"
	"testing"
)

const accMnemonic = "honey typical frozen dad grief oval glove flame predict again steak manage"
var accMnemoSeed = "781c85780b5e852094e2813cdbe6ca1b7e314ab83ff249a863136b1b3345fee698bcd831283671a25a4676d4d6a6c46c2d2f676fcb598584001ac924c94a8369"
var accPublic = []string{
	"a6eb79e8744c444886725187932d2069cbf2044caa4b1084ed97c1d76601d94f",
}
var accPrivate = []string{
	"781c85780b5e852094e2813cdbe6ca1b7e314ab83ff249a863136b1b3345fee6a6eb79e8744c444886725187932d2069cbf2044caa4b1084ed97c1d76601d94f",
}
var accAddress = []string{
	"0x932d2069cbf2044caa4b1084ed97c1d76601d94f",
}

func Test_Key(t *testing.T) {
	key,_ := hex.DecodeString(accPrivate[0])
	pub := hex.EncodeToString(wallet.PublicKey(key)[:])
	assert.Equal(t,pub,accPublic[0])
	assert.Equal(t,len(key),ed25519.PrivateKeySize)
	assert.Equal(t,strings.ToLower(wallet.Address(key).Hex()),accAddress[0])
}

func Test_MnemonicToSeed(t *testing.T) {
	bs := bip39.NewSeed(accMnemonic,"")
	assert.Equal(t, accMnemoSeed,hex.EncodeToString(bs))
}

func Test_Mnemonic(t *testing.T) {
	for i := range accPrivate {
		_, key := wallet.GenPair(0, accMnemonic, modern.DefaultSalt)
		assert.Equal(t,len(key),ed25519.PrivateKeySize)
		pri := hex.EncodeToString(key)
		assert.Equal(t,pri,accPrivate[i])
		pub := hex.EncodeToString(wallet.PublicKey(key))
		assert.Equal(t,pub,accPublic[i])
		assert.Equal(t,strings.ToLower(wallet.Address(key).Hex()),accAddress[i])
	}
}