package hash

import (
	"crypto/sha1"
	"fmt"
)

// SHA1Hasher uses SHA1 to hash passwords with passed salt
type SHA1Hasher struct {
	salt string
}

// NewSHA1Hasher is hasher constructor
func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

// Hash return SHA1 password hash
func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
