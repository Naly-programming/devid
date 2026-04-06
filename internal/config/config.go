package config

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// ErrNoIdentity is returned when identity.toml does not exist.
var ErrNoIdentity = errors.New("identity.toml not found - run devid init")

// homeDirFunc is overridable for testing.
var homeDirFunc = os.UserHomeDir

// SetHomeDir overrides the home directory resolution. Used in tests.
func SetHomeDir(fn func() (string, error)) {
	homeDirFunc = fn
}

// DevidDir returns the path to ~/.devid/, creating it if needed.
func DevidDir() (string, error) {
	home, err := homeDirFunc()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".devid")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// QueueDir returns the path to ~/.devid/queue/, creating it if needed.
func QueueDir() (string, error) {
	devid, err := DevidDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(devid, "queue")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// IdentityPath returns the path to ~/.devid/identity.toml.
func IdentityPath() (string, error) {
	devid, err := DevidDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(devid, "identity.toml"), nil
}

// Exists returns true if identity.toml exists on disk.
func Exists() bool {
	p, err := IdentityPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}

// Load reads identity.toml from disk and decodes it.
func Load() (*Identity, error) {
	p, err := IdentityPath()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return nil, ErrNoIdentity
	}
	var id Identity
	if _, err := toml.DecodeFile(p, &id); err != nil {
		return nil, err
	}
	return &id, nil
}

// Save writes the identity to ~/.devid/identity.toml atomically.
func Save(id *Identity) error {
	id.Meta.UpdatedAt = time.Now().UTC()

	p, err := IdentityPath()
	if err != nil {
		return err
	}

	tmp := p + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	if err := toml.NewEncoder(f).Encode(id); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}

	return os.Rename(tmp, p)
}
