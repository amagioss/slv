package cipher

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/nacl/box"
)

const (
	pub_key_prefix  = "slv_pubkey_"
	priv_key_prefix = "slv_privkey_"
	secret_prefix   = "slv_secret_"
)

type pub struct {
	key [32]byte
}

func checksum(data []byte) string {
	sumBytes := sha1.Sum(data)
	sum := base58.Encode((sumBytes[:]))
	return sum[len(sum)-6:]
}

func (pub *pub) Init(slv_public_key string) {
	keyStr := strings.Replace(slv_public_key, pub_key_prefix, "", 1)
	pub.key = [32]byte(base58.Decode(keyStr))
}

type priv struct {
	key [32]byte
	id  string
}

func (priv *priv) Init(slv_private_key string) {
	keyWithIdStr := strings.Replace(slv_private_key, priv_key_prefix, "", 1)
	keyWithIdSlice := strings.Split(keyWithIdStr, "_")
	priv.id = keyWithIdSlice[0]
	priv.key = [32]byte(base58.Decode(keyWithIdSlice[1]))
}

type KeyPair struct {
	pub  pub
	priv priv
}

func (kp *KeyPair) GetPublicKey() string {
	return fmt.Sprintf("%s%s", pub_key_prefix, base58.Encode(kp.pub.key[:]))
}

func (kp *KeyPair) GetPrivateKey() string {
	return fmt.Sprintf("%s%s_%s", priv_key_prefix, kp.priv.id, base58.Encode(kp.priv.key[:]))
}

func (kp *KeyPair) generateLocally() (err error) {
	var pub_key, priv_key *[32]byte
	pub_key, priv_key, err = box.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	kp.pub.key = *pub_key
	kp.priv.key = *priv_key
	kp.priv.id = checksum(kp.pub.key[:])
	return
}

func (kp *KeyPair) New() (err error) {
	err = kp.generateLocally()
	kp.priv.id = checksum(kp.pub.key[:])
	return
}

type Encrypter struct {
	public_key        pub
	ephemeral_keypair KeyPair
	shared_key        [32]byte
	nonce             [24]byte
}

func (e *Encrypter) New(slv_public_key string) (err error) {
	e.public_key.Init(slv_public_key)
	err = e.ephemeral_keypair.generateLocally()
	if err != nil {
		return err
	}
	box.Precompute(&e.shared_key, &e.ephemeral_keypair.priv.key, &e.public_key.key)
	return
}

func (e *Encrypter) encrypt(data []byte) (ciphertext []byte) {
	return box.SealAfterPrecomputation(nil, data, &e.nonce, &e.shared_key)
}

func (e *Encrypter) EncryptData(data []byte) (ciphertext string) {
	encrypted := e.encrypt(data)
	hash := checksum(data)
	return fmt.Sprintf("%s%s_%s_%s", secret_prefix, hash, base58.Encode(e.ephemeral_keypair.pub.key[:]), base58.Encode(encrypted))
}
