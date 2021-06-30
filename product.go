package myvmware

import (
	"encoding/json"
)

func (c *Client) Products() (data Products, err error) {
	res, err := c.HttpClient.Get(baseURL + "/channel/public/api/v1.0/products/getProductsAtoZ?isPrivate=true")
	if err != nil {
		return
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	return
}
