package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type SettingsFile struct {
	fileName string
	crypt    bool
}

func NewSettingsFile(filename string, crypt bool) SettingsFile {
	return SettingsFile{
		fileName: filename,
		crypt:    crypt,
	}
}

func (s *SettingsFile) Load(rv interface{}) error {
	encBytes, err := ioutil.ReadFile(s.fileName)
	if err != nil {
		return fmt.Errorf("readfile: %w", err)
	}

	var bytes []byte
	if s.crypt {
		bytes, err = decrypt(encBytes)
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
	} else {
		bytes = encBytes
	}

	err = json.Unmarshal(bytes, &rv)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}

func (s *SettingsFile) Save(v interface{}) error {
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	var encBytes []byte
	if s.crypt {
		encBytes, err = encrypt(bytes)
		if err != nil {
			return fmt.Errorf("encrypt: %w", err)
		}
	} else {
		encBytes = bytes
	}

	err = ioutil.WriteFile(s.fileName, encBytes, 0644)
	if err != nil {
		return fmt.Errorf("writefile: %w", err)
	}

	return nil
}

var pph = createHash("djk3@1opsw^$xF")

func createHash(key string) []byte {
	hash := sha256.Sum256([]byte(key))
	return hash[:]
}

func encrypt(data []byte) ([]byte, error) {
	block, _ := aes.NewCipher(pph)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("ReadFull: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(pph)
	if err != nil {
		return nil, fmt.Errorf("NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("NewGCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	return plaintext, nil
}
