package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/styu12/seungohcoin/utils"
)

const (
	fileName string = "seungohcoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address string
}

var w *wallet 

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	privBytes, err := x509.MarshalECPrivateKey(key) 
	utils.HandleError(err)
	err = os.WriteFile(fileName, privBytes, 0644)
	utils.HandleError(err)
}

func restoreKey() *ecdsa.PrivateKey {
	privBytes, err := os.ReadFile(fileName)
	utils.HandleError(err)
	key, err := x509.ParseECPrivateKey(privBytes)
	utils.HandleError(err)
	return key
}

func addressFromK(key *ecdsa.PrivateKey) string {
	z := append(key.X.Bytes(), key.Y.Bytes()...)
	address := fmt.Sprintf("%x", z)
	return address
}

func sign(payload string, w *wallet) string {
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsBytes)
	utils.HandleError(err)
	signature := append(r.Bytes(), s.Bytes()...)
	return fmt.Sprintf("%x", signature)
}



func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		}	else {
			key := createPrivKey()
			w.privateKey = key
			persistKey(key)
		}
		w.Address = addressFromK(w.privateKey)
	}
	return w
}