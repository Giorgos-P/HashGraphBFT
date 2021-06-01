package threshenc

import (
	"HashGraphBFT/logger"
	"HashGraphBFT/types"

	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

// SignMessage - Signs the message with the secret key
func SignMessage(message []byte) []byte {
	hash := sha256.New()
	_, err := hash.Write(message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	hashSum := hash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, SecretKey, crypto.SHA256, hashSum, nil)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	return signature
}

// VerifyMessage - Verifies if the message is signed by the correct public key
func VerifyMessage(message []byte, signature []byte, i int) bool {
	hash := sha256.New()
	_, err := hash.Write(message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	err = rsa.VerifyPSS(VerificationKeys[i], crypto.SHA256, hash.Sum(nil), signature, nil)
	if err != nil {
		logger.ErrLogger.Println(err)
		return false
	}
	return true
}

// VerifySignature -
func VerifySignature(ev *types.EventMessage) bool {
	//sign := ev.Signature
	k := *ev
	k.Signature = []byte("0")
	msg := fmt.Sprintf("%#v", k)
	return VerifyMessage([]byte(msg), ev.Signature, ev.Owner)
}

// CreateSignature -
func CreateSignature(ev *types.EventMessage) {

	//sign := ev.Signature
	k := *ev
	k.Signature = []byte("0")
	msg := fmt.Sprintf("%#v", k)
	ev.Signature = SignMessage([]byte(msg))

}
