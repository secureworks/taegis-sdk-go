package users

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	testTenantID = "123456789"
	testRole     = "scwxAnalyst"
	testUserID   = "1111"
	testEmail    = "noone@nowhere.com"
	testStatus   = "Registered"
)

var (
	uSvc             = NewUserSvc(client.NewClient())
	getUserResponses = map[string]string{
		"success": `{"errors": [], "data": {"tdruser": {"id": "1111", "user_id": "my_user_id", "email": "noone@nowhere.com", "status": "Registered", "roles": ["admin"], "tenants": [{"id": "5555"}], "tenants_v2": [{"id": "5555", "role": "admin"}, {"id": "123456789", "role": "tm-manager"}], "given_name": null, "family_name": null, "eula": {"date": null, "version": ""}, "created_at": "2017-11-30T20:49:07Z", "updated_at": "2020-06-23T18:10:58Z", "last_login": "2020-06-23T18:10:58Z"}}}`,
		"errors":  `{"errors": [{"message":"Field \"eula\" of type \"TDRUserLicense\" must have a selection of subfields. Did you mean \"eula { ... }\"?","locations":[{"line":16,"column":3}]}], "data": null}`,
	}
	findUsersResponses = map[string]string{
		"success": `{"errors": [], "data": {"tdrusers": [{"id": "1111", "user_id": "my_user_id", "email": "noone@nowhere.com", "status": "Registered", "roles": ["admin"], "tenants": [{"id": "5555"}], "tenants_v2": [{"id": "5555", "role": "admin"}, {"id": "123456789", "role": "tm-manager"}], "given_name": null, "family_name": null, "eula": {"date": null, "version": ""}, "created_at": "2017-11-30T20:49:07Z", "updated_at": "2020-06-23T18:10:58Z", "last_login": "2020-06-23T18:10:58Z"}]}}`,
		"errors":  `{"errors": [{"message":"Field \"eula\" of type \"TDRUserLicense\" must have a selection of subfields. Did you mean \"eula { ... }\"?","locations":[{"line":16,"column":3}]}], "data": null}`,
	}
)

func TestGetUserMissingAuth(t *testing.T) {
	envy.Set(UserServiceEnv, DefaultURL)
	getUserInput := &GetUserInput{}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "401 Unauthorized")
}

func TestGetUserMissingTenantContext(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	getUserInput := &GetUserInput{}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "tenant context is required")
}

func TestGetUserInvalidUserSvcURL(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	envy.Set(UserServiceEnv, "http://192.168.0.%31")
	getUserInput := &GetUserInput{}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid URL")
}

func TestGetUserBadUserServer(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	envy.Set(UserServiceEnv, "http://localhost:999999")
	getUserInput := &GetUserInput{}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid port")
}

func TestGetUserUserServerReturnsBadResponseCode(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	getUserInput := &GetUserInput{ID: testTenantID}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("error: %d", http.StatusUnprocessableEntity))
}

func TestGetUserBadResponseBody(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("status=success"))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	getUserInput := &GetUserInput{ID: testTenantID}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "error decoding get user response")
}

func TestGetUserErrorResponse(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(getUserResponses["errors"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	getUserInput := &GetUserInput{ID: testTenantID}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "1 error occurred")
	require.Contains(t, err.Error(), "must have a selection of subfields")
}

func TestGetUserHappyPath(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(getUserResponses["success"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	getUserInput := &GetUserInput{ID: testTenantID}
	out, err := uSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.NotNil(t, out)
	require.Nil(t, err)
	require.Equal(t, "1111", out.ID)
}

func TestGetUserHappyPathOauth(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	clientID := "testClientID"
	clientSecret := "testClientSecret"
	basicHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	mockAuthService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Basic "+basicHeader, r.Header.Get(common.AuthorizationHeader))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		_, err = w.Write([]byte("access_token=" + accessToken + "&token_type=bearer&expires_in=3600"))
		require.Nil(t, err)
	}))
	defer mockAuthService.Close()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     mockAuthService.URL,
	}

	ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, cleanhttp.DefaultClient())
	httpClient := config.Client(ctx)
	sdkClient := client.NewClient(client.WithHTTPClient(httpClient))
	userSvc := NewUserSvc(sdkClient)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(getUserResponses["success"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	getUserInput := &GetUserInput{ID: testTenantID}
	out, err := userSvc.GetUser(getUserInput, DefaultFields, graphql.RequestWithTenant(testTenantID))
	require.NotNil(t, out)
	require.Nil(t, err)
	require.Equal(t, "1111", out.ID)
}

func TestFindUsersMissingAuth(t *testing.T) {
	envy.Set(UserServiceEnv, DefaultURL)
	findUsersInput := &FindUsersInput{}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "401 Unauthorized")
}

