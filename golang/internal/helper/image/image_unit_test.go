//go:build unit
// +build unit

package image_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"

	imageHelper "github.com/dyrector-io/dyrectorio/golang/internal/helper/image"
	"github.com/dyrector-io/dyrectorio/golang/internal/pointer"
	"github.com/dyrector-io/dyrectorio/protobuf/go/agent"
)

type RegistryTestCase struct {
	Name        string
	Registry    *string
	RegistryUrl *string
	ExpectedUrl string
}

func TestRegistryWithTable(t *testing.T) {
	testCases := []RegistryTestCase{
		{
			Name:        "Test registry url",
			Registry:    pointer.NewPTR[string](""),
			RegistryUrl: pointer.NewPTR[string]("test"),
			ExpectedUrl: "test",
		},
		{
			Name:        "Test registry url priority",
			Registry:    pointer.NewPTR[string]("other"),
			RegistryUrl: pointer.NewPTR[string]("test"),
			ExpectedUrl: "test",
		},
		{
			Name:        "Test registry url empty",
			Registry:    nil,
			RegistryUrl: nil,
			ExpectedUrl: "",
		},
		{
			Name:        "Test registry url registry",
			Registry:    pointer.NewPTR[string]("other"),
			RegistryUrl: nil,
			ExpectedUrl: "other",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			if tC.RegistryUrl == nil {
				url := imageHelper.GetRegistryURL(tC.Registry, nil)
				assert.Equal(t, url, tC.ExpectedUrl)
			} else {
				auth := &imageHelper.RegistryAuth{URL: *tC.RegistryUrl}
				url := imageHelper.GetRegistryURL(tC.Registry, auth)
				assert.Equal(t, url, tC.ExpectedUrl)
			}
		})
	}
}

func TestProtoRegistryUrl(t *testing.T) {
	auth := &agent.RegistryAuth{
		Url: "test",
	}

	url := imageHelper.GetRegistryURLProto(nil, auth)
	assert.Equal(t, url, "test")
}

func TestProtoRegistryUrlPriority(t *testing.T) {
	registry := "other"
	auth := &agent.RegistryAuth{
		Url: "test",
	}

	url := imageHelper.GetRegistryURLProto(&registry, auth)
	assert.Equal(t, url, "test")
}

func TestProtoRegistryUrlRegistry(t *testing.T) {
	registry := "other"

	url := imageHelper.GetRegistryURLProto(&registry, nil)
	assert.Equal(t, url, "other")
}

func TestProtoRegistryUrlEmpty(t *testing.T) {
	url := imageHelper.GetRegistryURLProto(nil, nil)
	assert.Equal(t, url, "")
}

func TestExpandImageName(t *testing.T) {
	name, err := imageHelper.ExpandImageName("nginx")
	assert.NoError(t, err)
	assert.Equal(t, "docker.io/library/nginx:latest", name, "plain image is expanded to latest tag and it prefixing")

	name, err = imageHelper.ExpandImageName("nginx:tag")
	assert.NoError(t, err)
	assert.Equal(t, "docker.io/library/nginx:tag", name, "plain image name with tag keeps tag")

	name, err = imageHelper.ExpandImageName("my-reg.com/library/nginx")
	assert.NoError(t, err)
	assert.Equal(t, "my-reg.com/library/nginx:latest", name)

	name, err = imageHelper.ExpandImageName("my-reg.com/library/nginx:my-tag")
	assert.NoError(t, err)
	assert.Equal(t, "my-reg.com/library/nginx:my-tag", name)
}

func TestExpandImageNameWithTag(t *testing.T) {
	name, err := imageHelper.ExpandImageNameWithTag("nginx", "tag-1")
	assert.NoError(t, err)
	assert.Equal(t, "docker.io/library/nginx:tag-1", name)

	name, err = imageHelper.ExpandImageNameWithTag("nginx:tag", "tag-2")
	assert.NoError(t, err)
	assert.Equal(t, "docker.io/library/nginx:tag-2", name)

	name, err = imageHelper.ExpandImageNameWithTag("my-reg.com/library/nginx", "tag-3")
	assert.NoError(t, err)
	assert.Equal(t, "my-reg.com/library/nginx:tag-3", name)

	name, err = imageHelper.ExpandImageNameWithTag("my-reg.com/library/nginx:my-tag", "tag-4")
	assert.NoError(t, err)
	assert.Equal(t, "my-reg.com/library/nginx:tag-4", name)

	name, err = imageHelper.ExpandImageNameWithTag("my-reg.com/library/nginx", "-12@3%44-")
	assert.ErrorIs(t, err, imageHelper.ErrInvalidTag)
}

func TestSplitImageName(t *testing.T) {
	_, _, err := imageHelper.SplitImageName("nginx")
	assert.Error(t, err)

	name, tag, err := imageHelper.SplitImageName("docker.io/library/nginx:tag-2")
	assert.NoError(t, err)
	assert.Equal(t, "docker.io/library/nginx", name)
	assert.Equal(t, "tag-2", tag)

	name, tag, err = imageHelper.SplitImageName("my-reg.com/test/nginx:tag-3")
	assert.Equal(t, "my-reg.com/test/nginx", name)
	assert.NoError(t, err)
	assert.Equal(t, "tag-3", tag)

	name, tag, err = imageHelper.SplitImageName("my-reg.com/test/nginx")
	assert.Error(t, err)
}

func TestAuthConfigToBasicAuth(t *testing.T) {
	authConfig := "{\"username\":\"test-user\",\"password\":\"test-password\"}"
	encodedAuth := base64.URLEncoding.EncodeToString([]byte(authConfig))

	expectedBasicAuth := "test-user:test-password"
	expected := base64.URLEncoding.EncodeToString([]byte(expectedBasicAuth))

	basicAuth, err := imageHelper.AuthConfigToBasicAuth(encodedAuth)

	assert.NoError(t, err)
	assert.Equal(t, expected, basicAuth)
}

func TestParseDistribRefErr(t *testing.T) {
	_, err := imageHelper.ParseDistributionRef("invalid%image!123-name::")
	assert.NotNil(t, err, "invalid image name triggers image name parse error")
}

func TestParseDistribRef(t *testing.T) {
	name, err := imageHelper.ParseDistributionRef("nginx:latest")
	assert.Nil(t, err, "valid image name does not trigger any error, and it is expanded properly")
	assert.Equal(t, "docker.io/library/nginx:latest", name.String())
	name, err = imageHelper.ParseDistributionRef("nginx")
	assert.Nil(t, err, "valid image name does not trigger any error even without a tag, default tag is appended")
	assert.Equal(t, "docker.io/library/nginx:latest", name.String())
}
