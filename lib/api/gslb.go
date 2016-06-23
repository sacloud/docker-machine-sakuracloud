package api

// SetupGslbRecord create or update Gslb
func (c *APIClient) SetupGslbRecord(gslbName string, ip string) ([]string, error) {
	return c.client.GSLB.SetupGSLBRecord(gslbName, ip)
}

// DeleteGslbServer delete gslb server
func (c *APIClient) DeleteGslbServer(gslbName string, ip string) error {
	return c.client.GSLB.DeleteGSLBServer(gslbName, ip)
}