func TestFindUsersMissingTenantContext(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	findUsersInput := &FindUsersInput{}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "tenant context is required")
}

func TestFindUsersInvalidUserSvcURL(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	envy.Set(UserServiceEnv, "http://192.168.0.%31")
	findUsersInput := &FindUsersInput{}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid URL")
}

func TestFindUsersBadUserServer(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	envy.Set(UserServiceEnv, "http://localhost:999999")
	findUsersInput := &FindUsersInput{}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "invalid port")
}

func TestFindUsersUserServerReturnsBadResponseCode(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusUnprocessableEntity)
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{Email: testEmail}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("error: %d", http.StatusUnprocessableEntity))
}

func TestFindUsersBadResponseBody(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("status=success"))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{Email: testEmail}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "error decoding find users response")
}

func TestFindUsersErrorResponse(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(findUsersResponses["errors"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{Email: testEmail}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.Nil(t, out)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "1 error occurred")
	require.Contains(t, err.Error(), "must have a selection of subfields")
}

func TestFindUsersHappyPath(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(findUsersResponses["success"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{Email: testEmail}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.NotNil(t, out)
	require.Nil(t, err)
	require.Equal(t, 1, len(out))
	require.Equal(t, "1111", out[0].ID)
}

func TestFindUsersHappyPathOauth(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	clientID := "testClientID"
	clientSecret := "testClientSecret"
	basicHeader := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	mockAuthService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Basic "+basicHeader, r.Header.Get(common.AuthorizationHeader))

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
		_, err = w.Write([]byte("access_token=" + accessToken + "&token_type=bearer&expires_in=3600"))
		require.Nil(t, err)
	}))
	defer mockAuthService.Close()

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     mockAuthService.URL,
	}

	ctx := context.WithValue(context.TODO(), oauth2.HTTPClient, cleanhttp.DefaultClient())
	httpClient := config.Client(ctx)
	sdkClient := client.NewClient(client.WithHTTPClient(httpClient))
	userSvc := NewUserSvc(sdkClient)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(findUsersResponses["success"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{Email: testEmail}
	out, err := userSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithTenant(testTenantID))
	require.NotNil(t, out)
	require.Nil(t, err)
	require.Equal(t, 1, len(out))
	require.Equal(t, "1111", out[0].ID)
}

func TestFindUserFullInput(t *testing.T) {
	accessToken, err := createTestToken(testUserID, testTenantID, testRole)
	require.Nil(t, err)
	require.NotNil(t, accessToken)

	mockHandler := func() http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			headers := r.Header
			require.Equal(t, "Bearer "+accessToken, headers.Get(common.AuthorizationHeader))
			require.Equal(t, testTenantID, headers.Get(common.XTenantContextHeader))
			require.Equal(t, http.MethodPost, r.Method)

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(findUsersResponses["success"]))
		}
	}
	mockUserService := getMockUsersAPIServer(mockHandler)
	defer mockUserService.Close()

	findUsersInput := &FindUsersInput{
		Email:    testEmail,
		Role:     testRole,
		TenantID: testTenantID,
		Status:   testStatus,
		Page:     1,
		PerPage:  10,
	}
	out, err := uSvc.FindUsers(findUsersInput, DefaultFields, graphql.RequestWithToken(accessToken), graphql.RequestWithTenant(testTenantID))
	require.NotNil(t, out)
	require.Nil(t, err)
	require.Equal(t, 1, len(out))
	require.Equal(t, "1111", out[0].ID)
}

func createTestToken(userID string, tenantID string, role string) (string, error) {
	claims := jwt.MapClaims{}
	claims["sub"] = userID
	claims["https://missione/octolabs/io/tenantIds"] = tenantID
	claims["https://missione/octolabs/io/roles"] = role
	claims["https://missione/octolabs/io/tenant_v2"] = tenantID + ":" + role
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ss, err := token.SignedString([]byte("ctpx"))
	if err != nil {
		return "", fmt.Errorf("error creating access token: %s", err)
	}

	return ss, nil
}

func getMockUsersAPIServer(handler func() http.HandlerFunc) *httptest.Server {
	newServer := httptest.NewServer(handler())

	envy.Set(UserServiceEnv, newServer.URL)
	return newServer
}
