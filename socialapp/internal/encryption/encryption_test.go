package encryption

import (
	"testing"
)

func TestGenerateKeypair(t *testing.T) {
	privKey, publicKey := GenerateKeyPair(2048)

	privKeyString := string(PrivateKeyToBytes(privKey))
	pubKeyString := string(PublicKeyToBytes(publicKey))

	privKeyParsed := BytesToPrivateKey([]byte(privKeyString))
	pubKeyParsed := BytesToPublicKey([]byte(pubKeyString))

	message := []byte("hello world")
	encryptedMessage := EncryptWithPublicKey(message, pubKeyParsed)

	decryptedMessage := DecryptWithPrivateKey(encryptedMessage, privKeyParsed)
	if string(message) != string(decryptedMessage) {
		t.Error("message and decrypted message are not equal, message: ", string(message), " decrypted message: ", string(decryptedMessage))
	}

	signature, hash := SignWithPrivateKey(message, privKeyParsed)
	if err := VerifyWithPublicKey(signature, hash, pubKeyParsed); err != nil {
		t.Error("message is not correct")
	}
}
