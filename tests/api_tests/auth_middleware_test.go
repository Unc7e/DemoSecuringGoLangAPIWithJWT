package apitests

import (
	"DemoSecuringGOLangAPIWithJWT/core/authentication"
	"DemoSecuringGOLangAPIWithJWT/routers"
	"DemoSecuringGOLangAPIWithJWT/services"
	"DemoSecuringGOLangAPIWithJWT/settings"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type MiddlewaresTestSuite struct{}

var _ = Suite(&MiddlewaresTestSuite{})
var t *testing.T
var token string
var server *negroni.Negroni

func (s *MiddlewaresTestSuite) SetUpSuite(c *C) {
	os.Setenv("GO_ENV", "tests")
	settings.Init()
}

func (s *MiddlewaresTestSuite) SetUpTest(c *C) {
	authBackend := authentication.InitJWTAuthenticationBackend()
	assert.NotNil(t, authBackend)
	token, _ = authBackend.GenerateToken("1234")

	router := routers.InitRoutes()
	server = negroni.Classic()
	server.UseHandler(router)
}

func (s *MiddlewaresTestSuite) TestRequireTokenAuthentication(c *C) {
	resource := "/test/hello"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bear %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusOK)
}

func (s *MiddlewaresTestSuite) TestRequireTokenAuthenticationInvalidToken(c *C) {
	resource := "/test/hello"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", "token"))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func (s *MiddlewaresTestSuite) TestRequireTokenAuthenticationEmptyToken(c *C) {
	resource := "/test/hello"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", ""))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func (s *MiddlewaresTestSuite) TestRequireTokenAuthenticationWithoutToken(c *C) {
	resource := "/test/hello"

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}

func (s *MiddlewaresTestSuite) TestRequireTokenAuthenticationAfterLogout(c *C) {
	resource := "/test/hello"

	requestLogout, _ := http.NewRequest("GET", resource, nil)
	requestLogout.Header.Set("Authorization", fmt.Sprintf("Bearer v%", token))
	services.Logout(requestLogout)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", resource, nil)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
	server.ServeHTTP(response, request)

	assert.Equal(t, response.Code, http.StatusUnauthorized)
}
