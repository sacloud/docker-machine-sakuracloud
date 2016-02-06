package persist

import (
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/option"
)

// Store type of SettingStore Object
type Store interface {
	// Exists returns whether a setting exists or not
	Exists(name string) (bool, error)

	// List returns a list of all setting in the store
	List() ([]string, error)

	// Load loads a setting by name
	Load(name string) (*option.Config, error)

	// Remove removes a setting from the store
	Remove(name string) error

	// Save persists a setting in the store
	Save(*option.Config) error
}
