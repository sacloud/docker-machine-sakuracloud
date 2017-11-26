package sakuracloud

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/builder"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/version"
)

type APIClient struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	client            *api.Client
	initialized       bool
}

func NewAPIClient(token string, secret string, zone string) *APIClient {
	return &APIClient{
		AccessToken:       token,
		AccessTokenSecret: secret,
		Zone:              zone,
		client:            api.NewClient(token, secret, zone),
		initialized:       false,
	}
}

func (c *APIClient) Init() {
	if !c.initialized {
		c.client = api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)
		c.client.UserAgent = fmt.Sprintf("docker-machine-sakuracloud/%s", version.Version)
		c.initialized = true
	}
}

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

func (c *APIClient) ServerBuilder(osType, name, password string) *builder.PublicArchiveUnixServerBuilder {
	return builder.ServerPublicArchiveUnix(c.client, ostype.StrToOSType(osType), name, password)
}

func (c *APIClient) IsValidPlan(core int, memory int) (bool, error) {
	return c.client.Product.Server.IsValidPlan(core, memory)
}

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
