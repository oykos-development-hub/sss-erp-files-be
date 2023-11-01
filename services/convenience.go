package services

import (
	"github.com/oykos-development-hub/celeritas"
)

type BaseServiceImpl struct {
	App *celeritas.Celeritas
}

func (s BaseServiceImpl) RandomString(n int) string {
	return s.App.RandomString(n)
}

func (s BaseServiceImpl) Encrypt(text string) (string, error) {
	enc := celeritas.Encryption{Key: []byte(s.App.EncryptionKey)}

	encrypted, err := enc.Encrypt(text)
	if err != nil {
		return "", err
	}
	return encrypted, nil
}

func (s BaseServiceImpl) Decrypt(crypto string) (string, error) {
	enc := celeritas.Encryption{Key: []byte(s.App.EncryptionKey)}

	decrypted, err := enc.Decrypt(crypto)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}
