package option

import (
	"encoding/json"
	"fmt"
)

// NewConfig return new FileOption
func NewConfig(version int, name string) *Config {
	return &Config{
		Version: version,
		Name:    name,
		Values:  map[string]string{},
	}
}

// Config type of option value
type Config struct {
	Version int
	Name    string
	Values  map[string]string
}

// Get return value
func (o *Config) Get(key string) (string, error) {
	value, ok := o.Values[key]
	if ok {
		return value, nil
	}

	return "", nil
}

// Set set option vlaue
func (o *Config) Set(key string, value string) error {
	o.Values[key] = value
	return nil
}

// Clear clear option key
func (o *Config) Clear(key string) error {
	_, ok := o.Values[key]
	if ok {
		delete(o.Values, key)
	}
	return nil
}

// List return all options
func (o *Config) List() (map[string]string, error) {
	return o.Values, nil
}

// MigrateOption migrate option JSON
func MigrateOption(o *Config, data []byte) (*Config, bool, error) {
	migrationPerformed := false
	if err := json.Unmarshal(data, o); err != nil {
		return nil, migrationPerformed, fmt.Errorf("Error unmarshalling sakuracloud options: %s", err)
	}
	migrationPerformed = true

	return o, migrationPerformed, nil
}
