package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersionsSuccess(t *testing.T) {
	var versions map[string]APIVersions
	versions, err := basicClient.GetVersionsMap("vmware_tools", "VMTOOLS")
	assert.Nil(t, err)
	assert.Greater(t, len(versions), 1, "Expected response to contain at least 1 item")
}

func TestGetVersionsInvalidSubProduct(t *testing.T) {
	var versions map[string]APIVersions
	versions, err := basicClient.GetVersionsMap("vmware_tools", "DUMMY")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, versions, "Expected response to be empty")
}

func TestGetVersionsInvalidSlug(t *testing.T) {
	var versions map[string]APIVersions
	versions, err := basicClient.GetVersionsMap("mware_tools", "VMTOOLS")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, versions, "Expected response to be empty")
}

func TestFindVersion(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "VMTOOLS", "11.1.1")
	assert.Nil(t, err)
	assert.NotEmpty(t, foundVersion.Code, "Expected response not to be empty")
}

func TestFindVersionInvalidSlug(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("mware_tools", "VMTOOLS", "11.1.1")
	assert.ErrorIs(t, err, ErrorInvalidSlug)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionInvalidVersion(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "VMTOOL", "11.1.1")
	assert.ErrorIs(t, err, ErrorInvalidSubProduct)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionInvalidSubProduct(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "VMTOOLS", "666")
	assert.ErrorIs(t, err, ErrorInvalidVersion)
	assert.Empty(t, foundVersion.Code, "Expected response to be empty")
}

func TestFindVersionGlob(t *testing.T) {
	var foundVersion APIVersions
	foundVersion, err = basicClient.FindVersion("vmware_tools", "VMTOOLS", "10.2.*")
	assert.Nil(t, err)
	assert.Equal(t, foundVersion.Code, "VMTOOLS1025")
}
