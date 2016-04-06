package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testGslbName = "test_docker_machine_sakuracloud_gslb"

func TestGslbGet(t *testing.T) {
	item, err := client.getGslbCommonServiceItem(testGslbName)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testGslbName)

	//IPを追加して保存してみる
	item.Settings.GSLB.AddServer("8.8.8.8")
	item, err = client.updateGslbServers(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Settings.GSLB.Servers[0].IPAddress, "8.8.8.8")
	assert.Equal(t, item.Settings.GSLB.Servers[0].Weight, "1")
	assert.Equal(t, item.Settings.GSLB.Servers[0].Enabled, "True")

	//IPを追加して保存してみる(2個目)
	item.Settings.GSLB.AddServer("8.8.4.4")
	item, err = client.updateGslbServers(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Settings.GSLB.Servers[1].IPAddress, "8.8.4.4")
	assert.Equal(t, item.Settings.GSLB.Servers[1].Weight, "1")
	assert.Equal(t, item.Settings.GSLB.Servers[1].Enabled, "True")

}

func init() {
	setupFuncs = append(setupFuncs, cleanupGslbCommonServiceItem)
	tearDownFuncs = append(tearDownFuncs, cleanupGslbCommonServiceItem)
}

func cleanupGslbCommonServiceItem() {
	item, _ := client.getGslbCommonServiceItem(testGslbName)

	if item.ID != "" {
		client.deleteCommonServiceGslbItem(item)
	}
}
