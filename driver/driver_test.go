package driver

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/machine/commands/mcndirs"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

const (
	testStoreDir          = ".store-test"
	machineTestName       = "test-host"
	machineTestCaCert     = "test-cert"
	machineTestPrivateKey = "test-key"
)

type DriverOptionsMock struct {
	Data map[string]interface{}
}

func (d DriverOptionsMock) String(key string) string {
	if value, ok := d.Data[key]; ok {
		return value.(string)
	}
	return ""
}

func (d DriverOptionsMock) StringSlice(key string) []string {
	if value, ok := d.Data[key]; ok {
		return value.([]string)
	}
	return []string{}
}

func (d DriverOptionsMock) Int(key string) int {
	if value, ok := d.Data[key]; ok {
		return value.(int)
	}
	return 0
}

func (d DriverOptionsMock) Bool(key string) bool {
	if value, ok := d.Data[key]; ok {
		return value.(bool)
	}
	return false
}

func cleanup() error {
	return os.RemoveAll(testStoreDir)
}

func getTestStorePath() (string, error) {
	tmpDir, err := ioutil.TempDir("", "machine-test-")
	if err != nil {
		return "", err
	}
	mcndirs.BaseDir = tmpDir
	return tmpDir, nil
}

func getDefaultTestDriverFlags() *DriverOptionsMock {
	return &DriverOptionsMock{
		Data: map[string]interface{}{
			"name":                            "test",
			"url":                             "unix:///var/run/docker.sock",
			"sakuracloud-access-token":        "accesstoken",
			"sakuracloud-access-token-secret": "secret",
			"sakuracloud-region":              "is1a",
		},
	}
}

func getTestDriver() (*Driver, error) {
	storePath, err := getTestStorePath()
	if err != nil {
		return nil, err
	}
	defer cleanup()

	d := NewDriver(machineTestName, storePath)
	d.SetConfigFromFlags(getDefaultTestDriverFlags())
	drv := d.(*Driver)
	return drv, nil
}

func TestSetConfigFromFlags(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{
			"sakuracloud-access-token":        "token",
			"sakuracloud-access-token-secret": "secret",
			"sakuracloud-region":              "region",
		},
		CreateFlags: driver.GetCreateFlags(),
	}

	driver.SetConfigFromFlags(checkFlags)

	//assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
}
