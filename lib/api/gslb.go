package api

import (
	"encoding/json"
	"fmt"
	//	"strings"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
)

type searchGslbResponse struct {
	Total                  int                            `json:",omitempty"`
	From                   int                            `json:",omitempty"`
	Count                  int                            `json:",omitempty"`
	CommonServiceGslbItems []sakura.CommonServiceGslbItem `json:"CommonServiceItems,omitempty"`
}

type gslbRequest struct {
	CommonServiceGslbItem *sakura.CommonServiceGslbItem `json:"CommonServiceItem,omitempty"`
	From                  int                           `json:",omitempty"`
	Count                 int                           `json:",omitempty"`
	Sort                  []string                      `json:",omitempty"`
	Filter                map[string]interface{}        `json:",omitempty"`
	Exclude               []string                      `json:",omitempty"`
	Include               []string                      `json:",omitempty"`
}

type gslbResponse struct {
	*sakura.ResultFlagValue
	*sakura.CommonServiceGslbItem `json:"CommonServiceItem,omitempty"`
}

// SetupGslbRecord create or update Gslb
func (c *Client) SetupGslbRecord(gslbName string, ip string) ([]string, error) {

	gslbItem, err := c.getGslbCommonServiceItem(gslbName)

	if err != nil {
		return nil, err
	}
	gslbItem.Settings.GSLB.AddServer(ip)
	res, err := c.updateGslbServers(gslbItem)
	if err != nil {
		return nil, err
	}

	if gslbItem.ID == "" {
		return []string{res.Status.FQDN}, nil
	}
	return nil, nil

}

// DeleteGslbServer delete gslb server
func (c *Client) DeleteGslbServer(gslbName string, ip string) error {
	gslbItem, err := c.getGslbCommonServiceItem(gslbName)
	if err != nil {
		return err
	}
	gslbItem.Settings.GSLB.DeleteServer(ip)

	if gslbItem.HasGslbServer() {
		_, err = c.updateGslbServers(gslbItem)
		if err != nil {
			return err
		}

	} else {
		err = c.deleteCommonServiceGslbItem(gslbItem)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *Client) getGslbCommonServiceItem(gslbName string) (*sakura.CommonServiceGslbItem, error) {

	var (
		method = "GET"
		uri    = "commonserviceitem"
		body   = sakura.Request{
			Filter: map[string]interface{}{
				"Provider.Class": "gslb",
				"Name":           gslbName,
			},
		}
	)

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	uri = fmt.Sprintf("%s?%s", uri, bodyJSON)
	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	var gslb searchGslbResponse
	var gslbItem *sakura.CommonServiceGslbItem
	if err := json.Unmarshal(data, &gslb); err != nil {
		return nil, err
	}

	if gslb.Count > 0 {
		gslbItem = &gslb.CommonServiceGslbItems[0]
	} else {
		gslbItem = sakura.CreateNewGslbCommonServiceItem(gslbName)
	}

	return gslbItem, nil
}

func (c *Client) updateGslbServers(gslbItem *sakura.CommonServiceGslbItem) (*sakura.CommonServiceGslbItem, error) {

	var (
		method string
		uri    string
	)
	if gslbItem.ID == "" {
		method = "POST"
		uri = "/commonserviceitem"

	} else {
		method = "PUT"
		uri = fmt.Sprintf("/commonserviceitem/%s", gslbItem.ID)
	}
	n := gslbRequest{
		CommonServiceGslbItem: gslbItem,
	}

	data, err := c.newRequest(method, uri, n)
	if err != nil {
		return nil, err
	}
	var res gslbResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.CommonServiceGslbItem, nil
}

func (c *Client) deleteCommonServiceGslbItem(item *sakura.CommonServiceGslbItem) error {
	var (
		method string
		uri    string
	)
	method = "DELETE"
	uri = fmt.Sprintf("/commonserviceitem/%s", item.ID)

	_, err := c.newRequest(method, uri, item)
	if err != nil {
		return err
	}

	return nil

}
