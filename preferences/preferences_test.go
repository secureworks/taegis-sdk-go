package preferences

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

const (
	preferencesURL = "PREFERENCES_URL"
)

var (
	responses = map[string]string{
		"createSuccess":                    `{"errors": [], "data": {"createUserPreference": {"id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}}}`,
		"createTenantSuccess":              `{"errors": [], "data": {"createTenantPreference": {"id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}}}`,
		"getSuccess":                       `{"errors": [], "data": {"userPreferenceByKey": {"id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}}}`,
		"getNotificationPreferenceSuccess": `{"errors": [], "data": {"userNotificationPreference": {"id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}}}`,
		"getTenantSuccess":                 `{"errors": [], "data": {"listTenantPreferencesByKey": [{"id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}]}}`,
		"error":                            `{"errors": [{"message": "json body could not be decoded: invalid character ']' after object key:value pair"}],"data": null}`,
		"badTokenError":                    `{"errors": [{"message": "no authorization header", "path": ["createUserPreference"]}], "data": {"createUserPreference": null}}`,
		"JSONDecodeError":                  `{"errors": [], "data": {"createUserPreference": "id": "80853652-938f-41c0-9ec5-36aebbe9d17a", "user_id": "auth0|123", "key": "email", "preference_items": [{"key": "Mention", "value": "true"}]}}}`,
	}
)

func TestCreatePreferences_Success(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["createSuccess"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	userEmailPreferences := &UserEmail{
		Mention:             true,
		AssigneeChange:      true,
		AwaitingAction:      true,
		GroupMention:        true,
		GroupAssigneeChange: true,
		GroupAwaitingAction: true,
	}

	preferenceStorage := make(map[string]interface{})
	pStorageBytes, err := json.Marshal(userEmailPreferences)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	err = json.Unmarshal(pStorageBytes, &preferenceStorage)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	preferenceInput := &PreferencesInput{
		TenantID:    "123456789",
		Key:         EmailKey,
		Preferences: preferenceStorage,
	}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, fmt.Sprintf("%v", out.ID), "80853652-938f-41c0-9ec5-36aebbe9d17a")
	assert.Equal(t, out.UserID, "auth0|123")
	assert.Equal(t, out.PreferenceItems[0].Value, "true")
	assert.Nil(t, err)
}

func TestCreateTenantPreferences_Success(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["createTenantSuccess"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	userEmailPreferences := &UserEmail{
		Mention:             true,
		AssigneeChange:      true,
		AwaitingAction:      true,
		GroupMention:        true,
		GroupAssigneeChange: true,
		GroupAwaitingAction: true,
	}

	preferenceStorage := make(map[string]interface{})
	pStorageBytes, err := json.Marshal(userEmailPreferences)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	err = json.Unmarshal(pStorageBytes, &preferenceStorage)
	if err != nil {
		spew.Dump(err)
		os.Exit(1)
	}

	preferenceInput := &PreferencesInput{
		TenantID:    "123456789",
		Key:         EmailKey,
		Preferences: preferenceStorage,
	}

	out, err := preferenceSvc.CreateTenantPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, err)
	assert.NotNil(t, out)
	assert.Equal(t, fmt.Sprintf("%v", out.ID), "80853652-938f-41c0-9ec5-36aebbe9d17a")
	assert.Equal(t, out.UserID, "auth0|123")
	assert.Equal(t, out.PreferenceItems[0].Value, "true")
	assert.Nil(t, err)
}

