package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubProductsSlice(t *testing.T) {
	var subProducts []SubProduct
	subProducts, err = basicClient.GetSubProductsSlice("vmware_vsphere")
	assert.Nil(t, err)
	assert.Greater(t, len(subProducts), 16, "Expected response to contain at least 16 items")
}

func TestGetSubProductsMap(t *testing.T) {
	var subProducts map[string]SubProduct
	subProducts, err = basicClient.GetSubProductsMap("vmware_vsphere")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "vmtools")
}

func TestGetSubProductsMapInvalidSlug(t *testing.T) {
	var subProductMap map[string]SubProduct
	subProductMap, err = basicClient.GetSubProductsMap("vsphere")
	assert.NotNil(t, err)
	assert.Empty(t, subProductMap, "Expected response to be empty")
}

func TestGetSubProductsDetails(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "6_7")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProductDetails.Code, "Expected response to not be empty")
}

func TestGetSubProductsDetailsInvalidSubProduct(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "tools", "6_7")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}

func TestGetSubProductsDetailsInvalidMajorVersion(t *testing.T) {
	var subProductDetails DlgList
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "vmtools", "5_5")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProductMajorVersion)
	assert.Empty(t, subProductDetails.Code, "Expected response to be empty")
}
