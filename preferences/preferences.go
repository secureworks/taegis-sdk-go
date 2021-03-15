package preferences

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/secureworks/taegis-sdk-go/client"
	"github.com/secureworks/taegis-sdk-go/common"
	"github.com/secureworks/taegis-sdk-go/graphql"
	"github.com/gobuffalo/envy"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-multierror"
)

var (
	_          IPreferencesSvc = &PreferencesSvc{}
	DefaultURL                 = "https://api.ctpx.secureworks.com/graphql"
)

const (
	DefaultFields graphql.ResponseFields = `
        id
        created_at
        updated_at
        user_id
        key
        preference_items {key, value}
        `
	EmailKey = "email"
)

type UserEmail struct {
	Mention             bool `json:"mention"`
	AssigneeChange      bool `json:"assignee_change"`
	AwaitingAction      bool `json:"awaiting_action"`
	GroupMention        bool `json:"group_mention"`
	GroupAssigneeChange bool `json:"group_assignee_change"`
	GroupAwaitingAction bool `json:"group_awaiting_action"`
	ExportsCSVReady     bool `json:"exports_csv_ready"`
}

// IPreferencesSvc defines what the the Preferences API can do
type IPreferencesSvc interface {
	CreatePreferences(*PreferencesInput, graphql.ResponseFields) (*PreferencesOutput, error)
	CreateTenantPreferences(in *PreferencesInput, rf graphql.ResponseFields) (*PreferencesOutput, error)
	GetPreferencesByKey(*PreferencesInput, graphql.ResponseFields) (*PreferencesOutput, error)
	ListTenantPreferencesByKey(in *PreferencesInput, rf graphql.ResponseFields) ([]*PreferencesOutput, error)
}

// PreferencesSvc is the concrete implementation of the interface against the real api
type PreferencesSvc struct {
	client      *client.Client
	serviceName string
}

// NewPreferenceSvc takes a client from `client` package -- see examples/preferences.go for an example
func NewPreferencesSvc(c *client.Client, serviceName string) *PreferencesSvc {
	return &PreferencesSvc{client: c, serviceName: serviceName}
}

//NewRequestError takes a response body and a response code and returns a custom error
//to the CreateNotification caller for failed request
func NewRequestError(respBody io.Reader, respCode int) *RequestError {
	//reading 4kb of data only
	respBody = io.LimitReader(respBody, 1024*4)
	bodyBytes, err := ioutil.ReadAll(respBody)

	if err != nil {
		bodyString := fmt.Sprintf("ctpx-sdk-go/preferences: server responded with an error: %d - response body is invalid: %e", respCode, err)
		return &RequestError{
			RespBody: bodyString,
			Code:     http.StatusInternalServerError,
		}
	}

	bodyString := string(bodyBytes)
	return &RequestError{
		RespBody: bodyString,
		Code:     respCode,
	}
}

//RequestError is a custom error type to capture a request response body and a response code
type RequestError struct {
	RespBody string
	Code     int
}

func (b *RequestError) Error() string {
	if b.Code >= http.StatusInternalServerError {
		return fmt.Sprintf("ctpx-sdk-go/preferences: server responded with an error: %d - response body: %s", b.Code, b.RespBody)
	}
	return fmt.Sprintf("ctpx-sdk-go/preferences: server responded with bad request: %d - response body: %s", b.Code, b.RespBody)
}

func (t *PreferencesSvc) CreatePreferences(in *PreferencesInput, rf graphql.ResponseFields) (*PreferencesOutput, error) {
	var preferenceItems []PreferenceItem
	for k, v := range in.Preferences {
		item := PreferenceItem{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		}
		preferenceItems = append(preferenceItems, item)
	}

	query := fmt.Sprintf(`
        mutation createUserPreference ($newUserPreference: NewUserPreferenceInput) {
            createUserPreference (newUserPreference: $newUserPreference) 
            {%s}
        }`, rf)

	graphqlReq := graphql.NewRequest(query)
	newUserPreference := map[string]interface{}{
		"key":              in.Key,
		"preference_items": preferenceItems,
	}
	graphqlReq.Var("newUserPreference", newUserPreference)

	request, err := BuildRequest(graphqlReq, in.BearerToken, in.TenantID)

	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("ctpx-sdk-go/preferences: server connection error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		requestError := NewRequestError(resp.Body, resp.StatusCode)
		return nil, requestError
	}

	type createResponse struct {
		Data struct {
			Out *PreferencesOutput `json:"createUserPreference"`
		} `json:"data"`
		Error []graphql.Error `json:"errors"`
	}

	out := createResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	if len(out.Error) > 0 {
		var outErr error
		for _, e := range out.Error {
			outErr = multierror.Append(outErr, e)
		}
		return nil, outErr
	}

	return out.Data.Out, nil
}

func (t *PreferencesSvc) CreateTenantPreferences(in *PreferencesInput, rf graphql.ResponseFields) (*PreferencesOutput, error) {
	var preferenceItems []PreferenceItem
	for k, v := range in.Preferences {
		item := PreferenceItem{
			Key:   k,
			Value: fmt.Sprintf("%v", v),
		}
		preferenceItems = append(preferenceItems, item)
	}

	query := fmt.Sprintf(`
        mutation createTenantPreference ($newTenantPreference: NewTenantPrefernceInput) {
            createTenantPreference (newTenantPreference: $newTenantPreference) 
            {%s}
        }`, rf)

	graphqlReq := graphql.NewRequest(query)
	newUserPreference := map[string]interface{}{
		"key":              in.Key,
		"preference_items": preferenceItems,
	}
	graphqlReq.Var("newTenantPreference", newUserPreference)

	request, err := BuildRequest(graphqlReq, in.BearerToken, in.TenantID)

	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("ctpx-sdk-go/preferences: server connection error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		requestError := NewRequestError(resp.Body, resp.StatusCode)
		return nil, requestError
	}

	type createResponse struct {
		Data struct {
			Out *PreferencesOutput `json:"createTenantPreference"`
		} `json:"data"`
		Error []graphql.Error `json:"errors"`
	}

	out := createResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	if len(out.Error) > 0 {
		var outErr error
		for _, e := range out.Error {
			outErr = multierror.Append(outErr, e)
		}
		return nil, outErr
	}

	return out.Data.Out, nil
}

