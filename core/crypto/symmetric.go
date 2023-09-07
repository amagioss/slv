package crypto

// import (
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"fmt"
// 	"io"
// 	"strings"

// 	"github.com/shibme/slv/core/commons"
// 	"golang.org/x/crypto/argon2"
// )

// const (
// 	pwdToKeyArgon2SaltLength uint8  = 16
// 	pwdToKeyArgon2TimeCost   uint32 = 3
// 	pwdToKeyArgon2MemoryCost uint32 = 16 * 1024
// 	pwdToKeyArgon2Threads    uint8  = 1
// 	pwdToKeyArgon2KeyLength  uint32 = 32
// )

// func getSymmetricKey(passphrase, salt []byte) []byte {
// 	return argon2.Key(passphrase, salt, pwdToKeyArgon2TimeCost, pwdToKeyArgon2MemoryCost,
// 		pwdToKeyArgon2Threads, pwdToKeyArgon2KeyLength)
// }

// func symmetricEncrypt(plaintext []byte, passphrase []byte) ([]byte, error) {
// 	salt := make([]byte, pwdToKeyArgon2KeyLength)
// 	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
// 		return nil, err
// 	}
// 	key := getSymmetricKey(passphrase, salt)
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	aesgcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nonceSize := aesgcm.NonceSize()
// 	nonce := make([]byte, nonceSize)
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return nil, err
// 	}
// 	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
// 	ciphertext = append(nonce, ciphertext...)
// 	ciphertext = append(salt, ciphertext...)
// 	return ciphertext, nil
// }

// func SymmetricEncrypt(data []byte, passphrase string) (string, error) {
// 	ciphertext, err := symmetricEncrypt(data, []byte(passphrase))
// 	if err != nil {
// 		return "", err
// 	}
// 	return passwordProtectedAbbrev + commons.Encode(ciphertext), nil
// }

// func SymmetricEncryptPrivateKey(privateKey SecretKey, passphrase string) (string, error) {
// 	return SymmetricEncrypt(privateKey.toBytes(), passphrase)
// }

// func symmetricDecrypt(ciphertext []byte, passphrase []byte) ([]byte, error) {
// 	salt, ciphertext := ciphertext[:pwdToKeyArgon2SaltLength], ciphertext[pwdToKeyArgon2SaltLength:]
// 	key := getSymmetricKey(passphrase, salt)
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	aesgcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return nil, err
// 	}
// 	nonceSize := aesgcm.NonceSize()
// 	if len(ciphertext) < nonceSize {
// 		return nil, fmt.Errorf("ciphertext is too short")
// 	}
// 	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
// 	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return plaintext, nil
// }

// func SymmetricDecrypt(ciphertext, passphrase string) ([]byte, error) {
// 	if !strings.HasPrefix(ciphertext, passwordProtectedAbbrev) {
// 		return nil, ErrInvalidPassphraseEncryptedData
// 	}
// 	encodedCipherText := strings.TrimPrefix(ciphertext, passwordProtectedAbbrev)
// 	cipherBytes := commons.Decode(encodedCipherText)
// 	return symmetricDecrypt(cipherBytes, []byte(passphrase))
// }

// func SymmetricDecryptPrivateKey(ciphertext, passphrase string) (*SecretKey, error) {
// 	plainBytes, err := SymmetricDecrypt(ciphertext, passphrase)
// 	if err != nil {
// 		return nil, err
// 	}
// 	key, err := keyFromBytes(plainBytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &SecretKey{
// 		key: key,
// 	}, nil
// }
