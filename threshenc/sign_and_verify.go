package threshenc

import (
	"HashGraphBFT/logger"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
