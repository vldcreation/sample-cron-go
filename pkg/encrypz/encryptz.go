package encrypz

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// WFN is a function type to write key pair to file
// @param privateKeyName, publicKeyName
// @return error
type WFN func(keyName string, source []byte) error

type symetricEncryption struct {
	// private key instance of rsa.PrivateKey
	// use to decrypt data
	privateKey *rsa.PrivateKey
	// public key instance of rsa.PublicKey
	// use to encrypt data
	publicKey *rsa.PublicKey
}

type symetricImpl interface {
	encrypt(data []byte) []byte
	decrypt(chiper []byte) []byte
}

func NewSymetricEncryption() *symetricEncryption {
	privKey, pubKey, err := lookupKeyPair("id_rsa.pem", "id_rsa.pub")
	if err != nil {
		log.Fatalf("failed to lookup key pair: %v", err)
	}

	return &symetricEncryption{
		privateKey: privKey,
		publicKey:  pubKey,
	}
}

// generate key pair using rsa library
// @param bits
// given bit size
// @param wfn
// write function to write key pair to file
func genKeyPairWithWrite(bits int, wfn WFN) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	prvKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	prvKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(prvKey),
	})

	err = wfn("id_rsa.pem", prvKeyPem)
	if err != nil {
		return nil, nil, err
	}

	// extract public component from private key
	pub := prvKey.Public()

	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("invalid public key type")
	}

	pubKeyPem := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(pubKey),
	})

	err = wfn("id_rsa.pub", pubKeyPem)
	if err != nil {
		return nil, nil, err
	}

	return prvKey, pubKey, nil
}

// extract key pair from file
// @return privateKey, publicKey
func extractKeyPair(privPath, pubPath string) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	prv, err := ioutil.ReadFile(privPath)
	if err != nil {
		log.Printf("failed to read private key file: %v", err)
		return nil, nil, err
	}

	privPem, _ := pem.Decode(prv)
	if privPem == nil {
		log.Printf("failed to decode private key file: %v\n please use: openssl genrsa -out id_rsa.pem 2048 to generate private key", err)
		log.Println("openssl rsa -in key.pem -outform PEM -pubout -out id_rsa.pub")
		log.Println("Or you can run: bash gen.sh")
		return nil, nil, err
	}

	var privBt []byte

	if privPem.Type == "RSA PRIVATE KEY" {
		privBt = privPem.Bytes
	} else {
		return nil, nil, errors.New("invalid private key type")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privBt)
	if err != nil {
		return nil, nil, err
	}

	pub, err := ioutil.ReadFile(pubPath)
	if err != nil {
		return nil, nil, err
	}

	pubPem, _ := pem.Decode(pub)
	if pubPem == nil {
		return nil, nil, err
	}

	var pubBt []byte

	if pubPem.Type == "PUBLIC KEY" {
		pubBt = pubPem.Bytes
	} else {
		return nil, nil, errors.New("invalid public key type")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(pubBt)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}

// lookup key pair from file
// check if file exists
// if exists, extract the key pair
func lookupKeyPair(privPath, pubPath string) (privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, err error) {
	_, err = statMultipleFiles(privPath, pubPath)
	if err != nil {
		// if one of the files does not exist
		// generate key pair
		return genKeyPairWithWrite(2048, func(keyName string, source []byte) error {
			return os.WriteFile(keyName, source, 0666)
		})
	}

	// if both files exist
	// extract key pair
	return extractKeyPair(privPath, pubPath)
}

func statMultipleFiles(paths ...string) (bool, error) {
	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

func statFile(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func (s *symetricEncryption) encrypt(data []byte) ([]byte, error) {
	// implement encryption
	// i just use rsa library
	// and use as simple as possible
	chiper, err := rsa.EncryptPKCS1v15(rand.Reader, s.publicKey, data)
	if err != nil {
		log.Printf("failed to encrypt data: %v", err)
		return nil, err
	}

	return chiper, nil
}

func (s *symetricEncryption) decrypt(chiper []byte) ([]byte, error) {
	oaep, err := rsa.DecryptPKCS1v15(rand.Reader, s.privateKey, chiper)
	if err != nil {
		log.Printf("failed to decrypt data: %v", err)
		return nil, err
	}
	return oaep, nil
}
