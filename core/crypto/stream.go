package crypto

// func (publicKey *PublicKey) encryptStream(dst io.Writer, src io.Reader) error {
// 	header := append([]byte{byte(*publicKey.version), byte(*publicKey.keyType)}, publicKey.toBytes()...)
// 	if _, err := dst.Write(header); err != nil {
// 		return err
// 	}
// 	return publicKey.pubKey.EncryptStream(dst, src, true)
// }

// func (secretKey *SecretKey) decryptStream(dst io.Writer, src io.Reader) error {
// 	publicKey, err := secretKey.PublicKey()
// 	if err != nil {
// 		return err
// 	}
// 	if err != nil || !ciphered.IsEncryptedBy(publicKey) {
// 		return nil, errSecretKeyMismatch
// 	}
// 	data, err = secretKey.privKey.Decrypt(*ciphered.ciphertext)
// 	if err != nil {
// 		return nil, errDecryptionFailed
// 	}
// 	return
// }
