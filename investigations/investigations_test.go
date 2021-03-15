package investigations

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/stretchr/testify/assert"
)

const (
	ApplicationJSON = "application/json"
	ContentType     = "Content-Type"
)

var (
	responses = struct {
		success string
		errors  string
	}{
		success: `{"error": [], "data": {"investigation": {"ID": "1234", "TenantID": "1", "Type": "Threat Hunt", "Priority" : 0}}}`,
		errors:  `{"errors": [{"message":"expected at least one definition, found }","locations":[{"line":4,"column":36}]}], "data": null}`,
	}
)

func TestDefaultFields(t *testing.T) {
	// sanity check to make sure people don't change default response fields willy-nilly

	// striping all space to prevent weird formating in comparison
	noSpace := regexp.MustCompile(`\s+`)
	expected := noSpace.ReplaceAllString(string(DefaultFields), "")
	actual := noSpace.ReplaceAllString(string(graphql.ResponseFields(`
                id
                created_at
                updated_at
                tenant_id
                description
                status
		key_findings
                created_by
                assignee_id
		genesis_alerts {
                        id
                }
                genesis_events {
                        id
                }
                alerts {
                        id
                }
                events {
                        id
                }
				priority
				type`)), "")

	assert.Equal(t, expected, actual)
}

func TestGetInvestigation_Success(t *testing.T) {
	investigationSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses.success))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err := investigationSvc.GetInvestigation(&GetInvestigationInput{ID: "1234", TenantID: "1"}, DefaultFields)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out.ID, "1234")
	assert.Equal(t, out.Priority, 0)
	assert.Equal(t, out.Type, "Threat Hunt")
}

func TestGetInvestigation_SuccessWithBearer(t *testing.T) {
	investigationSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer test-token")
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses.success))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err := investigationSvc.GetInvestigation(
		&GetInvestigationInput{ID: "1234", TenantID: "1"},
		DefaultFields,
		graphql.RequestWithToken("test-token"),
	)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, out.ID, "1234")
	assert.Equal(t, out.Priority, 0)
	assert.Equal(t, out.Type, "Threat Hunt")
}

func TestGetInvestigation_URLErrors(t *testing.T) {
	DefaultURL = "foo"
	tenantSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	out, err := tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestGetInvestigation_Errors(t *testing.T) {
	tenantSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses.errors))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err := tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, ``)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, "1 error occurred:\n\t* message: expected at least one definition, found }\n\n", err.Error())
}

func TestGetInvestigation_SuccessBadResponseBody(t *testing.T) {
	tenantSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("status=success"))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err := tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))
}

func TestGetInvestigation_ServerErrors(t *testing.T) {
	tenantSvc := NewInvestigationSvc(client.NewClient(), "test-app")
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(``))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err := tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, ``)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, "server responded with an error: 500", err.Error())

	fakeHandler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))

		r.Context().Done()
	}
	server = httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	server.Close()

	out, err = tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server connection error")

	fakeHandler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, ApplicationJSON, r.Header.Get(ContentType))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{''`))
	}

	server = httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	out, err = tenantSvc.GetInvestigation(&GetInvestigationInput{ID: "1234"}, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, "server responded with an error: 500", err.Error())
}
