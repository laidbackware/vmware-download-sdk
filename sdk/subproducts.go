package sdk

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

type SubProductDetails struct {
	NiceName string
	DlgList  DlgList
}

var ErrorInvalidSubProduct = errors.New("subproduct: invalid subproduct requested")
var ErrorInvalidSubProductMajorVersion = errors.New("subproduct: invalid major version requested")

// TODO may switch to slice, as having the value as a code may not be useful
func (c *Client) GetSubProductsMap(slug string) (data map[string]map[string]SubProductDetails, err error) {
	c.EnsureProductDetailMap()
	if err != nil {
		return
	}

	if _, ok := ProductDetailMap[slug]; !ok {
		err = ErrorInvalidSlug
		return
	}

	subProductMap := make(map[string]map[string]SubProductDetails)

	majorVersions, _ := c.GetMajorVersionsSlice(slug)
	if err != nil {
		return
	}

	// Iterate major product versions and extract all unique products
	// All version information is stripped
	for _, majorVersion := range majorVersions {
		dlgEditionsList, _ := c.GetDlgEditionsList(slug, majorVersion)
		for _, dlgEdition := range dlgEditionsList {
			for _, product := range dlgEdition.DlgList {
				// Remove versions from the code and name to allow to be generic
				re := regexp.MustCompile(`[0-9]+.*`)
				productCode := re.ReplaceAllString(product.Code, "")
				productCode = strings.TrimSuffix(strings.TrimSuffix(productCode, "_"), "-")
				productName := re.ReplaceAllString(product.Name, "")
				productName = strings.TrimSpace(productName)
				subProductDetails := SubProductDetails{
					NiceName: productName,
					DlgList:  product,
				}

				if len(subProductMap[productCode]) == 0 {
					majorMap := map[string]SubProductDetails{
						majorVersion: subProductDetails,
					}
					subProductMap[productCode] = majorMap
				} else {
					subProductMap[productCode][majorVersion] = subProductDetails
				}
			}
		}
	}

	data = subProductMap
	return
}

func (c *Client) GetSubProductsSlice(slug string) (data []map[string]string, err error) {
	subProductMap, _ := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	// Sort keys to output sorted slice
	keys := make([]string, len(subProductMap))
	i := 0
	for k := range subProductMap {
		keys[i] = k
		// keys = append(keys, k)
		i++
	}

	sort.Strings(keys)
	// TODO file to return nice name
	for k := range subProductMap {
		product := map[string]string{
			"code": k,
			// "name": v.niceName,
		}
		data = append(data, product)
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

// TODO need to return specific version for minor version!
func (c *Client) GetSubProductDetails(slug, subProduct, majorVersion string) (data SubProductDetails, err error) {
	subProducts, err := c.GetSubProductsMap(slug)
	if err != nil {
		return
	}

	if subProductVersions, ok := subProducts[subProduct]; ok {
		if subProductDetails, ok := subProductVersions[majorVersion]; ok {
			data = subProductDetails
		} else {
			err = ErrorInvalidSubProductMajorVersion
		}

	} else {
		err = ErrorInvalidSubProduct
	}

	return
}
