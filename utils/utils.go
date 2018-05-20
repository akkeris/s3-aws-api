package utils

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io"
	"net/http"
	"os"
)

var Region string
var key []byte

func Init() string {

	vaulttoken := os.Getenv("VAULT_TOKEN")
	vaultaddr := os.Getenv("VAULT_ADDR")
	s3secret := os.Getenv("S3_SECRET")
	vaultaddruri := vaultaddr + "/v1/" + s3secret
	vreq, err := http.NewRequest("GET", vaultaddruri, nil)
	vreq.Header.Add("X-Vault-Token", vaulttoken)
	vclient := &http.Client{}
	vresp, err := vclient.Do(vreq)
	if err != nil {
		fmt.Println(err)
	}
	defer vresp.Body.Close()
	bodyj, err := simplejson.NewFromReader(vresp.Body)
	if err != nil {
		fmt.Println(err)
	}
	keystring, _ := bodyj.Get("data").Get("key").String()
	key = []byte(keystring)

	Region = os.Getenv("REGION")
	toreturn, _ := bodyj.Get("data").Get("brokerdb").String()
	return toreturn
}

func Encrypt(plaintext string) string {

	text := []byte(plaintext)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(crand.Reader, iv); err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], text)
	return base64.StdEncoding.EncodeToString(ciphertext)
}

func Decrypt(b64 string) string {
	text, _ := base64.StdEncoding.DecodeString(b64)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if len(text) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	return string(text)
}