func TestGetPreferences_Success(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["getSuccess"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{
		BearerToken: os.Getenv("ACCESS_TOKEN"),
		TenantID:    "123456789",
		Key:         EmailKey,
	}

	out, err := preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.NotNil(t, out)
	assert.Equal(t, fmt.Sprintf("%v", out.ID), "80853652-938f-41c0-9ec5-36aebbe9d17a")
	assert.Equal(t, out.UserID, "auth0|123")
	assert.Equal(t, out.PreferenceItems[0].Value, "true")
	assert.Nil(t, err)
}

func TestGetNotificationPreferences_Success(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["getNotificationPreferenceSuccess"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{
		BearerToken: os.Getenv("ACCESS_TOKEN"),
		TenantID:    "123456789",
		Key:         EmailKey,
		UserID:      "auth0|123",
	}

	out, err := preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.NotNil(t, out)
	assert.Equal(t, fmt.Sprintf("%v", out.ID), "80853652-938f-41c0-9ec5-36aebbe9d17a")
	assert.Equal(t, out.UserID, "auth0|123")
	assert.Equal(t, out.PreferenceItems[0].Value, "true")
	assert.Nil(t, err)
}

func TestListTenantPreferences_Success(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["getTenantSuccess"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{
		BearerToken: os.Getenv("ACCESS_TOKEN"),
		TenantID:    "123456789",
		Key:         EmailKey,
	}

	out, err := preferenceSvc.ListTenantPreferencesByKey(preferenceInput, DefaultFields)

	assert.NotNil(t, out)
	assert.Equal(t, fmt.Sprintf("%v", out[0].ID), "80853652-938f-41c0-9ec5-36aebbe9d17a")
	assert.Equal(t, out[0].UserID, "auth0|123")
	assert.Equal(t, out[0].PreferenceItems[0].Value, "true")
	assert.Nil(t, err)
}

func TestPreferences_Failure(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["error"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "json body could not be decoded: invalid character"))

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "json body could not be decoded: invalid character"))

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "json body could not be decoded: invalid character"))
}

func TestPreferences_BadTokenFailure(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["badTokenError"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "no authorization header"))

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "no authorization header"))

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "no authorization header"))
}

func TestPreferences_JSONDecodeError(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responses["JSONDecodeError"]))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))
}

func TestPreferences_SuccessBadResponseBody(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("status=createSuccess"))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "error decoding response"))
}

func Test_URLErrors(t *testing.T) {
	DefaultURL = "foo"
	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}

func TestPreferences_ServerErrors(t *testing.T) {
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(``))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt := client.NewClient()
	preferenceSvc := NewPreferencesSvc(clnt, "preferences example")

	preferenceInput := &PreferencesInput{}

	out, err := preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with an error: 500 - response body: ")

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with an error: 500 - response body: ")

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	expected := "2 errors occurred:\n\t* ctpx-sdk-go/graphql: " +
		"server responded with an error: 500\n\t* ctpx-sdk-go/graphql: error decoding response: EOF\n\n"
	assert.Equal(t, expected, err.Error())

	fakeHandler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		r.Context().Done()
	}
	server = httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	server.Close()

	clnt = client.NewClient()
	preferenceSvc = NewPreferencesSvc(clnt, "preferences example")

	preferenceInput = &PreferencesInput{}

	out, err = preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server connection error")

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server connection error")

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)

	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server connection error")

	fakeHandler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{''}`))
	}

	server = httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt = client.NewClient()
	preferenceSvc = NewPreferencesSvc(clnt, "preferences example")

	preferenceInput = &PreferencesInput{}

	out, err = preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with an error: 500 - response body: {''}")

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with an error: 500 - response body: {''}")

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	expected = "2 errors occurred:\n\t* ctpx-sdk-go/graphql: server responded with an error: 500\n\t* " +
		"ctpx-sdk-go/graphql: error decoding response: invalid character '\\'' looking for beginning of object key string\n\n"
	assert.Equal(t, expected, err.Error())

	fakeHandler = func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{''}`))
	}

	server = httptest.NewServer(http.HandlerFunc(fakeHandler))
	DefaultURL = server.URL
	defer server.Close()

	clnt = client.NewClient()
	preferenceSvc = NewPreferencesSvc(clnt, "preferences example")

	preferenceInput = &PreferencesInput{}

	out, err = preferenceSvc.CreatePreferences(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with bad request: 400 - response body: {''}")

	out, err = preferenceSvc.GetPreferencesByKey(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "ctpx-sdk-go/preferences: server responded with bad request: 400 - response body: {''}")

	out, err = preferenceSvc.GetNotificationPreferences(preferenceInput, DefaultFields)
	assert.Nil(t, out)
	assert.NotNil(t, err)
	expected = "2 errors occurred:\n\t* ctpx-sdk-go/graphql: server responded with an error: 400\n\t* " +
		"ctpx-sdk-go/graphql: error decoding response: invalid character '\\'' looking for beginning of object key string\n\n"
	assert.Equal(t, expected, err.Error())
}

func TestPreferences_NewRequestError(t *testing.T) {
	r := strings.NewReader("my request")
	requestError := NewRequestError(r, 0)

	assert.NotNil(t, requestError)
}
