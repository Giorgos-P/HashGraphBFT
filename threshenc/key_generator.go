package threshenc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"strconv"
)

// GenerateKeys - Generates the key pair and saves them locally
func GenerateKeys(N int, folder string) {
	for i := 0; i < N; i++ {
		privateKey, publicKey, err := generateKeyPair(2048)
		if err != nil {
			log.Fatal(err)
		}

		privateKeyBytes := exportRSAPrivateKeyAsPEM(privateKey)
		publicKeyBytes, err := exportRSAPublicKeyAsPEM(publicKey)
		if err != nil {
			log.Fatal(err)
		}

		secretFile := folder + "secret_" + strconv.Itoa(i) + ".pem"
		err = writeKeyToFile(privateKeyBytes, secretFile)
		if err != nil {
			log.Fatal(err)
		}

		verificationFile := folder + "verification_" + strconv.Itoa(i) + ".key"
		err = writeKeyToFile(publicKeyBytes, verificationFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// generateKeyPair - Creates a RSA Key Pair of specified byte size
func generateKeyPair(bitSize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Validate Key
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// exportRSAPrivateKeyAsPEM - Exports a rsa.PrivateKey as PEM
func exportRSAPrivateKeyAsPEM(privateKey *rsa.PrivateKey) []byte {
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		},
	)
	return privateKeyPEM
}

// exportRSAPublicKeyAsPEM - Exports a rsa.PublicKey as PEM
func exportRSAPublicKeyAsPEM(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	return publicKeyPEM, nil
}

// writeKeyToFile - Writes the key bytes to a file
func writeKeyToFile(keyBytes []byte, file string) error {
	err := ioutil.WriteFile(file, keyBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
