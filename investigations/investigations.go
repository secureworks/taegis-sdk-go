package investigations

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secureworks/tdr-sdk-go/client"
	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/graphql"
	"github.com/gobuffalo/envy"
	"github.com/hashicorp/go-multierror"
)

var (
	_             IInvestigationSvc      = &InvestigationSvc{}
	DefaultURL                           = "https://api.ctpx.secureworks.com/graphql"
	DefaultFields graphql.ResponseFields = `
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
		type
		`
)

// IInvestigationSvc defines what the the Investigation API can do
type IInvestigationSvc interface {
	GetInvestigation(*GetInvestigationInput, graphql.ResponseFields, ...graphql.RequestOption) (*InvestigationOutput, error)
}

// InvestigationsSvc is the concrete implementation of the interface against the real api
type InvestigationSvc struct {
	client      *client.Client
	serviceName string
}

// NewInvestigationsSvc takes a client from `client` package -- see examples/notifications.go for an example
func NewInvestigationSvc(c *client.Client, serviceName string) *InvestigationSvc {
	return &InvestigationSvc{client: c, serviceName: serviceName}
}

func (t *InvestigationSvc) GetInvestigation(in *GetInvestigationInput, rf graphql.ResponseFields, opts ...graphql.RequestOption) (*InvestigationOutput, error) {
	query := fmt.Sprintf(`query getInvestigation($id: ID!){
			investigation(investigation_id: $id) {
				%s
			}
		}
		`, rf)

	graphqlReq := graphql.NewRequest(query, opts...)
	graphqlReq.Var("id", in.ID)

	//TODO: Update GetInvestigation to use graphql.ExecuteQuery(...)
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(graphqlReq)
	if err != nil {
		return nil, err
	}

	reqBody := bytes.NewReader(buf.Bytes())
	investigationsURL := envy.Get("INVESTIGATIONS_URL", DefaultURL)

	request, err := http.NewRequest(http.MethodPost, investigationsURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("malformed request error: %w", err)
	}
	if _, ok := graphqlReq.Header[common.AuthorizationHeader]; ok {
		request.Header.Add(common.AuthorizationHeader, graphqlReq.Header.Get(common.AuthorizationHeader))
	}
	request.Header.Add(common.XTenantContextHeader, in.TenantID)

	resp, err := t.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("server connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("server responded with an error: %d", resp.StatusCode)
	}

	type createResponse struct {
		Data struct {
			Out *InvestigationOutput `json:"investigation"`
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

type GetInvestigationInput struct {
	TenantID string
	ID       string
}

type InvestigationOutput struct {
	ID            string                `json:"id"`
	CreatedAt     string                `json:"created_at"`
	CreatedBy     string                `json:"created_by"`
	UpdatedAt     string                `json:"updated_at"`
	TenantID      string                `json:"tenant_id"`
	Description   string                `json:"description"`
	Status        string                `json:"status"`
	KeyFindings   string                `json:"key_findings"`
	AssigneeID    string                `json:"assignee_id"`
	GenesisAlerts []GenesisAlertsOutput `json:"genesis_alerts"`
	GenesisEvents []GenesisEventsOutput `json:"genesis_events"`
	Alerts        []AlertsOutput        `json:"alerts"`
	Events        []EventsOutput        `json:"events"`
	Priority      int                   `json:"priority"`
	Type          string                `json:"type"`
}

type GenesisAlertsOutput struct {
	ID string `json:"id"`
}

type GenesisEventsOutput struct {
	ID string `json:"id"`
}

type AlertsOutput struct {
	ID string `json:"id"`
}

type EventsOutput struct {
	ID string `json:"id"`
}
