package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"

	"github.com/styu12/seungohcoin/utils"
)

const (
	hashedMessage string = "c33084feaa65adbbbebd0c9bf292a26ffc6dea97b170d88e501ab4865591aafd"
	privateKey string = "307702010104207b7ca24302be7041ab9cf77c5af50d717a0f48eef3a5a119a3787372660cbdfda00a06082a8648ce3d030107a14403420004a1990fc424697b0a65859b21c4c4a617dd638e8ecaff80557a5b7fdb5c280673a02d021525a9a03ea10b5027d324942ef4a2a00b281456f6c60419c9d382688c"
	signature string = "9dd7a643b151b5075eb21f470504f56b7552fbb259e716e84bda3a0a263264f49fb9d4f2f80c0f29db41dac1b245b9a708c95d6da003fd4bc7727da404e1e6e8"
)

func Start() {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	keyAsBytes, err := x509.MarshalECPrivateKey(privateKey)
	fmt.Printf("%x\n\n\n\n", keyAsBytes)

	hashAsBytes, err := hex.DecodeString(hashedMessage)
	utils.HandleError(err)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashAsBytes)
	utils.HandleError(err)
	signature := append(r.Bytes(), s.Bytes()...)
	fmt.Printf("%x\n\n",signature)

	ok := ecdsa.Verify(&privateKey.PublicKey, hashAsBytes, r, s)
	fmt.Println(ok)
}