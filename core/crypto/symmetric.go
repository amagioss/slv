package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/argon2"
)

func symmetricEncrypt(plaintext []byte, passphrase []byte) ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key := argon2.Key(passphrase, salt, 7, 32*1024, uint8(runtime.NumCPU()), 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(nonce, ciphertext...)
	ciphertext = append(salt, ciphertext...)
	return ciphertext, nil
}

func SymmetricEncrypt(data []byte, passphrase string) (string, error) {
	ciphertext, err := symmetricEncrypt(data, []byte(passphrase))
	if err != nil {
		return "", err
	}
	return passphraseProtectedPrefix + base58.Encode(ciphertext), nil
}

func SymmetricEncryptPrivateKey(privateKey PrivateKey, passphrase string) (string, error) {
	return SymmetricEncrypt(privateKey.toBytes(), passphrase)
}

func symmetricDecrypt(ciphertext []byte, passphrase []byte) ([]byte, error) {
	salt, ciphertext := ciphertext[:16], ciphertext[16:]
	key := argon2.Key(passphrase, salt, 7, 32*1024, uint8(runtime.NumCPU()), 32)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func SymmetricDecrypt(ciphertext, passphrase string) ([]byte, error) {
	if !strings.HasPrefix(ciphertext, passphraseProtectedPrefix) {
		return nil, ErrInvalidPassphraseEncryptedData
	}
	encodedCipherText := strings.TrimPrefix(ciphertext, passphraseProtectedPrefix)
	cipherBytes := base58.Decode(encodedCipherText)
	return symmetricDecrypt(cipherBytes, []byte(passphrase))
}

func SymmetricDecryptPrivateKey(ciphertext, passphrase string) (*PrivateKey, error) {
	plainBytes, err := SymmetricDecrypt(ciphertext, passphrase)
	if err != nil {
		return nil, err
	}
	key, err := keyFromBytes(plainBytes)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		key: key,
	}, nil
}
