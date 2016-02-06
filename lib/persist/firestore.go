package persist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/option"
)

const (
	// OptionFileVersion version
	OptionFileVersion = 1
)

// Filestore Store of settings
type Filestore struct {
	Path string
}

// NewFilestore initialize Filestore
func NewFilestore(path string) *Filestore {
	return &Filestore{
		Path: path,
	}
}

// GetDriversDir return setting file dir
func (s Filestore) GetDriversDir() string {
	return filepath.Join(s.Path, "drivers", "sakuracloud")
}

func (s Filestore) saveToFile(data []byte, file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ioutil.WriteFile(file, data, 0600)
	}

	tmpfi, err := ioutil.TempFile(filepath.Dir(file), "options.json.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfi.Name())

	if err = ioutil.WriteFile(tmpfi.Name(), data, 0600); err != nil {
		return err
	}

	if err = tmpfi.Close(); err != nil {
		return err
	}

	if err = os.Remove(file); err != nil {
		return err
	}

	if err = os.Rename(tmpfi.Name(), file); err != nil {
		return err
	}
	return nil
}

// Save persists a setting in the store
func (s Filestore) Save(opt *option.Config) error {
	data, err := json.MarshalIndent(opt, "", "    ")
	if err != nil {
		return err
	}

	optPath := filepath.Join(s.GetDriversDir(), opt.Name)

	// Ensure that the directory we want to save to exists.
	if err := os.MkdirAll(optPath, 0700); err != nil {
		return err
	}

	return s.saveToFile(data, filepath.Join(optPath, "options.json"))
}

// Remove removes a setting from the store
func (s Filestore) Remove(name string) error {
	optPath := filepath.Join(s.GetDriversDir(), name)
	return os.RemoveAll(optPath)
}

// List returns a list of all setting in the store
func (s Filestore) List() ([]string, error) {
	dir, err := ioutil.ReadDir(s.GetDriversDir())
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	optNames := []string{}

	for _, file := range dir {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") {
			optNames = append(optNames, file.Name())
		}
	}

	return optNames, nil
}

// Exists returns whether a setting exists or not
func (s Filestore) Exists(name string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.GetDriversDir(), name, "options.json"))

	if os.IsNotExist(err) {
		return false, nil
	} else if err == nil {
		return true, nil
	}

	return false, err
}

func (s Filestore) loadConfig(o *option.Config) error {
	data, err := ioutil.ReadFile(filepath.Join(s.GetDriversDir(), o.Name, "options.json"))
	if err != nil {
		return err
	}

	// Remember the machine name so we don't have to pass it through each
	// struct in the migration.
	name := o.Name

	migratedOption, migrationPerformed, err := option.MigrateOption(o, data)
	if err != nil {
		return fmt.Errorf("Error getting migrated option: %s", err)
	}

	*o = *migratedOption

	o.Name = name

	// If we end up performing a migration, we should save afterwards so we don't have to do it again on subsequent invocations.
	if migrationPerformed {
		if err := s.saveToFile(data, filepath.Join(s.GetDriversDir(), o.Name, "options.json.bak")); err != nil {
			return fmt.Errorf("Error attempting to save backup after migration: %s", err)
		}

		if err := s.Save(o); err != nil {
			return fmt.Errorf("Error saving config after migration was performed: %s", err)
		}
	}

	return nil
}

// Load loads a setting by name
func (s Filestore) Load(name string) (*option.Config, error) {
	optPath := filepath.Join(s.GetDriversDir(), name)

	if _, err := os.Stat(optPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("setting does not exists : %s", name)
	}

	option := s.CreateNewConfig(name)

	if err := s.loadConfig(option); err != nil {
		return nil, err
	}

	return option, nil
}

// CreateNewConfig return create new option
func (s *Filestore) CreateNewConfig(name string) *option.Config {
	return option.NewConfig(OptionFileVersion, name)
}
