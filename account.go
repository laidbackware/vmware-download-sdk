package myvmware

import (
	"encoding/json"
	"strings"
)

func (c *Client) CurrentUser() (data CurrentUser, err error) {
	res, err := c.HttpClient.Get(baseURL + "/vmwauth/loggedinuser")
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	return
}

func (c *Client) AccountInfo() (data AccountInfo, err error) {
	payload := `{"rowLimit": 10}`
	res, err := c.HttpClient.Post(baseURL+"/channel/api/v1.0/ems/accountinfo", "application/json", strings.NewReader(payload))
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	return
}
