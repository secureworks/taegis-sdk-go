package client

import (
	"net/http"
	"os"
	"time"

	"github.com/secureworks/tdr-sdk-go/common"
	"github.com/secureworks/tdr-sdk-go/log"

	"github.com/hashicorp/go-cleanhttp"
	"moul.io/http2curl"
)

const (
	defaultTimeout = 5000 * time.Millisecond
)

// Client is the default wrapper around an http client
// extra helpers for our APIs around using bearer tokens
type Client struct {
	HTTPTimeout time.Duration
	CommandName string
	client      *http.Client
	header      http.Header
	Logger      log.Logger
	bearer      *string
	tenant      *string
}

// Do will run the HTTP request and add a bearer if the client was setup with a token
// optionally you can set a env var of `CURL_DEBUG=true` to get curl output printed out for
// easier debugging
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.addHeaderValues(req.Header)
	if os.Getenv("CURL_DEBUG") == "true" {
		command, err := http2curl.GetCurlCommand(req)
		if err != nil {
			c.Logger.Warn().WithError(err).Msg("curl not found")
		} else {
			c.Logger.Debug().WithStr("command", command.String()).Msg("curl command")
		}
	}
	return c.client.Do(req)
}

func (c *Client) addHeaderValues(header http.Header) {
	_, ok := header[common.AuthorizationHeader] //Only add the bearer token if there isn't one there already
	if c.bearer != nil && !ok {                 //APIs may not know how to respond if there is more than one entry - there are multiple ways to set one.
		header.Add(common.AuthorizationHeader, "Bearer "+*c.bearer)
	}
	_, ok = header[common.XTenantContextHeader]
	if c.tenant != nil && !ok {
		header.Add(common.XTenantContextHeader, *c.tenant)
	}
	header.Add("Content-Type", "application/json")
	if c.header != nil {
		for k, v := range c.header {
			if _, ok := header[k]; ok { //Only add the header if it isn't already set
				continue
			}
			for _, vn := range v {
				header.Add(k, vn)
			}
		}
	}
}

func (c *Client) Header() http.Header {
	header := http.Header{}
	c.addHeaderValues(header)
	return header
}

// NewClient creates a new default client with no options set
// This is typically just passed into a service call like
// 	notifications.NewNotificationSvc(client.NewClient())
func NewClient(opts ...Option) *Client {
	client := Client{
		HTTPTimeout: defaultTimeout,
		CommandName: "default-ctpx-sdk-go",
		Logger:      log.Noop(),
	}

	for _, opt := range opts {
		opt(&client)
	}

	if client.client == nil {
		client.client = cleanhttp.DefaultClient()
		client.client.Timeout = client.HTTPTimeout
	} else if client.client.Timeout == 0 {
		client.client.Timeout = client.HTTPTimeout
	}

	return &client
}

// Option is used for setting client options for the SDK
type Option func(*Client)

// WithLogger sets the underlying logger to be used with any Client passing through the returned Option
func WithLogger(in log.Logger) Option {
	return func(c *Client) {
		c.Logger = in
	}
}

// WithBearerToken lets you set a token to include with every request
func WithBearerToken(in string) Option {
	return func(c *Client) {
		c.bearer = &in
	}
}

// WithTenant lets you set a tenant id with every request
func WithTenant(in string) Option {
	return func(c *Client) {
		c.tenant = &in
	}
}

// WithHeader lets you add additional header values to all requests
func WithHeader(header http.Header) Option {
	return func(c *Client) {
		c.header = header
	}
}

// WithHTTPTimeout sets how long the request has to finish, this defaults to 5 seconds
func WithHTTPTimeout(in time.Duration) Option {
	return func(c *Client) {
		c.HTTPTimeout = in
	}
}

// WithHTTPClient sets the underlying http client for use with requests, overrides default of http.DefaultClient
func WithHTTPClient(hc *http.Client) Option {
	if hc == nil {
		hc = cleanhttp.DefaultClient()
	}
	return func(c *Client) {
		c.client = hc
	}
}
