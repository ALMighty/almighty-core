package login_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/fabric8-services/fabric8-wit/account"
	"github.com/fabric8-services/fabric8-wit/app"
	config "github.com/fabric8-services/fabric8-wit/configuration"
	"github.com/fabric8-services/fabric8-wit/gormapplication"
	"github.com/fabric8-services/fabric8-wit/gormtestsupport"
	. "github.com/fabric8-services/fabric8-wit/login"
	"github.com/fabric8-services/fabric8-wit/resource"
	"github.com/fabric8-services/fabric8-wit/test/token"
	"github.com/goadesign/goa"
	goajwt "github.com/goadesign/goa/middleware/security/jwt"
	"github.com/goadesign/goa/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
)

type serviceBlackBoxTest struct {
	gormtestsupport.DBTestSuite
	loginService  KeycloakOAuthService
	oauth         *oauth2.Config
	configuration *config.ConfigurationData
}

func TestRunServiceBlackBoxTest(t *testing.T) {
	suite.Run(t, &serviceBlackBoxTest{DBTestSuite: gormtestsupport.NewDBTestSuite("../config.yaml")})
}

// SetupSuite overrides the DBTestSuite's function but calls it before doing anything else
// The SetupSuite method will run before the tests in the suite are run.
// It sets up a database connection for all the tests in this suite without polluting global space.
func (s *serviceBlackBoxTest) SetupSuite() {
	s.DBTestSuite.SetupSuite()

	req := &http.Request{Host: "api.service.domain.org"}
	authEndpoint, err := s.Configuration.GetKeycloakEndpointAuth(req)
	if err != nil {
		panic(err)
	}
	tokenEndpoint, err := s.Configuration.GetKeycloakEndpointToken(req)
	if err != nil {
		panic(err)
	}
	s.oauth = &oauth2.Config{
		ClientID:     s.Configuration.GetKeycloakClientID(),
		ClientSecret: s.Configuration.GetKeycloakSecret(),
		Scopes:       []string{"user:email"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authEndpoint,
			TokenURL: tokenEndpoint,
		},
	}

	userRepository := account.NewUserRepository(s.DB)
	identityRepository := account.NewIdentityRepository(s.DB)
	app := gormapplication.NewGormDB(s.DB)
	s.loginService = NewKeycloakOAuthProvider(identityRepository, userRepository, token.TokenManager, app)
}

func (s *serviceBlackBoxTest) TestKeycloakAuthorizationRedirect() {
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize"),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	// The user clicks login while on WIT UI.
	// Therefore the referer would be an WIT URL.
	refererUrl := "https://wit-url.example.org/path"
	req.Header.Add("referer", refererUrl)

	prms := url.Values{}
	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	r := &http.Request{Host: "demo.api.openshift.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")

	assert.Equal(s.T(), 307, rw.Code)
	assert.Contains(s.T(), rw.Header().Get("Location"), s.oauth.Endpoint.AuthURL)
	assert.NotEqual(s.T(), rw.Header().Get("Location"), "")
}

func (s *serviceBlackBoxTest) getValidRedirectURLs() string {
	r := &http.Request{Host: "domain.com"}
	whitelist, err := s.Configuration.GetValidRedirectURLs(r)
	require.Nil(s.T(), err)
	return whitelist
}

func (s *serviceBlackBoxTest) TestKeycloakAuthorizationRedirectsToRedirectParam() {
	rw := httptest.NewRecorder()
	redirect := "https://url.example.org/pathredirect"
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize?redirect="),
	}
	parameters := url.Values{}
	if redirect != "" {
		parameters.Add("redirect", redirect)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, parameters)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	r := &http.Request{Host: "api.domain.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")

	assert.Equal(s.T(), 307, rw.Code)
	assert.Contains(s.T(), rw.Header().Get("Location"), s.oauth.Endpoint.AuthURL)
	assert.NotEqual(s.T(), rw.Header().Get("Location"), "")
}

func (s *serviceBlackBoxTest) TestKeycloakAuthorizationWithNoRefererAndRedirectParamFails() {
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize"),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	prms := url.Values{}
	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	r := &http.Request{Host: "api.domain.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")
	assert.Equal(s.T(), 400, rw.Code)
}

