package api

// SetupDnsRecord get dns zone commonserviceitem id
func (c *APIClient) SetupDnsRecord(zoneName string, hostName string, ip string) ([]string, error) {

	return c.client.DNS.SetupDNSRecord(zoneName, hostName, ip)
}

// DeleteDnsRecord delete dns record
func (c *APIClient) DeleteDnsRecord(zoneName string, hostName string, ip string) error {
	return c.client.DNS.DeleteDNSRecord(zoneName, hostName, ip)
}
