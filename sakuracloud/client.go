package sakuracloud

import (
	"fmt"
	"github.com/sacloud/docker-machine-sakuracloud/version"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

// APIClient client for SakuraCloud API
type APIClient struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	Region            string // 後方互換
	client            *api.Client
	initialized       bool
}

// NewAPIClient returns new APIClient
func NewAPIClient(token string, secret string, zone string) *APIClient {
	return &APIClient{
		AccessToken:       token,
		AccessTokenSecret: secret,
		Zone:              zone,
		client:            api.NewClient(token, secret, zone),
		initialized:       false,
	}
}

// Init initialize APIClient
func (c *APIClient) Init() {
	if !c.initialized {
		if c.Zone == "" && c.Region != "" {
			c.Zone = c.Region
		}
		c.client = api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)
		c.client.UserAgent = fmt.Sprintf("docker-machine-sakuracloud/%s", version.Version)
		c.initialized = true
	}
}

// ValidateClientConfig validates client config
func (c *APIClient) ValidateClientConfig() error {
	if c.client.AccessToken == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token")
	}

	if c.client.AccessTokenSecret == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token-secret")
	}
	if c.client.Zone == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-zone")
	}
	return nil
}

// ServerBuilder returns new server builder
func (c *APIClient) ServerBuilder(osType, name, password string) builder.PublicArchiveUnixServerBuilder {
	return builder.ServerPublicArchiveUnix(c.client, ostype.StrToOSType(osType), name, password)
}

// IsValidPlan validates plan
func (c *APIClient) IsValidPlan(core int, memory int) (bool, error) {
	return c.client.Product.Server.IsValidPlan(core, memory, sacloud.PlanDefault)
}

// IsExistsPacketFilter returns true is PakcetFilter is exists
func (c *APIClient) IsExistsPacketFilter(id int64) (bool, error) {
	pf, err := c.client.PacketFilter.Read(id)
	if err != nil {
		if e, ok := err.(api.Error); ok && e.ResponseCode() == 404 {
			return false, nil
		}
		return false, err
	}
	return pf != nil, nil
}
