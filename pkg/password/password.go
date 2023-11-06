// The password package provides utility functionality for comparing and hashing passwords.
// Hashing is done based on the [Argon2 hashing algorithm](https://en.wikipedia.org/wiki/Argon2).
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Error definitions
var (
	ErrIncorrectHash    = errors.New("incorrectly formatted hash provided")
	ErrIncorrectSaltLen = errors.New("salt length should be a positive integer")
)

// Default config values
const (
	DefaultTime       uint32 = 1
	DefaultMemory     uint32 = 64 * 1024
	DefaultThreads    uint8  = 4
	DefaultKeyLength  uint32 = 4
	DefaultSaltLength int    = 16

	AmountOfHashParts int = 6
)

// Config password manager configuration
type Config struct {
	Time    uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen int
}

// IManager password manager interface
type IManager interface {
	GenerateHash(password string) (string, error)
	ComparePassword(password, hash string) (bool, error)
}

// NewDefaultManager creates a new default manager
func NewDefaultManager() IManager {
	return NewManager(Config{
		Time:    DefaultTime,
		Memory:  DefaultMemory,
		Threads: DefaultThreads,
		KeyLen:  DefaultKeyLength,
		SaltLen: DefaultSaltLength,
	})
}

// NewManager creates a new password manager for given config
func NewManager(conf Config) IManager {
	return &manager{
		conf: conf,
	}
}

// manager implementation of IManager
type manager struct {
	conf Config
}

// GenerateHash generates a new hash for given password with config of manager
func (m *manager) GenerateHash(password string) (string, error) {
	if m.conf.SaltLen <= 0 {
		return "", ErrIncorrectSaltLen
	}
	// Generate a Salt
	salt := make([]byte, m.conf.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, m.conf.Time, m.conf.Memory, m.conf.Threads, m.conf.KeyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, m.conf.Memory, m.conf.Time, m.conf.Threads, b64Salt, b64Hash)

	return full, nil
}

// ComparePassword checks if password equals given hash
func (m *manager) ComparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != AmountOfHashParts {
		return false, ErrIncorrectHash
	}

	c := &Config{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &c.Memory, &c.Time, &c.Threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	c.KeyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, c.Time, c.Memory, c.Threads, c.KeyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}
