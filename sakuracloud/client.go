package api

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/version"
)

type APIClient struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	Client            *api.Client
	initialized       bool
}

func NewAPIClient(token string, secret string, zone string) *APIClient {
	return &APIClient{
		AccessToken:       token,
		AccessTokenSecret: secret,
		Zone:              zone,
		Client:            api.NewClient(token, secret, zone),
		initialized:       true,
	}
}

func (c *APIClient) Init() {
	if !c.initialized {
		c.Client = api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Zone)
		c.Client.UserAgent = fmt.Sprintf("docker-machine-sakuracloud/%s", version.Version)
	}
}

func (c *APIClient) ValidateClientConfig() error {
	if c.Client.AccessToken == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token")
	}

	if c.Client.AccessTokenSecret == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token-secret")
	}
	return nil
}

func (c *APIClient) IsValidPlan(core int, memory int) (bool, error) {
	return c.Client.Product.Server.IsValidPlan(core, memory)
}
