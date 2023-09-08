package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"strings"

	"github.com/shibme/slv/core/commons"
	"golang.org/x/crypto/argon2"
)

const (
	pwdToKeyArgon2SaltLength uint8  = 16
	pwdToKeyArgon2TimeCost   uint32 = 3
	pwdToKeyArgon2MemoryCost uint32 = 16 * 1024
	pwdToKeyArgon2Threads    uint8  = 1
	pwdToKeyArgon2KeyLength  uint32 = 32

	passwordProtectedDataAbbrev = "PPD" // PP = Password Protected Data
	passwordProtectedKeyAbbrev  = "PPK" // PP = Password Protected Key
)

func getSymmetricKey(passphrase, salt []byte) []byte {
	return argon2.Key(passphrase, salt, pwdToKeyArgon2TimeCost, pwdToKeyArgon2MemoryCost,
		pwdToKeyArgon2Threads, pwdToKeyArgon2KeyLength)
}

func symmetricEncrypt(plaintext []byte, passphrase []byte) ([]byte, error) {
	salt := make([]byte, pwdToKeyArgon2KeyLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	key := getSymmetricKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	ciphertext = append(nonce, ciphertext...)
	ciphertext = append(salt, ciphertext...)
	return ciphertext, nil
}

func SymmetricEncryptData(data []byte, passphrase string) (string, error) {
	ciphertext, err := symmetricEncrypt(data, []byte(passphrase))
	if err != nil {
		return "", err
	}
	return passwordProtectedDataAbbrev + "_" + commons.Encode(ciphertext), nil
}

func SymmetricEncryptKey(secretKey SecretKey, passphrase string) (string, error) {
	ciphertext, err := symmetricEncrypt(secretKey.toBytes(), []byte(passphrase))
	if err != nil {
		return "", err
	}
	return passwordProtectedKeyAbbrev + "_" + commons.Encode(ciphertext), nil
}

func symmetricDecrypt(ciphertext []byte, passphrase []byte) ([]byte, error) {
	salt, ciphertext := ciphertext[:pwdToKeyArgon2SaltLength], ciphertext[pwdToKeyArgon2SaltLength:]
	key := getSymmetricKey(passphrase, salt)
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
		return nil, ErrSymmetricCiphertextShort
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func symmetricDecryptStr(cipherStr, passphrase, abbrev string) ([]byte, error) {
	sliced := strings.Split(cipherStr, "_")
	if len(sliced) != 2 || sliced[0] != abbrev {
		return nil, ErrInvalidSymmetricCipherdataFormat
	}
	ciphertext := commons.Decode(sliced[1])
	return symmetricDecrypt(ciphertext, []byte(passphrase))
}

func SymmetricDecryptData(cipherStr, passphrase string) ([]byte, error) {
	return symmetricDecryptStr(cipherStr, passphrase, passwordProtectedDataAbbrev)
}

func SymmetricDecryptSecretKey(ciphertext, passphrase string) (*SecretKey, error) {
	decrypted, err := symmetricDecryptStr(ciphertext, passphrase, passwordProtectedKeyAbbrev)
	if err == nil {
		var secretKey SecretKey
		err = secretKey.fromBytes(decrypted)
		if err == nil {
			return &secretKey, nil
		}
	}
	return nil, err
}
