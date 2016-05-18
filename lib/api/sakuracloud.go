package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	sakuraCloudAPIRoot                = "https://secure.sakura.ad.jp/cloud/zone"
	sakuraCloudAPIRootSuffix          = "api/cloud/1.1"
	sakuraCloudPublicImageSearchWords = "Ubuntu%20Server%2016.04%20LTS%2064bit"
)

// Client type of sakuracloud api client config values
type Client struct {
	AccessToken       string
	AccessTokenSecret string
	Region            string
}

// NewClient create new sakuracloud api client
func NewClient(token, tokenSecret, region string) *Client {
	return &Client{AccessToken: token, AccessTokenSecret: tokenSecret, Region: region}
}

func (c *Client) getEndpoint() string {
	return fmt.Sprintf("%s/%s/%s", sakuraCloudAPIRoot, c.Region, sakuraCloudAPIRootSuffix)
}

func (c *Client) isOkStatus(code int) bool {
	codes := map[int]bool{
		200: true,
		201: true,
		202: true,
		204: true,
		305: false,
		400: false,
		401: false,
		403: false,
		404: false,
		405: false,
		406: false,
		408: false,
		409: false,
		411: false,
		413: false,
		415: false,
		500: false,
		503: false,
	}
	return codes[code]
}

func (c *Client) newRequest(method, uri string, body interface{}) ([]byte, error) {
	var (
		client = &http.Client{}
		url    = fmt.Sprintf("%s/%s", c.getEndpoint(), uri)
		err    error
		req    *http.Request
	)

	if body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, bytes.NewBuffer(bodyJSON))

	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Error with request: %v - %q", url, err)
	}

	req.SetBasicAuth(c.AccessToken, c.AccessTokenSecret)
	req.Method = method

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if !c.isOkStatus(resp.StatusCode) {
		return nil, fmt.Errorf("Error in response: %s", data)
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}
