package driver

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

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

	driver.SetConfigFromFlags(checkFlags) // nolint

	//assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
}
