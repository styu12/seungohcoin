package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
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

func restoreBigInts(hexString string) (*big.Int, *big.Int, error) {
	hexAsBytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := hexAsBytes[:len(hexAsBytes)/2]
	secondHalfBytes := hexAsBytes[len(hexAsBytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleError(err)
	x, y, err := restoreBigInts(address)
	utils.HandleError(err)
	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleError(err)
	publicKey := ecdsa.PublicKey {
		Curve: elliptic.P256(),
		X: x,
		Y: y,
	}
	ok := ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
	return ok
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