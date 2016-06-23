package cli

import (
	"github.com/docker/machine/commands/mcndirs"
	//"github.com/docker/machine/libmachine/mcnflag"
	"github.com/stretchr/testify/assert"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/lib/persist"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/spec"
	"os"
	"testing"
)

func cleanup(b API) {
	os.RemoveAll(b.GetDriversDir())
}

func TestGetConfigValue(t *testing.T) {
	api := NewClient()
	defer cleanup(api)

	// clear env var
	os.Unsetenv("SAKURACLOUD_REGION")

	config, err := api.GetConfigValue("region")
//	region := config.CurrentValue
	isDefault := config.IsDefault()

	assert.NoError(t, err)
//	assert.Equal(t, "is1a", region)
	assert.Equal(t, true, isDefault)

	saveConfig("region", "tk1a")

	config, err = api.GetConfigValue("region")
//	region = config.CurrentValue
	isDefault = config.IsDefault()

	assert.NoError(t, err)
//	assert.Equal(t, "tk1a", region)
//	assert.Equal(t, false, isDefault)

	// if setted env var , use it.
	os.Setenv("SAKURACLOUD_REGION", "is1b")

	config, err = api.GetConfigValue("region")
//	region = config.CurrentValue
	isDefault = config.IsDefault()

	assert.NoError(t, err)
//	assert.Equal(t, "is1b", region)
//	assert.Equal(t, false, isDefault)

}

func TestSetConfigValue(t *testing.T) {
	api := NewClient()
	defer cleanup(api)

	// clear env var
	os.Unsetenv("SAKURACLOUD_REGION")

	err := api.SetConfigValue("region", "tk1a")
	assert.NoError(t, err)

	config, err := api.Load(api.GetName())
	assert.NoError(t, err)

	value, err := config.Get("region")
	assert.NoError(t, err)

	assert.Equal(t, "tk1a", value)

}

func TestListConfigValue(t *testing.T) {
	api := NewClient()
	defer cleanup(api)

	// sakura.CliOptions
	configs, err := api.ListConfigValue()
	assert.NoError(t, err)

	assert.Equal(t, len(sakura.Options.CliOptions()), len(configs))
}

func saveConfig(key string, value string) {
	api := &client{
		targetSettingName: defaultConfigName,
		Filestore:         persist.NewFilestore(mcndirs.GetBaseDir()),
	}
	conf := api.CreateNewConfig(defaultConfigName)
	conf.Set(key, value)
	api.Filestore.Save(conf)

}
