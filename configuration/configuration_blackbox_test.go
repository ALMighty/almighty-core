package configuration

import (
	"fmt"
	"os"
	"testing"

	"net/http"

	"github.com/almighty/almighty-core/resource"
	"github.com/goadesign/goa"
	"github.com/stretchr/testify/assert"
)

var req *goa.RequestData

func TestMain(m *testing.M) {
	resetConfiguration()

	req = &goa.RequestData{
		Request: &http.Request{Host: "api.service.domain.org"},
	}
	os.Exit(m.Run())
}

func resetConfiguration() {
	if err := Setup(""); err != nil {
		panic(fmt.Errorf("Failed to setup the configuration: %s", err.Error()))
	}
}

func TestGetKeycloakEndpointAuthDevModeOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	url, err := GetKeycloakEndpointAuth(req)
	assert.Nil(t, err)
	// In dev mode it's always the defualt value regardless of the request
	assert.Equal(t, "http://sso.demo.almighty.io/auth/realms/demo/protocol/openid-connect/auth", url)
}

func TestGetKeycloakEndpointAuthSetByEnvVaribaleOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	env := os.Getenv("ALMIGHTY_KEYCLOAK_ENDPOINT_AUTH")
	defer func() {
		os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_AUTH", env)
		resetConfiguration()
	}()

	os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_AUTH", "authEndpoint")
	resetConfiguration()

	url, err := GetKeycloakEndpointAuth(req)
	assert.Nil(t, err)
	assert.Equal(t, "authEndpoint", url)
}

func TestGetKeycloakEndpointTokenOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	url, err := GetKeycloakEndpointToken(req)
	assert.Nil(t, err)
	// In dev mode it's always the defualt value regardless of the request
	assert.Equal(t, "http://sso.demo.almighty.io/auth/realms/demo/protocol/openid-connect/token", url)
}

func TestGetKeycloakEndpointTokenSetByEnvVaribaleOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	env := os.Getenv("ALMIGHTY_KEYCLOAK_ENDPOINT_TOKEN")
	defer func() {
		os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_TOKEN", env)
		resetConfiguration()
	}()

	os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_TOKEN", "tokenEndpoint")
	resetConfiguration()

	url, err := GetKeycloakEndpointToken(req)
	assert.Nil(t, err)
	assert.Equal(t, "tokenEndpoint", url)
}

func TestGetKeycloakEndpointUserInfoOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	t.Parallel()

	url, err := GetKeycloakEndpointUserInfo(req)
	assert.Nil(t, err)
	// In dev mode it's always the defualt value regardless of the request
	assert.Equal(t, "http://sso.demo.almighty.io/auth/realms/demo/protocol/openid-connect/userinfo", url)
}

func TestGetKeycloakEndpointUserInfoSetByEnvVaribaleOK(t *testing.T) {
	resource.Require(t, resource.UnitTest)
	env := os.Getenv("ALMIGHTY_KEYCLOAK_ENDPOINT_USERINFO")
	defer func() {
		os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_USERINFO", env)
		resetConfiguration()
	}()

	os.Setenv("ALMIGHTY_KEYCLOAK_ENDPOINT_USERINFO", "userinfoEndpoint")
	resetConfiguration()

	url, err := GetKeycloakEndpointUserInfo(req)
	assert.Nil(t, err)
	assert.Equal(t, "userinfoEndpoint", url)
}
