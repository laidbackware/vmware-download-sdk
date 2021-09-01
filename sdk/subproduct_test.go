package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSubProductsSlice(t *testing.T) {
	var subProducts []map[string]string
	subProducts, err = basicClient.GetSubProductsSlice("vmware_vsphere")
	assert.Nil(t, err)
	assert.Greater(t, len(subProducts), 16, "Expected response to contain at least 16 items")
}

func TestGetSubProductsMap(t *testing.T) {
	var subProducts map[string]map[string]SubProductDetails
	subProducts, err = basicClient.GetSubProductsMap("vmware_vsphere")
	assert.Nil(t, err)
	assert.Contains(t, subProducts, "VMTOOLS")
}

func TestGetSubProductsMapInvalidSlug(t *testing.T) {
	var subProductMap map[string]map[string]SubProductDetails
	subProductMap, err = basicClient.GetSubProductsMap("vsphere")
	assert.NotNil(t, err)
	assert.Empty(t, subProductMap, "Expected response to be empty")
}

func TestGetSubProductsDetails(t *testing.T) {
	var subProductDetails SubProductDetails
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "VMTOOLS", "6_7")
	assert.Nil(t, err)
	assert.NotEmpty(t, subProductDetails.NiceName, "Expected response to not be empty")
}

func TestGetSubProductsDetailsInvalidSubProduct(t *testing.T) {
	var subProductDetails SubProductDetails
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "TOOLS", "6_7")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, subProductDetails.NiceName, "Expected response to be empty")
}

func TestGetSubProductsDetailsInvalidMajorVersion(t *testing.T) {
	var subProductDetails SubProductDetails
	subProductDetails, err = basicClient.GetSubProductDetails("vmware_vsphere", "VMTOOLS", "5_5")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProductMajorVersion)
	assert.Empty(t, subProductDetails.NiceName, "Expected response to be empty")
}
