package api

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

const testDomain = "docker-machine-sakuracloud.com"

var client *Client

func TestMain(m *testing.M) {
	//環境変数にトークン/シークレットがある場合のみテスト実施
	accessToken := os.Getenv("SAKURACLOUD_ACCESS_TOKEN")
	accessTokenSecret := os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET")

	if accessToken == "" || accessTokenSecret == "" {
		log.Fatal("Please Set ENV 'SAKURACLOUD_ACCESS_TOKEN' and 'SAKURACLOUD_ACCESS_TOKEN_SECRET'")
		os.Exit(0) // exit normal
	}
	region := os.Getenv("SAKURACLOUD_REGION")
	if region == "" {
		region = "is1a"
	}
	client = NewClient(accessToken, accessTokenSecret, region)

	cleanupCommonServiceItem(testDomain)

	ret := m.Run()

	cleanupCommonServiceItem(testDomain)

	os.Exit(ret)
}

func cleanupCommonServiceItem(domainName string) {
	item, _ := client.getDnsCommonServiceItem(testDomain)

	if item.ID != "" {
		client.deleteDnsServiceItem(item)
	}
}

func TestGetUbuntuArchiveID(t *testing.T) {
	id, err := client.GetUbuntuArchiveID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	t.Logf("ubuntu archive ID : %s", id)
}

func TestUpdateDnsCommonServiceItem(t *testing.T) {
	item, err := client.getDnsCommonServiceItem(testDomain) //存在しないため新たに作る
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testDomain)

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

func TestSetupDnsRecord(t *testing.T) {
	item, err := client.getDnsCommonServiceItem(testDomain) //存在しないため新たに作る
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testDomain)

	//IPを追加して保存してみる
	item.Settings.DNS.AddDnsRecordSet("test1", "192.168.0.1")

	item, err = client.updateDnsRecord(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testDomain)
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].Name, "test1")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].RData, "192.168.0.1")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[0].Type, "A")

	//IPを追加して保存してみる(２個目)
	item.Settings.DNS.AddDnsRecordSet("test2"+testDomain, "192.168.0.2")

	item, err = client.updateDnsRecord(item)
	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, item.Name, testDomain)
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].Name, "test2")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].RData, "192.168.0.2")
	assert.Equal(t, item.Settings.DNS.ResourceRecordSets[1].Type, "A")

}
