package api

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var client *APIClient
var setupFuncs []func() = []func(){}
var tearDownFuncs []func() = []func(){}

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
		region = "is1b"
	}
	client = NewAPIClient(accessToken, accessTokenSecret, region)

	for _, f := range setupFuncs {
		f()
	}

	ret := m.Run()

	for _, f := range tearDownFuncs {
		f()
	}

	os.Exit(ret)
}

func TestGetUbuntuArchiveID(t *testing.T) {
	id, err := client.GetUbuntuArchiveID()
	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	t.Logf("ubuntu archive ID : %s", id)
}

