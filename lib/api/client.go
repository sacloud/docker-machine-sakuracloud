package api

import (
	"fmt"
	"github.com/sacloud/libsacloud/api"
	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
	"github.com/yamamoto-febc/docker-machine-sakuracloud/version"
)

type APIClient struct {
	AccessToken       string
	AccessTokenSecret string
	Region            string
	client            *api.Client
	initialized       bool
}

func NewAPIClient(token string, secret string, zone string) *APIClient {
	return &APIClient{
		AccessToken:       token,
		AccessTokenSecret: secret,
		Region:            zone,
		client:            api.NewClient(token, secret, zone),
		initialized:       true,
	}
}

func (c *APIClient) Init() {
	if !c.initialized {
		c.client = api.NewClient(c.AccessToken, c.AccessTokenSecret, c.Region)
		c.client.UserAgent = fmt.Sprintf("docker-machine-sakuracloud/%s", version.Version)
	}
}

func (c *APIClient) ValidateClientConfig() error {
	if c.client.AccessToken == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token")
	}

	if c.client.AccessTokenSecret == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token-secret")
	}
	return nil
}

func (c *APIClient) IsValidPlan(core int, memory int) (bool, error) {
	return c.client.Product.Server.IsValidPlan(core, memory)
}

func (c *APIClient) GetUbuntuArchiveID() (string, error) {
	res, err := c.client.Archive.FindByOSType(ostype.Ubuntu)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", res.ID), nil

}

func (c *APIClient) NewServer() *sacloud.Server {
	return c.client.Server.New()
}

func (c *APIClient) NewDisk() *sacloud.Disk {
	return c.client.Disk.New()
}
