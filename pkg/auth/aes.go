package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"sync"
)

/*func main() {
	se, err := AesEncrypt("")
	fmt.Println(se, err)
	sd, err := AesDecrypt(se)
	fmt.Println(sd, err)
}*/

var (
	commonkey = []byte("Muxihackeveryday")
	syncMutex sync.Mutex
)

func SetAesKey(key string) {
	syncMutex.Lock()
	defer syncMutex.Unlock()
	commonkey = []byte(key)
}
func AesEncrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(commonkey)
	if err != nil {
		return "", err
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:],
		[]byte(plaintext))
	return hex.EncodeToString(ciphertext), nil

}

func AesDecrypt(d string) (string, error) {
	ciphertext, err := hex.DecodeString(d)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(commonkey)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	//	fmt.Println(len(ciphertext), len(iv))
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}
