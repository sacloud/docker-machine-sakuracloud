package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestServerFuncs(t *testing.T) {

	newItem := client.client.Server.New()
	newItem.Name = "test_docker_machine_sakuracloud"
	newItem.Description = "before"
	newItem.SetServerPlanByID("1001")   // 1core 1GBメモリ
	newItem.AddPublicNWConnectedParam() //公開セグメントに接続

	server, err := client.Create(newItem, "")
	assert.NotEmpty(t, server)
	assert.NoError(t, err)

	//err = client.ConnectPacketFilterToSharedNIC(server, "112800627722")
	//assert.NoError(t, err)

	ip, err := client.GetIP(server.ID, false)
	assert.NotEmpty(t, ip)
	assert.NoError(t, err)

	err = client.Delete(server.ID, []string{})
	assert.NoError(t, err)

}
