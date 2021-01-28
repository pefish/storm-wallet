package kto

import (
	"crypto/ed25519"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	AddressSize = 47
	KTOPrefix   = "Kto"
)

// Wallet include address and private key
type Wallet struct {
	PrivateKey []byte
	Address    string
}

// GenKeyPairWithSeedAndPath create wif and address by seed and path
func GenKeyPairWithSeedAndPath(seed, path string) (address, wif string, err error) {
	m := md5.New()
	m.Write([]byte(seed + path))
	result := hex.EncodeToString(m.Sum(nil))
	pubKey, privateKey, err := ed25519.GenerateKey(strings.NewReader(result))
	if err != nil {
		return "", "", err
	}
	address = PublicKeyToAddress(pubKey)
	if len(address) != AddressSize {
		return "", "", errors.New(`address size error`)
	}
	return address, PrivateKeyToWif(privateKey), nil
}

// PublicKeyToAddress encode public key to string address
func PublicKeyToAddress(publicKey []byte) string {
	pubStr := Encode(publicKey)
	return KTOPrefix + pubStr
}

// PrivateKeyToWif ecode bytes private key to string wif
func PrivateKeyToWif(privateKey []byte) string {
	return Encode(privateKey)
}
