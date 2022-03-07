package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/styu12/seungohcoin/utils"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	
	message := "I love you"

	hash := utils.Hash(message)

	hashAsBytes, err := hex.DecodeString(hash)
	utils.HandleError(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleError(err)
	fmt.Printf("R: %d\nS: %d\n",r,s)
}