func (t *PreferencesSvc) ListTenantPreferencesByKey(in *PreferencesInput, rf graphql.ResponseFields) ([]*PreferencesOutput, error) {
	query := `query listTenantPreferencesByKey($key: String!) {
			listTenantPreferencesByKey(key: $key){
				preference_items{
					key
					value
				}
				tenant_id
			}
		}`
	graphqlReq := graphql.NewRequest(query)
	graphqlReq.Var("key", in.Key)

	request, err := BuildRequest(graphqlReq, in.BearerToken, in.TenantID)

	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("ctpx-sdk-go/preferences: server connection error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		requestError := NewRequestError(resp.Body, resp.StatusCode)
		return nil, requestError
	}

	type createResponse struct {
		Data struct {
			Out []*PreferencesOutput `json:"listTenantPreferencesByKey"`
		} `json:"data"`
		Error []graphql.Error `json:"errors"`
	}

	out := createResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	if len(out.Error) > 0 {
		var outErr error
		for _, e := range out.Error {
			outErr = multierror.Append(outErr, e)
		}
		return nil, outErr
	}

	return out.Data.Out, nil
}

func (t *PreferencesSvc) GetPreferencesByKey(in *PreferencesInput, rf graphql.ResponseFields) (*PreferencesOutput, error) {
	query := fmt.Sprintf(`
        query userPreferenceByKey ($key: String!) {
            userPreferenceByKey (key: $key) 
            {%s}
        }`, rf)

	graphqlReq := graphql.NewRequest(query)
	graphqlReq.Var("key", in.Key)

	request, err := BuildRequest(graphqlReq, in.BearerToken, in.TenantID)

	if err != nil {
		return nil, err
	}

	resp, err := t.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("ctpx-sdk-go/preferences: server connection error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		requestError := NewRequestError(resp.Body, resp.StatusCode)
		return nil, requestError
	}

	type createResponse struct {
		Data struct {
			Out *PreferencesOutput `json:"userPreferenceByKey"`
		} `json:"data"`
		Error []graphql.Error `json:"errors"`
	}

	out := createResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("error decoding response: %s", err.Error())
	}

	if len(out.Error) > 0 {
		var outErr error
		for _, e := range out.Error {
			outErr = multierror.Append(outErr, e)
		}
		return nil, outErr
	}

	return out.Data.Out, nil
}

func (t *PreferencesSvc) GetNotificationPreferences(in *PreferencesInput, rf graphql.ResponseFields) (*PreferencesOutput, error) {
	query := fmt.Sprintf(`
        query userNotificationPreference ($userID: String!) {
            userNotificationPreference (userID: $userID) 
            {%s}
        }`, rf)

	graphqlReq := graphql.NewRequest(query)
	graphqlReq.Var("userID", in.UserID)

	h := http.Header{}
	h.Add(common.AuthorizationHeader, "Bearer "+in.BearerToken)
	h.Add(common.XTenantContextHeader, in.TenantID)
	h.Add("Content-Type", "application/json")

	out := &struct {
		Out *PreferencesOutput `json:"userNotificationPreference"`
	}{}

	qc := &graphql.QueryConfig{
		ServerURL:  envy.Get("PREFERENCES_URL", DefaultURL),
		HClient:    t.client,
		Request:    graphqlReq,
		Header:     h,
		EscapeHTML: false,
		Output:     out,
	}

	err := graphql.ExecuteQueryContext(context.Background(), qc)

	if err != nil {
		return nil, err
	}

	return out.Out, nil
}

func BuildRequest(graphqlReq *graphql.Request, bearerToken string, tenantID string) (*http.Request, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(graphqlReq)

	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(buf.Bytes())
	preferencesURL := envy.Get("PREFERENCES_URL", DefaultURL)
	request, err := http.NewRequest(http.MethodPost, preferencesURL, reqBody)

	if err != nil {
		return nil, fmt.Errorf("ctpx-sdk-go/preferences: malformed request error: %w", err)
	}

	request.Header.Add(common.AuthorizationHeader, "Bearer "+bearerToken)
	request.Header.Add(common.XTenantContextHeader, tenantID)
	request.Header.Add("Content-Type", "application/json")

	return request, nil
}

type PreferencesInput struct {
	BearerToken string `json:"-"`
	UserID      string
	TenantID    string
	Key         string
	Preferences map[string]interface{}
}

type PreferencesOutput struct {
	ID              uuid.UUID        `json:"id" db:"id"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at" db:"updated_at"`
	UserID          string           `json:"user_id" db:"user_id"`     //only used for userPreferences
	TenantID        string           `json:"tenant_id" db:"tenant_id"` // only used for tenantPreferences
	Key             string           `json:"key" db:"key"`
	PreferenceItems []PreferenceItem `json:"preference_items"`
}

type PreferenceItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
