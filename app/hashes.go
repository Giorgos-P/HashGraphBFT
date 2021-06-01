package app

import (
	"HashGraphBFT/types"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
)

func makeHash(ev *types.EventMessage) string {
	return SHA384(fmt.Sprintf("%#v", ev))
}

// MD5 - md5 hash
func MD5(text string) string {
	algorithm := md5.New()
	return stringHasher(algorithm, text)

}

// SHA1 hashes using sha1 algorithm
func SHA1(text string) string {
	algorithm := sha1.New()
	return stringHasher(algorithm, text)
}

// SHA256 hashes using sha256 algorithm
func SHA256(text string) string {
	algorithm := sha256.New()
	return stringHasher(algorithm, text)
}

// SHA384 hashes using sha512 algorithm
func SHA384(text string) string {
	algorithm := sha512.New384()
	return stringHasher(algorithm, text)
}

// SHA512 hashes using sha512 algorithm
func SHA512(text string) string {
	algorithm := sha512.New()
	return stringHasher(algorithm, text)
}

func stringHasher(algorithm hash.Hash, text string) string {
	algorithm.Write([]byte(text))
	return hex.EncodeToString(algorithm.Sum(nil))
}

// func test() {

// 	var text string
// 	text = "abc"

// 	//hash := app.MD5(text)
// 	hash := MD5(text)
// 	fmt.Println(text)
// 	fmt.Println(hash)
// }
