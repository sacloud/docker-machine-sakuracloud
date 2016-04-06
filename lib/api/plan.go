package api

import (
	"encoding/json"
	"fmt"
	sakura "github.com/yamamoto-febc/docker-machine-sakuracloud/lib/cloud/resources"
	"strconv"
)

// IsValidPlan return valid plan
func (c *Client) IsValidPlan(core int, memGB int) (bool, error) {
	//assert args
	if core <= 0 {
		return false, fmt.Errorf("Invalid Parameter: CPU Core")
	}
	if memGB <= 0 {
		return false, fmt.Errorf("Invalid Parameter: Memory Size(GB)")
	}

	planID, _ := strconv.ParseInt(fmt.Sprintf("%d%03d", memGB, core), 10, 64)

	var (
		method = "GET"
		uri    = fmt.Sprintf("product/server/%d", planID)
	)

	data, err := c.newRequest(method, uri, nil)
	if err != nil {
		return false, err
	}
	var plan sakura.Response
	if err := json.Unmarshal(data, &plan); err != nil {
		return false, err
	}

	if plan.ServerPlan != nil {
		return true, nil
	}

	return false, fmt.Errorf("Server Plan[%d] Not Found", planID)

}
