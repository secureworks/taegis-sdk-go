package graphql

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/secureworks/taegis-sdk-go/log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	defaultBufferSize     = 1024
	disconnectionDeadline = 5 * time.Second
)

var (
	errSubscriptionFailedSendingMessage = errors.New("subscription failed sending message on channel")
	errSubscriptionCompleted            = errors.New("subscription ended")
	errOmitMessage                      = errors.New("omit message")
)

type Subscription struct {
	ch              chan *Message
	conn            *websocket.Conn
	u               *url.URL
	header          http.Header
	bufferSize      int
	sendKAMsgs      bool
	log             log.Logger
	responseCreator func() interface{}
	query           string
	vars            map[string]interface{}

	//indicates that the reader goroutine is done
	readerDone chan struct{}
}

type Message struct {
	Payload interface{}
	Err     error
}

type KeepAliveMessage struct{}

type SubscriptionOption func(s *Subscription)

func SubscriptionSendKAMessages(s *Subscription) {
	s.sendKAMsgs = true
}

func SubscriptionWithLog(l log.Logger) SubscriptionOption {
	return func(s *Subscription) {
		s.log = l
	}
}

func SubscriptionWithVars(vars map[string]interface{}) SubscriptionOption {
	return func(s *Subscription) {
		s.vars = vars
	}
}

func SubscriptionWithTenant(tenantID string) SubscriptionOption {
	return func(s *Subscription) {
		c := http.Cookie{
			Name:    "x-tenant-context",
			Value:   tenantID,
			Expires: time.Now().Add(1 * time.Hour),
		}
		s.header.Add("Cookie", c.String())
	}
}

func SubscriptionWithToken(token string) SubscriptionOption {
	return func(s *Subscription) {
		c := http.Cookie{
			Name:    "access_token",
			Value:   token,
			Expires: time.Now().Add(1 * time.Hour),
		}
		s.header.Add("Cookie", c.String())
	}
}

func SubscriptionWithHeader(header http.Header) SubscriptionOption {
	return func(s *Subscription) {
		s.header = header
	}
}

func NewSubscription(ctx context.Context, u *url.URL, query string, responseCreator func() interface{}, opts ...SubscriptionOption) (*Subscription, error) {
	if !strings.HasPrefix(strings.TrimSpace(query), "subscription") {
		return nil, errors.New("query must be a subscription")
	}

	s := &Subscription{u: u,
		bufferSize:      defaultBufferSize,
		log:             log.Noop(),
		header:          http.Header{},
		responseCreator: responseCreator,
		query:           query,
		readerDone:      make(chan struct{}, 1),
	}
	for _, opt := range opts {
		opt(s)
	}

	var (
		err  error
		resp *http.Response
	)
	s.conn, resp, err = websocket.DefaultDialer.DialContext(ctx, s.u.String(), s.header)
	if err != nil {
		s.conn = nil
		if resp != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("websocket error with status: %d", resp.StatusCode))
		}
		return nil, errors.Wrap(err, "websocket error")
	}

	defer func() {
		if err != nil {
			_ = s.Shutdown(ctx)
		}
	}()

	s.ch = make(chan *Message, s.bufferSize)

	go s.startWSReader()
	if err = s.connect(ctx); err != nil {
		return nil, err
	}

	if err != nil {
		s.log.WithError(err).WithFields(map[string]interface{}{
			"url":   s.u.String(),
			"query": s.query,
			"vars":  s.vars,
		}).Error().Msg("failed connecting to sub")
		return nil, err
	}

	return s, nil
}

func (s *Subscription) Messages() <-chan *Message {
	return s.ch
}

func (s *Subscription) informMessageReceived(m *Message) error {
	select {
	case s.ch <- m:
		return nil
	default:
		return errSubscriptionFailedSendingMessage
	}
}

func (s *Subscription) Shutdown(ctx context.Context) error {
	if s.conn == nil {
		return errors.New("subscription is already down")
	}
	s.log.Debug().Msg("sub close called")

	s.disconnect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.readerDone:
	}
	s.conn = nil
	s.log.Debug().Msg("sub goroutines done")
	s.log.Info().Msg("sub closed")
	return nil
}

func (s *Subscription) disconnect() {
	defer func() {
		_ = s.conn.Close()
	}()
	err := s.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(disconnectionDeadline))
	if err != nil {
		s.log.WithError(err).Warn().Msg("failed sending close message")
	}
}
