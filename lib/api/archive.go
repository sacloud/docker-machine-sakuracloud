package api

import (
	"encoding/json"
	"errors"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
)

// GetUbuntuArchiveID get ubuntu archive id
func (c *Client) GetUbuntuArchiveID() (string, error) {

	var (
		method = "GET"
		uri    = "archive"
		body   = sakura.Request{
			Filter: map[string]interface{}{
				"Name":  sakuraCloudPublicImageSearchWords,
				"Scope": "shared",
			},
			Include: []string{"ID", "Name"},
		}
	)

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	uri = fmt.Sprintf("%s?%s", uri, bodyJSON)
	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return "", err
	}
	var ubuntu sakura.SearchResponse
	if err := json.Unmarshal(data, &ubuntu); err != nil {
		return "", err
	}

	//すでに登録されている場合
	if ubuntu.Count > 0 {
		return ubuntu.Archives[0].ID, nil
	}

	return "", errors.New("Archive'Ubuntu Server 14.04 LTS 64bit' not found.")
}
