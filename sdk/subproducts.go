package sdk

import (
	"errors"
	// "fmt"
	// "reflect"
	"regexp"
	"sort"
	"strings"
)

type SubProduct struct {
	ProductName string
	ProductCode string
	DlgListByVersion map[string]DlgList
}


type SubProductSliceElement struct {
	Name string
	Description string
}

var ErrorInvalidSubProduct = errors.New("subproduct: invalid subproduct requested")
var ErrorInvalidSubProductMajorVersion = errors.New("subproduct: invalid major version requested")

func (c *Client) GetSubProductsMap(slug string) (data map[string]SubProduct, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	subProductMap := make(map[string]SubProduct)

	majorVersions, _ := c.GetMajorVersionsSlice(slug)
	if err != nil {
		return
	}

	// Iterate major product versions and extract all unique products
	// All version information is stripped
	for _, majorVersion := range majorVersions {
		dlgEditionsList, _ := c.GetDlgEditionsList(slug, majorVersion)
		for _, dlgEdition := range dlgEditionsList {
			for _, dlgList := range dlgEdition.DlgList {
				// Remove versions from the productCode and productName to allow to be generic
				re := regexp.MustCompile(`[0-9]+.*`)
				productCode := re.ReplaceAllString(dlgList.Code, "")
				productCode = strings.ToLower(strings.TrimSuffix(strings.TrimSuffix(productCode, "_"), "-"))
				productName := re.ReplaceAllString(dlgList.Name, "")
				productName = strings.TrimSpace(productName)

				// Initalize the struct for a product code for the first time
				if _, ok := subProductMap[productCode]; !ok {
					subProductMap[productCode] = SubProduct{
						ProductName: productName,
						ProductCode: productCode,
						DlgListByVersion: make(map[string]DlgList),
					}
				}

				subProductMap[productCode].DlgListByVersion[majorVersion] = dlgList
			}
		}
	}

	data = subProductMap
	return
}

func (c *Client) GetSubProductsSlice(slug string) (data []SubProduct, err error) {
	subProductMap, _ := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	// Sort keys to output sorted slice
	keys := make([]string, len(subProductMap))
	i := 0
	for key := range subProductMap {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
		
	// Append to array using sorted keys to fetch from map
	for _, key := range keys {
		data = append(data, subProductMap[key])
	}

	return
}

func (c *Client) ValidateSubProduct(slug, subProduct string) (err error) {
	subProductMap, _ := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	if _, ok := subProductMap[subProduct]; !ok {
		err = ErrorInvalidSubProduct
	}

	return
}

func (c *Client) GetSubProductDetails(slug, subProduct, majorVersion string) (data DlgList, err error) {
	subProducts, err := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	if subProduct, ok := subProducts[subProduct]; ok {
		if dlgList, ok := subProduct.DlgListByVersion[majorVersion]; ok {
			data = dlgList
		} else {
			err = ErrorInvalidSubProductMajorVersion
		}

	} else {
		err = ErrorInvalidSubProduct
	}

	return
}
