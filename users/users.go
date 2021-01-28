package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/nulls"
	"github.com/hashicorp/go-multierror"
)

const (
	UserServiceEnv = "USER_SVC_URL"
	// DefaultURL server URL for the API
	DefaultURL = "https://api.ctpx.secureworks.com/graphql"
	// DefaultFields that the API can respond with, can be permanently set at the package
	// or used as an example to define your own
	DefaultFields graphql.ResponseFields = `
		id
		user_id
		email
		status
		roles
		tenants {
			id
		}
		tenants_v2 {
			id
			role
		}
		family_name
		given_name
		phone_number
		created_at
		updated_at
		last_login
		eula {
			date
			version
		}
		timezone
	`
)

var (
	_ IUserSvc = &UserSvc{}
)

type GetUserInput struct {
	ID string
}

type FindUsersInput struct {
	Email    string
	Role     string
	TenantID string
	Status   string
	Page     int
	PerPage  int
}

type TenantOutput struct {
	ID string `json:"id"`
}

type TenantV2Output struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

type EulaOutput struct {
	Date    time.Time `json:"date"`
	Version string    `json:"version"`
}

type UserOutput struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	LastLogin       nulls.Time       `json:"last_login"`
	Status          string           `json:"status"`
	Email           string           `json:"email"`
	EmailNormalized string           `json:"email_normalized"`
	FamilyName      nulls.String     `json:"family_name"`
	GivenName       nulls.String     `json:"given_name"`
	PhoneNumber     nulls.String     `json:"phone_number"`
	Roles           []string         `json:"roles"`
	Tenants         []TenantOutput   `json:"tenants"`
	TenantsV2       []TenantV2Output `json:"tenants_v2"`
	Eula            *EulaOutput      `json:"eula,omitempty"`
	Timezone        nulls.String     `json:"timezone"`
}

type IUserSvc interface {
	GetUser(*GetUserInput, graphql.ResponseFields, ...graphql.RequestOption) (*UserOutput, error)
	FindUsers(*FindUsersInput, graphql.ResponseFields, ...graphql.RequestOption) ([]*UserOutput, error)
}

type UserSvc struct {
	client *client.Client
}

type getUserResponse struct {
	Data struct {
		Out *UserOutput `json:"tdruser"`
	} `json:"data"`
	Error []graphql.Error `json:"errors"`
}

type findUsersResponse struct {
	Data struct {
		Out []*UserOutput `json:"tdrusers"`
	} `json:"data"`
	Error []graphql.Error `json:"errors"`
}

// NewNotificationSvc takes a client from `client` package -- see examples/notifications.go for an example
func NewUserSvc(c *client.Client) IUserSvc {
	return &UserSvc{client: c}
}

func (u *UserSvc) GetUser(in *GetUserInput, rf graphql.ResponseFields, opts ...graphql.RequestOption) (*UserOutput, error) {
	query := fmt.Sprintf(`query ($id: ID!) {
			tdruser(id: $id) {
				%s
			}
		}
		`, rf)

	graphqlReq := graphql.NewRequest(query, opts...)
	graphqlReq.Var("id", in.ID)

	//TODO: Update GetUser to use graphql.ExecuteQuery(...)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(graphqlReq)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(buf.Bytes())
	userSvcURL := envy.Get(UserServiceEnv, DefaultURL)

	request, err := http.NewRequest(http.MethodPost, userSvcURL, reqBody)
	if err != nil {
		return nil, err
	}
	if _, ok := graphqlReq.Header[common.AuthorizationHeader]; ok {
		request.Header.Add(common.AuthorizationHeader, graphqlReq.Header.Get(common.AuthorizationHeader))
	}
	if _, ok := graphqlReq.Header[common.XTenantContextHeader]; ok {
		request.Header.Add(common.XTenantContextHeader, graphqlReq.Header.Get(common.XTenantContextHeader))
	} else {
		return nil, errors.New("tenant context is required")
	}

	resp, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server responded with an error: %d, %s", resp.StatusCode, resp.Status))
	}

	out := getUserResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, errors.New(fmt.Sprintf("error decoding get user response: %s", err.Error()))
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

func (u *UserSvc) FindUsers(in *FindUsersInput, rf graphql.ResponseFields, opts ...graphql.RequestOption) ([]*UserOutput, error) {
	query := fmt.Sprintf(`query ($email: String, $role: String, $tenantID: ID, $status: String, $page: Int, $perPage: Int) {
			tdrusers(email: $email, role: $role, tenantID: $tenantID, status: $status, page: $page, perPage: $perPage) {
				%s
			}
		}
		`, rf)

	graphqlReq := graphql.NewRequest(query, opts...)
	if in.Email != "" {
		graphqlReq.Var("email", in.Email)
	}
	if in.Role != "" {
		graphqlReq.Var("role", in.Role)
	}
	if in.TenantID != "" {
		graphqlReq.Var("tenantID", in.TenantID)
	}
	if in.Status != "" {
		graphqlReq.Var("status", in.Status)
	}
	if in.Page != 0 {
		graphqlReq.Var("page", in.Page)
	}
	if in.PerPage != 0 {
		graphqlReq.Var("perPage", in.PerPage)
	}

	//TODO: Update FindUsers to use graphql.ExecuteQuery(...)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(graphqlReq)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(buf.Bytes())
	userSvcURL := envy.Get(UserServiceEnv, DefaultURL)

	request, err := http.NewRequest(http.MethodPost, userSvcURL, reqBody)
	if err != nil {
		return nil, err
	}
	if _, ok := graphqlReq.Header[common.AuthorizationHeader]; ok {
		request.Header.Add(common.AuthorizationHeader, graphqlReq.Header.Get(common.AuthorizationHeader))
	}
	if _, ok := graphqlReq.Header[common.XTenantContextHeader]; ok {
		request.Header.Add(common.XTenantContextHeader, graphqlReq.Header.Get(common.XTenantContextHeader))
	} else {
		return nil, errors.New("tenant context is required")
	}

	resp, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server responded with an error: %d, %s", resp.StatusCode, resp.Status))
	}

	out := findUsersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, errors.New(fmt.Sprintf("error decoding find users response: %s", err.Error()))
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
