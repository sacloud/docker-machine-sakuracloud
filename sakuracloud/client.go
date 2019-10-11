package sakuracloud

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/sacloud/docker-machine-sakuracloud/version"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"github.com/sacloud/libsacloud/v2/utils/server"
)

// APIClient client for SakuraCloud API
type APIClient struct {
	AccessToken       string
	AccessTokenSecret string
	Zone              string
	Region            string // 後方互換

	/*
		Note: 以下エクスポートしていない項目は適切にconfig storeからのUnmarshalに対応すること
			  -> 例: APIClient.Init()
	*/

	caller   sacloud.APICaller
	initOnce sync.Once
}

// NewAPIClient returns new APIClient
func NewAPIClient(token string, secret string, zone string) *APIClient {
	caller := initCaller(token, secret)

	return &APIClient{
		AccessToken:       token,
		AccessTokenSecret: secret,
		Zone:              zone,
		caller:            caller,
	}
}

func initCaller(token, secret string) sacloud.APICaller {
	return &sacloud.Client{
		AccessToken:            token,
		AccessTokenSecret:      secret,
		UserAgent:              fmt.Sprintf("docker-machine-sakuracloud/%s", version.Version),
		AcceptLanguage:         sacloud.APIDefaultAcceptLanguage,
		DefaultTimeoutDuration: sacloud.APIDefaultTimeoutDuration,
		RetryMax:               sacloud.APIDefaultRetryMax,
		RetryInterval:          sacloud.APIDefaultRetryInterval,
		HTTPClient:             http.DefaultClient,
	}
}

// ServerBuilderClient returns a api client for ServerBuilder
func (c *APIClient) ServerBuilderClient() *server.BuildersAPIClient {
	return server.NewBuildersAPIClient(c.caller)
}

// Init initialize APIClient
func (c *APIClient) Init() {
	c.initOnce.Do(func() {
		if c.Zone == "" && c.Region != "" {
			c.Zone = c.Region
		}
	})
	if c.caller == nil {
		c.caller = initCaller(c.AccessToken, c.AccessTokenSecret)
	}
}

// ValidateClientConfig validates client config
func (c *APIClient) ValidateClientConfig() error {
	c.Init()

	if c.AccessToken == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token")
	}

	if c.AccessTokenSecret == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-access-token-secret")
	}
	if c.Zone == "" {
		return fmt.Errorf("Missing required setting - --sakuracloud-zone")
	}
	return nil
}

// IsValidPlan validates plan
func (c *APIClient) IsValidPlan(core int, memory int) (bool, error) {
	plan, err := server.FindPlan(context.Background(), sacloud.NewServerPlanOp(c.caller), c.Zone, &server.FindPlanRequest{
		CPU:      core,
		MemoryGB: memory,
	})
	if err != nil {
		return false, err
	}
	exists := plan != nil
	return exists, nil
}

// IsExistsPacketFilter returns true is PakcetFilter is exists
func (c *APIClient) IsExistsPacketFilter(id types.ID) (bool, error) {
	pf, err := sacloud.NewPacketFilterOp(c.caller).Read(context.Background(), c.Zone, id)
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return pf != nil, nil
}
