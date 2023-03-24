package cipher

import (
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/nacl/box"
)

type Encrypter struct {
	public_key        [32]byte
	ephemeral_keypair KeyPair
	shared_key        [32]byte
	nonce             [24]byte
}

func (e *Encrypter) New(slv_public_key string) (err error) {
	var pub key
	pub.init(slv_public_key, pub_key_prefix)
	e.public_key = pub.key
	err = e.ephemeral_keypair.generateLocally()
	if err != nil {
		return err
	}
	box.Precompute(&e.shared_key, &e.public_key, &e.ephemeral_keypair.priv.key)
	return
}

func (e *Encrypter) encrypt(data []byte) (ciphered_data []byte) {
	return box.SealAfterPrecomputation(nil, data, &e.nonce, &e.shared_key)
}

func (e *Encrypter) EncryptData(data []byte) (ciphertext string) {
	encrypted := e.encrypt(data)
	hash := checksum(data)
	return fmt.Sprintf("%s%s_%s_%s", secret_prefix, hash, base58.Encode(e.ephemeral_keypair.pub.key[:]), base58.Encode(encrypted))
}

type Decrypter struct {
	private_key key
	nonce       [24]byte
}

func (d *Decrypter) New(slv_private_key string) {
	d.private_key.init(slv_private_key, priv_key_prefix)
}

func (d *Decrypter) decrypt(ciphered_data []byte, ephemeral_pub_key [32]byte) (decrypted_data []byte, err error) {
	decrypted_data, success := box.Open(nil, ciphered_data, &d.nonce, &ephemeral_pub_key, &d.private_key.key)
	if !success {
		return nil, error_decryption
	}
	return decrypted_data, nil
}

func (d *Decrypter) DecryptData(ciphertext string) (decrypted_data []byte, err error) {
	payload := strings.Replace(ciphertext, secret_prefix, "", 1)
	ciphered_data_split := strings.Split(payload, "_")
	ephemeral_pub_key := [32]byte(base58.Decode(ciphered_data_split[1]))
	ciphered_data_bytes := base58.Decode(ciphered_data_split[2])
	return d.decrypt(ciphered_data_bytes, ephemeral_pub_key)
}
