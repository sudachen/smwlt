package fu

import (
	"crypto/aes"
	"crypto/cipher"
)

func AesXor(key []byte, data []byte) (result []byte, err error) {
	if block, err := aes.NewCipher(key); err == nil {
		iv := [16]byte{}
		iv[15] = 5
		result = make([]byte, len(data))
		cipher.NewCTR(block, iv[:]).XORKeyStream(result, data)
	}
	return
}
