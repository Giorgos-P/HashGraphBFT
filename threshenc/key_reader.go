package threshenc

import (
	"HashGraphBFT/variables"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"strconv"
)

var (
	// SecretKey - The server's secret key
	SecretKey *rsa.PrivateKey

	// VerificationKeys - All servers verification keys
	VerificationKeys map[int]*rsa.PublicKey
)

// ReadKeys - Reads the keys from local files
func ReadKeys(folder string) {
	secretFile := folder + "secret_" + strconv.Itoa(variables.ID) + ".pem"
	sKey, err := readKeyFromFile(secretFile)
	if err != nil {
		log.Fatal(err)
	}

	SecretKey, err = parseRSAPrivateKeyFromPEM(sKey)
	if err != nil {
		log.Fatal(err)
	}

	VerificationKeys = make(map[int]*rsa.PublicKey, variables.N)
	for i := 0; i < variables.N; i++ {
		verificationFile := folder + "verification_" + strconv.Itoa(i) + ".key"
		vKey, err := readKeyFromFile(verificationFile)
		if err != nil {
			log.Fatal(err)
		}

		VerificationKeys[i], err = parseRSAPublicKeyFromPEM(vKey)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// parseRSAPrivateKeyFromPEM - Parses a rsa.PrivateKey from PEM
func parseRSAPrivateKeyFromPEM(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("Failed to parse PEM block containing the private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// parseRSAPublicKeyFromPEM - Parses a rsa.PublicKey from PEM
func parseRSAPublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("Failed to parse PEM block containing the public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Key type is not rsa.PublicKey")
	}
}

// readKeyFromFile - Reads the key bytes from a file and returns it as a string
func readKeyFromFile(file string) (string, error) {
	key, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(key), nil
}
