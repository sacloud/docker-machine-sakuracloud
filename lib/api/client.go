package api

import (
	"fmt"
	"github.com/yamamoto-febc/libsacloud/api"
)

const (
	sakuraCloudPublicImageSearchWords = "Ubuntu Server 16.04 LTS 64bit"
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

	res, err := c.client.Archive.
		WithNameLike(sakuraCloudPublicImageSearchWords).
		WithSharedScope().
		Include("ID").
		Include("Name").
		Find()

	if err != nil {
		return "", err
	}

	//すでに登録されている場合
	if res.Count > 0 {
		return res.Archives[0].ID, nil
	}

	return "", fmt.Errorf("Archive'%s' not found.", sakuraCloudPublicImageSearchWords)
}
