package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
	"strings"
)

type searchDnsResponse struct {
	Total                 int                           `json:",omitempty"`
	From                  int                           `json:",omitempty"`
	Count                 int                           `json:",omitempty"`
	CommonServiceDnsItems []sakura.CommonServiceDnsItem `json:"CommonServiceItems,omitempty"`
}

type dnsRequest struct {
	CommonServiceDnsItem *sakura.CommonServiceDnsItem `json:"CommonServiceItem,omitempty"`
	From                 int                          `json:",omitempty"`
	Count                int                          `json:",omitempty"`
	Sort                 []string                     `json:",omitempty"`
	Filter               map[string]interface{}       `json:",omitempty"`
	Exclude              []string                     `json:",omitempty"`
	Include              []string                     `json:",omitempty"`
}
type dnsResponse struct {
	*sakura.ResultFlagValue
	*sakura.CommonServiceDnsItem `json:"CommonServiceItem,omitempty"`
}

// SetupDnsRecord get dns zone commonserviceitem id
func (c *Client) SetupDnsRecord(zoneName string, hostName string, ip string) ([]string, error) {

	dnsItem, err := c.getDnsCommonServiceItem(zoneName)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(hostName, zoneName) {
		hostName = strings.Replace(hostName, zoneName, "", -1)
	}

	dnsItem.Settings.DNS.AddDnsRecordSet(hostName, ip)

	res, err := c.updateDnsRecord(dnsItem)
	if err != nil {
		return nil, err
	}

	if dnsItem.ID == "" {
		return res.Status.NS, nil
	}

	return nil, nil

}

// DeleteDnsRecord delete dns record
func (c *Client) DeleteDnsRecord(zoneName string, hostName string, ip string) error {
	dnsItem, err := c.getDnsCommonServiceItem(zoneName)
	if err != nil {
		return err
	}
	dnsItem.Settings.DNS.DeleteDnsRecordSet(hostName, ip)

	if dnsItem.HasDnsRecord() {
		_, err = c.updateDnsRecord(dnsItem)
		if err != nil {
			return err
		}

	} else {
		err = c.deleteCommonServiceDnsItem(dnsItem)
		if err != nil {
			return err
		}

	}
	return nil
}

func (c *Client) getDnsCommonServiceItem(zoneName string) (*sakura.CommonServiceDnsItem, error) {

	var (
		method = "GET"
		uri    = "commonserviceitem"
		body   = sakura.Request{
			Filter: map[string]interface{}{
				"Name":           zoneName,
				"Provider.Class": "dns",
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
	var dnsZone searchDnsResponse
	if err := json.Unmarshal(data, &dnsZone); err != nil {
		return nil, err
	}

	//すでに登録されている場合
	var dnsItem *sakura.CommonServiceDnsItem
	if dnsZone.Count > 0 {
		dnsItem = &dnsZone.CommonServiceDnsItems[0]
	} else {
		dnsItem = sakura.CreateNewDnsCommonServiceItem(zoneName)
	}

	return dnsItem, nil
}

func (c *Client) updateDnsRecord(dnsItem *sakura.CommonServiceDnsItem) (*sakura.CommonServiceDnsItem, error) {

	var (
		method string
		uri    string
	)
	if dnsItem.ID == "" {
		method = "POST"
		uri = "/commonserviceitem"

	} else {
		method = "PUT"
		uri = fmt.Sprintf("/commonserviceitem/%s", dnsItem.ID)
	}
	n := dnsRequest{
		CommonServiceDnsItem: dnsItem,
	}

	data, err := c.newRequest(method, uri, n)
	if err != nil {
		return nil, err
	}
	var res dnsResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.CommonServiceDnsItem, nil
}

func (c *Client) deleteCommonServiceDnsItem(item *sakura.CommonServiceDnsItem) error {
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
