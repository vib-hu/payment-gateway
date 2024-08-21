package encryption

type Encryption interface {
	Encrypt(plaintext string) (string, error)
	EncryptMultiple(plaintexts map[string]string) (map[string]string, error)
	Decrypt(ciphertext string) (string, error)
}
