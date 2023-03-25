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
	pub_key_prefix       = "slv_pubkey_"
	priv_key_prefix      = "slv_privkey_"
	secret_prefix        = "slv_secret_"
	encrypted_key_prefix = "slv_enckey_"
)

func checksum(data []byte) string {
	sumBytes := sha1.Sum(data)
	sum := base58.Encode((sumBytes[:]))
	return sum[len(sum)-6:]
}

type key struct {
	id  string
	key [32]byte
}

func (key *key) init(slv_key_str, prefix string) {
	key_with_id_str := strings.Replace(slv_key_str, prefix, "", 1)
	key_with_id_split := strings.Split(key_with_id_str, "_")
	key.id = key_with_id_split[0]
	key.key = [32]byte(base58.Decode(key_with_id_split[1]))
}

func (key *key) to_string(prefix string) string {
	return fmt.Sprintf("%s%s_%s", prefix, key.id, base58.Encode(key.key[:]))
}

type KeyPair struct {
	pub  key
	priv key
}

func (kp *KeyPair) GetPublicKey() string {
	return kp.pub.to_string(pub_key_prefix)
}

func (kp *KeyPair) GetPrivateKey() string {
	return kp.priv.to_string(priv_key_prefix)
}

func (kp *KeyPair) generateLocally() (err error) {
	var pub_key, priv_key *[32]byte
	pub_key, priv_key, err = box.GenerateKey(rand.Reader)
	if err != nil {
		return
	}
	kp.pub.key = *pub_key
	kp.priv.key = *priv_key
	sum := checksum(kp.pub.key[:])
	kp.pub.id = sum
	kp.priv.id = sum
	return
}

func (kp *KeyPair) New() (err error) {
	err = kp.generateLocally()
	kp.priv.id = checksum(kp.pub.key[:])
	return
}
