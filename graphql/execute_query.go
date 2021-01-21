package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/secureworks/tdr-sdk-go/log"

	"github.com/secureworks/tdr-sdk-go/common"

	"github.com/hashicorp/go-multierror"
)

//HTTPClient is any client that can perform HTTP requests.
//It is often but not always github.com/secureworks/tdr-sdk-go/client
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
	Header() http.Header
}

//QueryConfig takes all necessary or optional params for this package's Execute functions.
type QueryConfig struct {
	ServerURL  string
	HClient    HTTPClient
	Request    *Request
	Header     http.Header //DEPRECATED - DO NOT USE. Use graphql.RequestOptions instead
	EscapeHTML bool
	LimitRead  int64
	Output     interface{}
	logger     log.Logger
}

func (qc *QueryConfig) isValid() bool {
	v := qc == nil || qc.HClient == nil || qc.ServerURL == "" || qc.Request == nil
	return !v
}

//ExecuteQueryContext takes a context for HTTP request control and the given QueryConfig.
//It executes the graphql request against the proveded HTTPClient, returning an error or unmarshalling into QueryConfig.Output if provided.
//HClient, ServerURL, and Request are required in QueryConfig.
func ExecuteQueryContext(ctx context.Context, qc *QueryConfig) error {
	_, err := executeQueryContext(ctx, qc, false)
	return err
}

func ExecuteQueryContextWithTenant(ctx context.Context, qc *QueryConfig) (string, error) {
	return executeQueryContext(ctx, qc, true)
}

func executeQueryContext(ctx context.Context, qc *QueryConfig, enforceTenant bool) (string, error) {
	if ctx == nil || !qc.isValid() {
		return "", errors.New("ctpx-sdk-go/graphql: nil ctx or config to ExecuteQueryContext")
	}

	buf := bytes.NewBuffer(make([]byte, 0, 256))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(qc.EscapeHTML)

	err := enc.Encode(qc.Request)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, qc.ServerURL, buf)
	if err != nil {
		return "", err
	}

	for k := range qc.Request.Header { //RequestOption headers
		request.Header.Add(k, qc.Request.Header.Get(k))
	}

	for k := range qc.Header { //DEPRECATED - Rules client is still using this. It WILL be removed once Rules has been updated
		if _, ok := qc.Request.Header[k]; !ok { //If it already exists don't add it
			request.Header.Add(k, qc.Header.Get(k))
		}
	}

	tenant := request.Header.Get(common.XTenantContextHeader)
	if enforceTenant && tenant == "" {
		//check if client has a tenant defined
		tenant = qc.HClient.Header().Get(common.XTenantContextHeader)
		if tenant == "" {
			return "", errors.New("ctpx-sdk-go/graphql: request or client must specify tenant option")
		}
	}

	resp, err := qc.HClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("ctpx-sdk-go/graphql: server connection error: %w", err)
	}
	defer resp.Body.Close()
	var outErr error
	if resp.StatusCode >= http.StatusBadRequest {
		outErr = multierror.Append(outErr, fmt.Errorf("ctpx-sdk-go/graphql: server responded with an error: %d", resp.StatusCode))
	}

	graphqlResp := Response{
		Data: qc.Output,
	}

	var (
		r        io.Reader = resp.Body
		jsonBody bytes.Buffer
	)

	r = io.TeeReader(r, &jsonBody)

	if qc.LimitRead > 0 {
		r = io.LimitReader(r, qc.LimitRead)
	}

	if err := json.NewDecoder(r).Decode(&graphqlResp); err != nil {
		outErr = multierror.Append(outErr, fmt.Errorf("ctpx-sdk-go/graphql: error decoding response: %w", err))
		return "", outErr
	}

	for _, e := range graphqlResp.Error {
		outErr = multierror.Append(outErr, e)
	}

	if qc.logger != nil {
		qc.logger.Debug().WithError(outErr).WithFields(map[string]interface{}{
			"json": jsonBody.String(),
			"resp": graphqlResp,
		}).Msg("graphql resp")
	}

	return tenant, outErr
}

//ExecuteQuery is shorthand for ExecuteQueryContext with the given args.
//Look to ExecuteQueryContext for more details.
func ExecuteQuery(cli HTTPClient, serverURL string, graphqlReq *Request, out interface{}) error {
	return ExecuteQueryContext(context.Background(), &QueryConfig{
		HClient:   cli,
		ServerURL: serverURL,
		Request:   graphqlReq,
		Output:    out,
		logger:    graphqlReq.logger,
	})
}

func ExecuteQueryWithTenant(cli HTTPClient, serverURL string, graphqlReq *Request, out interface{}) (string, error) {
	return ExecuteQueryContextWithTenant(context.Background(), &QueryConfig{
		HClient:   cli,
		ServerURL: serverURL,
		Request:   graphqlReq,
		Output:    out,
		logger:    graphqlReq.logger,
	})
}

//ExecuteQueryEscapeHTML is shorthand for ExecuteQueryContext with the given args.
//Look to ExecuteQueryContext for more details.
func ExecuteQueryEscapeHTML(cli HTTPClient, serverURL string, graphqlReq *Request, escapeHTML bool, out interface{}) error {
	return ExecuteQueryContext(context.Background(), &QueryConfig{
		HClient:    cli,
		ServerURL:  serverURL,
		Request:    graphqlReq,
		EscapeHTML: escapeHTML,
		Output:     out,
		logger:     graphqlReq.logger,
	})
}
