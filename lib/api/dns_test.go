package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDnsDomain = "docker-machine-sakuracloud.com"

func TestUpdateDnsCommonServiceItem(t *testing.T) {
	item, err := client.getDnsCommonServiceItem(testDnsDomain) //存在しないため新たに作る
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testDnsDomain)

	//IPを追加して保存してみる
	item.Settings.DNS.AddDnsRecordSet("test1", "192.168.0.1")

	item, err = client.updateDnsRecord(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].Name, "test1")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].RData, "192.168.0.1")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].Type, "A")

	//IPを追加して保存してみる(２個目)
	item.Settings.DNS.AddDnsRecordSet("test2", "192.168.0.2")

	item, err = client.updateDnsRecord(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].Name, "test2")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].RData, "192.168.0.2")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].Type, "A")

}

func init() {
	setupFuncs = append(setupFuncs, cleanupDnsCommonServiceItem)
	tearDownFuncs = append(tearDownFuncs, cleanupDnsCommonServiceItem)
}

func cleanupDnsCommonServiceItem() {
	item, _ := client.getDnsCommonServiceItem(testDnsDomain)

	if item.ID != "" {
		client.deleteCommonServiceDnsItem(item)
	}
}