func (s *serviceBlackBoxTest) TestKeycloakAuthorizationWithNoValidRefererFails() {
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize"),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	// Not whitelisted redirect fails
	prms := url.Values{}
	prms.Add("redirect", "http://notauthorized.domain.com")

	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	r := &http.Request{Host: "api.domain.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", config.DefaultValidRedirectURLs, "")
	assert.Equal(s.T(), 400, rw.Code)

	// openshift.io redirects pass
	rw = httptest.NewRecorder()
	prms = url.Values{}
	prms.Add("redirect", "https://openshift.io/somepath")

	goaCtx = goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err = app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", config.DefaultValidRedirectURLs, "")
	assert.Equal(s.T(), 307, rw.Code)
	assert.Contains(s.T(), rw.Header().Get("Location"), s.oauth.Endpoint.AuthURL)
	assert.NotEqual(s.T(), rw.Header().Get("Location"), "")

	// Any redirects pass in Dev mode.
	rw = httptest.NewRecorder()
	prms = url.Values{}
	prms.Add("redirect", "https://anydoamin.io/somepath")

	goaCtx = goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err = app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	if err != nil {
		panic("invalid test data " + err.Error()) // bug
	}

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")
	assert.Equal(s.T(), 307, rw.Code)
	assert.Contains(s.T(), rw.Header().Get("Location"), s.oauth.Endpoint.AuthURL)
	assert.NotEqual(s.T(), rw.Header().Get("Location"), "")
}

func (s *serviceBlackBoxTest) TestKeycloakLinkRedirect() {
	keycloakLinkRedirect(s, "", "")
	keycloakLinkRedirect(s, "", "https://some.redirect.io")
	keycloakLinkRedirect(s, "github", "")
	keycloakLinkRedirect(s, "github", "https://some.redirect.io")
	keycloakLinkRedirect(s, "openshift-v3", "")
	keycloakLinkRedirect(s, "openshift-v3", "https://some.redirect.io")
}

func keycloakLinkRedirect(s *serviceBlackBoxTest, provider string, redirect string) {
	rw := httptest.NewRecorder()
	p := "/api/login/link"

	parameters := url.Values{}
	if redirect != "" {
		parameters.Add("redirect", redirect)
	}
	if provider != "" {
		parameters.Add("provider", provider)
	}

	req, err := http.NewRequest("GET", p, nil)
	require.Nil(s.T(), err)

	referrerUrl := "https://example.org/path"
	req.Header.Add("referer", referrerUrl)

	ss := uuid.NewV4().String()
	claims := jwt.MapClaims{}
	claims["session_state"] = &ss
	token := &jwt.Token{Claims: claims}
	ctx := goajwt.WithJWT(context.Background(), token)
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LinkTest"), rw, req, parameters)

	linkCtx, err := app.NewLinkLoginContext(goaCtx, req, goa.New("LoginService"))
	require.Nil(s.T(), err)

	r := &http.Request{Host: "api.example.org"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)
	clientID := s.Configuration.GetKeycloakClientID()

	err = s.loginService.Link(linkCtx, brokerEndpoint, clientID, s.getValidRedirectURLs())
	require.Nil(s.T(), err)

	assert.Equal(s.T(), 307, rw.Code)
	redirectLocation := rw.Header().Get("Location")
	if provider == "" {
		provider = "github"
		assert.Contains(s.T(), redirectLocation, "next%3Dopenshift-v3")
	} else {
		assert.NotContains(s.T(), redirectLocation, "next%3D")
	}
	location := brokerEndpoint + "/" + provider + "/link?"
	assert.Contains(s.T(), redirectLocation, location)
}

func keycloakLinkCallbackRedirect(s *serviceBlackBoxTest, next string) {
	rw := httptest.NewRecorder()
	p := "/api/login/linkcallback"

	parameters := url.Values{}
	parameters.Add("state", uuid.NewV4().String())
	parameters.Add("sessionState", uuid.NewV4().String())
	if next != "" {
		parameters.Add("next", next)
	}
	req, err := http.NewRequest("GET", p, nil)
	require.Nil(s.T(), err)

	referrerUrl := "https://sso.example.org/path"
	req.Header.Add("referer", referrerUrl)

	goaCtx := goa.NewContext(goa.WithAction(context.Background(), "LinkcallbackTest"), rw, req, parameters)

	linkCtx, err := app.NewLinkcallbackLoginContext(goaCtx, req, goa.New("LoginService"))
	require.Nil(s.T(), err)

	r := &http.Request{Host: "api.example.org"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)
	clientID := s.Configuration.GetKeycloakClientID()

	err = s.loginService.LinkCallback(linkCtx, brokerEndpoint, clientID)
	if next != "" {
		require.Nil(s.T(), err)
		assert.Equal(s.T(), 307, rw.Code)
		redirectLocation := rw.Header().Get("Location")
		assert.NotContains(s.T(), redirectLocation, "next%3D")
		location := brokerEndpoint + "/openshift-v3/link?"
		assert.Contains(s.T(), redirectLocation, location)
	} else {
		require.NotNil(s.T(), err)
	}
}

