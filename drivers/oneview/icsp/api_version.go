package icsp

import (
	"encoding/json"

	"github.com/docker/machine/drivers/oneview/rest"
	"github.com/docker/machine/log"
)

// URLEndPoint .....
type URLEndPoint string

// URLEndPoint export this constant
const URLEndPoint = "/rest/version"

// APIVersion struct
type APIVersion struct {
	CurrentVersion int `json:"currentVersion,omitempty"`
	MinimumVersion int `json:"minimumVersion,omitempty"`
}

// GetAPIVersion - returns the api version for OneView server
// returns structure APIVersion
func (c *ICSPClient) GetAPIVersion() (APIVersion, error) {
	var (
		uri        = URLEndPoint
		apiversion APIVersion
	)

	//c.AuthHeaders := map[string]interface{}{"auth": []interface{}{auth}}
	c.SetAuthHeaderOptions(c.GetAuthHeaderMap())
	data, err := c.RestAPICall(rest.GET, uri, nil)
	if err != nil {
		return apiversion, err
	}

	log.Debugf("GetAPIVersion %s", data)
	if err := json.Unmarshal([]byte(data), &apiversion); err != nil {
		return apiversion, err
	}
	return apiversion, err
}
