package encryption

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"log/slog"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		slog.Error("failed to generate rsa keypair", "error", err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		slog.Error("failed to marshal public key", "error", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		slog.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			slog.Error("failed to decrypt pem block", "error", err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		slog.Error("failed to parse private key", "error", err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		slog.Info("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			slog.Error("failed to decrypt pem block", "error", err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		slog.Error("failed to parse public key", "error", err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		slog.Error("failed to cast public key")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		slog.Error("failed to encrypt with public key", "error", err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		slog.Error("failed to decrypt with private key", "error", err)
	}
	return plaintext
}

// SignWithPrivateKey decrypts data with private key
func SignWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) (signature []byte, hash []byte) {
	hash = sha512.New().Sum(ciphertext)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA512, hash)
	// DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		slog.Error("failed to sign with private key", "error", err)
	}
	return signature, hash
}

// VerifyWithPublicKey decrypts data with private key
func VerifyWithPublicKey(signature []byte, hash []byte, pub *rsa.PublicKey) error {
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA512, hash, signature)
	if err != nil {
		slog.Error("failed to verify signature", "error", err)
	}
	slog.Info("verified")
	return nil
}
