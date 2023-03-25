package cipher

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/nacl/box"
)

var error_decryption = errors.New("Something went wrong while attempting to decrypt!")

func compress(data []byte) (compressed_data []byte, err error) {
	var compressed_data_buffer bytes.Buffer
	gzWriter := gzip.NewWriter(&compressed_data_buffer)
	_, err = gzWriter.Write(data)
	if err == nil {
		err = gzWriter.Close()
		if err == nil {
			return compressed_data_buffer.Bytes(), nil
		}
	}
	return
}

func compress_if_required(data []byte) (processed_data []byte, err error) {
	data_size := len(data)
	compressed_data, err := compress(data)
	if err == nil {
		compressed_size := len(compressed_data)
		if data_size <= compressed_size {
			processed_data = append([]byte{0}, data...)
		} else {
			processed_data = append([]byte{1}, compressed_data...)
		}
	}
	return processed_data, err
}

func decompress(compressed_data []byte) (data []byte, err error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(compressed_data))
	if err == nil {
		defer gzReader.Close()
		return ioutil.ReadAll(gzReader)
	}
	return
}

func decompress_if_required(possibly_compressed_data []byte) (data []byte, err error) {
	processed_data := possibly_compressed_data[1:]
	if possibly_compressed_data[0] == 1 {
		return decompress(processed_data)
	}
	return processed_data, nil
}

func newNonce() (nonce [24]byte, err error) {
	_, err = rand.Read(nonce[0:24])
	if err != nil {
		return
	}
	return
}

type Encrypter struct {
	public_key        [32]byte
	ephemeral_keypair KeyPair
	shared_key        [32]byte
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

func (e *Encrypter) encrypt(data []byte) (ciphered_data []byte, err error) {
	nonce, err := newNonce()
	if err == nil {
		processed_data, err := compress_if_required(data)
		if err == nil {
			ciphered_data = append(e.ephemeral_keypair.pub.key[:], nonce[:]...)
			ciphered_data = append(ciphered_data, box.SealAfterPrecomputation(nil, processed_data, &nonce, &e.shared_key)...)
		}
	}
	return
}

func (e *Encrypter) EncryptData(data []byte) (ciphertext string, err error) {
	encrypted, err := e.encrypt(data)
	if err == nil {
		hash := checksum(data)
		ciphertext = fmt.Sprintf("%s%s_%s", secret_prefix, hash, base58.Encode(encrypted))
	}
	return
}

type Decrypter struct {
	private_key key
}

func (d *Decrypter) New(slv_private_key string) {
	d.private_key.init(slv_private_key, priv_key_prefix)
}

func (d *Decrypter) decrypt(ciphered_data []byte) (decrypted_data []byte, err error) {
	ephemeral_pub_key := [32]byte(ciphered_data[0:32])
	nonce := [24]byte(ciphered_data[32:56])
	encrypted := ciphered_data[56:]
	decrypted_data, success := box.Open(nil, encrypted, &nonce, &ephemeral_pub_key, &d.private_key.key)
	if !success {
		return nil, error_decryption
	}
	return decompress_if_required(decrypted_data)
}

func (d *Decrypter) DecryptData(ciphertext string) (decrypted_data []byte, err error) {
	payload := strings.Replace(ciphertext, secret_prefix, "", 1)
	ciphered_data_split := strings.Split(payload, "_")
	ciphered_data_bytes := base58.Decode(ciphered_data_split[1])
	if err == nil {
		return d.decrypt(ciphered_data_bytes)
	}
	return
}