func (s *serviceBlackBoxTest) TestInvalidState() {
	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize"),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	// The OAuth 'state' is sent as a query parameter by calling /api/login/authorize?code=_SOME_CODE_&state=_SOME_STATE_
	// The request originates from Keycloak after a valid authorization by the end user.
	// This is not where the redirection should happen on failure.
	refererKeycloakUrl := "https://keycloak-url.example.org/path-of-login"
	req.Header.Add("referer", refererKeycloakUrl)

	prms := url.Values{
		"state": {},
		"code":  {"doesnt_matter_what_is_here"},
	}
	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	require.Nil(s.T(), err)

	r := &http.Request{Host: "demo.api.openshift.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")
	assert.Equal(s.T(), 401, rw.Code)
}

func (s *serviceBlackBoxTest) TestInvalidOAuthAuthorizationCode() {

	// When a valid referrer talks to our system and provides
	// an invalid OAuth2.0 code, the access token exchange
	// fails. In such a scenario, there is response redirection
	// to the valid referer, ie, the URL where the request originated from.
	// Currently, this should be something like https://demo.openshift.io/somepage/

	resource.Require(s.T(), resource.Database)

	// Setup request context
	rw := httptest.NewRecorder()
	u := &url.URL{
		Path: fmt.Sprintf("/api/login/authorize"),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		panic("invalid test " + err.Error()) // bug
	}

	// The user clicks login while on WIT UI.
	// Therefore the referer would be an WIT URL.
	refererUrl := "https://wit-url.example.org/path"
	req.Header.Add("referer", refererUrl)

	prms := url.Values{}
	ctx := context.Background()
	goaCtx := goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err := app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))
	require.Nil(s.T(), err)

	r := &http.Request{Host: "demo.api.openshift.io"}
	brokerEndpoint, err := s.Configuration.GetKeycloakEndpointBroker(r)
	require.Nil(s.T(), err)

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")

	assert.Equal(s.T(), 307, rw.Code) // redirect to keycloak login page.

	locationString := rw.HeaderMap["Location"][0]
	locationUrl, err := url.Parse(locationString)
	require.Nil(s.T(), err)

	allQueryParameters := locationUrl.Query()

	// Avoiding panics.
	assert.NotNil(s.T(), allQueryParameters)
	assert.NotNil(s.T(), allQueryParameters["state"][0])

	returnedState := allQueryParameters["state"][0]

	prms = url.Values{
		"state": {returnedState},
		"code":  {"INVALID_OAUTH2.0_CODE"},
	}
	ctx = context.Background()
	rw = httptest.NewRecorder()

	req, err = http.NewRequest("GET", u.String(), nil)

	// The OAuth code is sent as a query parameter by calling /api/login/authorize?code=_SOME_CODE_&state=_SOME_STATE_
	// The request originates from Keycloak after a valid authorization by the end user.
	// This is not where the redirection should happen on failure.
	refererKeycloakUrl := "https://keycloak-url.example.org/path-of-login"
	req.Header.Add("referer", refererKeycloakUrl)
	require.Nil(s.T(), err)

	goaCtx = goa.NewContext(goa.WithAction(ctx, "LoginTest"), rw, req, prms)
	authorizeCtx, err = app.NewAuthorizeLoginContext(goaCtx, req, goa.New("LoginService"))

	err = s.loginService.Perform(authorizeCtx, s.oauth, brokerEndpoint, "", "", s.getValidRedirectURLs(), "")

	locationString = rw.HeaderMap["Location"][0]
	locationUrl, err = url.Parse(locationString)
	require.Nil(s.T(), err)

	allQueryParameters = locationUrl.Query()
	assert.Equal(s.T(), 307, rw.Code) // redirect to WIT page where login was clicked.
	// Avoiding panics.
	assert.NotNil(s.T(), allQueryParameters)
	assert.NotNil(s.T(), allQueryParameters["error"])
	assert.NotEqual(s.T(), allQueryParameters["error"][0], "")

	returnedErrorReason := allQueryParameters["error"][0]
	assert.NotEmpty(s.T(), returnedErrorReason)
	assert.NotContains(s.T(), locationString, refererKeycloakUrl)
	assert.Contains(s.T(), locationString, refererUrl)
}
