package connectors

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/secureworks/tdr-sdk-go/common"

	"github.com/secureworks/tdr-sdk-go/graphql"
)

const (
	methodsVarName    = "methods"
	allTenantsVarName = "all"
)

type Subscription interface {
	Next(ctx context.Context) (*Connector, error)
	io.Closer
}

type subscription struct {
	sub  *graphql.Subscription
	msgs <-chan *graphql.Message
}

func (s *subscription) Next(ctx context.Context) (*Connector, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg, ok := <-s.msgs:
		if !ok {
			return nil, io.EOF
		}
		if msg.Err != nil {
			return nil, msg.Err
		}
		switch t := msg.Payload.(type) {
		case graphql.KeepAliveMessage:
			return &Connector{}, nil
		case *connectorCreatedEvent:
			return t.ConnectorEvent, nil
		case *connectorUpdatedEvent:
			return t.ConnectorEvent, nil
		case *connectorDeletedEvent:
			return t.ConnectorEvent, nil
		default:
			return nil, fmt.Errorf("unexpected event payload type %T", msg.Payload)
		}
	}
}

func (s *subscription) Close() error {
	return s.sub.Shutdown(context.Background())
}

type connectorCreatedEvent struct {
	ConnectorEvent *Connector `json:"connectorCreated"`
}
type connectorUpdatedEvent struct {
	ConnectorEvent *Connector `json:"connectorUpdated"`
}
type connectorDeletedEvent struct {
	ConnectorEvent *Connector `json:"connectorDeleted"`
}

func (o *connectorSvc) createSub(ctx context.Context, connectorMethods common.IDs, allTenants bool, query string, creator func() interface{}, options ...graphql.SubscriptionOption) (Subscription, error) {
	u, err := url.Parse(o.url)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws":
		break
	case "wss":
		break
	default:
		return nil, fmt.Errorf("invalid protocol: %s", u.Scheme)
	}

	gSub, err := graphql.NewSubscription(ctx, u,
		graphql.AddVarNamesToQuery(query, methodsVarName, allTenantsVarName),
		creator,
		append(options, graphql.SubscriptionWithVars(
			map[string]interface{}{methodsVarName: connectorMethods, allTenantsVarName: allTenants},
		))...)
	if err != nil {
		return nil, err
	}

	sub := &subscription{sub: gSub, msgs: gSub.Messages()}
	return sub, nil
}

func (o *connectorSvc) ConnectorCreated(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error) {
	return o.createSub(ctx, connectorMethods, allTenants, connectorCreatedSubQuery, func() interface{} {
		return &connectorCreatedEvent{}
	}, options...)
}

func (o *connectorSvc) ConnectorDeleted(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error) {
	return o.createSub(ctx, connectorMethods, allTenants, connectorDeletedSubQuery, func() interface{} {
		return &connectorDeletedEvent{}
	}, options...)
}

func (o *connectorSvc) ConnectorUpdated(ctx context.Context, connectorMethods common.IDs, allTenants bool, options ...graphql.SubscriptionOption) (Subscription, error) {
	return o.createSub(ctx, connectorMethods, allTenants, connectorUpdatedSubQuery, func() interface{} {
		return &connectorUpdatedEvent{}
	}, options...)
}
