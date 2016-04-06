package resources

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDnsRecordSets(t *testing.T) {
	records := DnsRecordSets{}
	assert.True(t, len(records.ResourceRecordSets) == 0)

	records.AddDnsRecordSet("test1", "192.168.0.1")
	assert.True(t, len(records.ResourceRecordSets) == 1)
	assert.Equal(t, records.ResourceRecordSets[0].RData, "192.168.0.1")
	//t.Logf("records:%#v", records)

	records.AddDnsRecordSet("test2", "192.168.0.2")
	assert.True(t, len(records.ResourceRecordSets) == 2)
	assert.Equal(t, records.ResourceRecordSets[1].RData, "192.168.0.2")
	//t.Logf("records:%#v", records)

	records.AddDnsRecordSet("test1", "192.168.0.3")
	assert.True(t, len(records.ResourceRecordSets) == 2)
	assert.Equal(t, records.ResourceRecordSets[0].RData, "192.168.0.3")
	//t.Logf("records:%#v", records)

	records.DeleteDnsRecordSet("test1", "192.168.0.1")
	assert.True(t, len(records.ResourceRecordSets) == 2)
	assert.Equal(t, records.ResourceRecordSets[0].RData, "192.168.0.3")

	records.DeleteDnsRecordSet("test1", "192.168.0.3")
	assert.True(t, len(records.ResourceRecordSets) == 1)
	assert.Equal(t, records.ResourceRecordSets[0].RData, "192.168.0.2")

}
