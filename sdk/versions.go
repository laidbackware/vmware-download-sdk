package sdk

import (
	"errors"
	"sort"
	"strings"
)

type APIVersions struct {
	Code         string
	MajorVersion string
}

var ErrorInvalidVersionGlob = errors.New("versions: invalid glob. No versions found")

func (c *Client) GetVersionMap(slug, subProductName string) (data map[string]APIVersions, err error) {
	data = make(map[string]APIVersions)

	majorVersions, err := c.GetMajorVersionsSlice(slug)
	if err != nil {
		return
	}

	err = c.ValidateSubProduct(slug, subProductName)
	if err != nil {
		return
	}
	subProductName = strings.ToUpper(subProductName)

	// Iterate over all major versions of product to collect actual versions
	for _, majorVersion := range majorVersions {
		dlgEditionsList, _ := c.GetDlgEditionsList(slug, majorVersion)
		var foundProduct DlgList

		// Iterate each edition of each major version of product
		for _, edition := range dlgEditionsList {
			for _, product := range edition.DlgList {
				if strings.HasPrefix(product.Code, subProductName) {
					foundProduct = product
					break
				}
			}

			// When matching product is found pull nice version name and API code
			if foundProduct.Name != "" {
				var dlgHeader DlgHeader
				dlgHeader, err = c.GetDlgHeader(foundProduct.Code, foundProduct.ProductID)
				if err != nil {
					return
				}

				for _, version := range dlgHeader.Versions {
					aPIVersions := APIVersions{
						Code:         version.ID,
						MajorVersion: majorVersion,
					}
					data[version.Name] = aPIVersions
				}
			} else {
				err = ErrorInvalidSubProduct
			}
		}
	}

	return
}

func (c *Client) FindVersion(slug, subProduct, version string) (data APIVersions, err error) {
	versionMap, err := c.GetVersionMap(slug, subProduct)
	if err != nil {
		return
	}

	var searchVersion string
	if strings.Contains(version, "*") {
		searchVersion, err = c.FindVersionFromGlob(slug, subProduct, version, versionMap)
		if err != nil {
			return
		}
	} else {
		searchVersion = version
	}
	if _, ok := versionMap[searchVersion]; !ok {
		err = ErrorInvalidVersion
		return
	}

	data = versionMap[searchVersion]
	return
}

func (c *Client) FindVersionFromGlob(slug, subProduct, versionGlob string, versionMap map[string]APIVersions) (version string, err error) {
	// Ensure only one glob is defined
	globCount := strings.Count(versionGlob, "*")
	if globCount == 0 {
		err = ErrorNoGlob
		return
	} else if globCount > 1 {
		err = ErrorMultipleGlob
		return
	}

	// Extract prefix by removing *
	versionPrefix := strings.Split(versionGlob, "*")[0]

	sortedKeys := sortVersionMapKeys(versionMap)

	// Check if only * is provided as strings.Split returns empty if separator is found
	if versionPrefix == "" {
		// return the first entry, which is the highest number.
		version = sortedKeys[0]
		return
	} else {
		// return the first entry matching the prefix
		for _, key := range sortedKeys {
			if strings.HasPrefix(key, versionPrefix) {
				version = key
				return
			}
		}
	}

	err = ErrorInvalidVersionGlob
	return
}

func (c *Client) GetVersionSlice(slug, subProductName string) (data []string, err error) {
	versionMap, _ := c.GetVersionMap(slug, subProductName)
	if err != nil {return}

	data = sortVersionMapKeys(versionMap)

	return
}

func sortVersionMapKeys(versionMap map[string]APIVersions) (keys []string) {
	// Extract all keys which are the version strings and reverse sort them
	// This means the versions will go from high to low
	keys = make([]string, len(versionMap))
	i := 0
	for key := range versionMap {
		keys[i] = key
		i++
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return
}