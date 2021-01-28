package playbooks

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/secureworks/tdr-sdk-go/common"

	"github.com/secureworks/tdr-sdk-go/graphql"
)

const playbookIDs = "playbookIds"

type Subscription interface {
	Next(ctx context.Context) (*PlaybookInstance, error)
	io.Closer
}

type subscription struct {
	sub  *graphql.Subscription
	msgs <-chan *graphql.Message
}

func (s *subscription) Next(ctx context.Context) (*PlaybookInstance, error) {
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
			return &PlaybookInstance{}, nil
		case *playbookInstanceCreatedEvent:
			return t.PlaybookInstanceEvent, nil
		case *playbookInstanceUpdatedEvent:
			return t.PlaybookInstanceEvent, nil
		case *playbookInstanceDeletedEvent:
			return t.PlaybookInstanceEvent, nil
		default:
			return nil, fmt.Errorf("unexpected event payload type %T", msg.Payload)
		}
	}
}

func (s *subscription) Close() error {
	return s.sub.Shutdown(context.Background())
}

type playbookInstanceCreatedEvent struct {
	PlaybookInstanceEvent *PlaybookInstance `json:"playbookInstanceCreated"`
}
type playbookInstanceUpdatedEvent struct {
	PlaybookInstanceEvent *PlaybookInstance `json:"playbookInstanceUpdated"`
}
type playbookInstanceDeletedEvent struct {
	PlaybookInstanceEvent *PlaybookInstance `json:"playbookInstanceDeleted"`
}

func (playbookService *playbookSvc) createSub(ctx context.Context, playbooks common.IDs, query string, creator func() interface{}, options ...graphql.SubscriptionOption) (Subscription, error) {
	u, err := url.Parse(playbookService.url)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "http", "localhost", "ws":
		u.Scheme = "ws"
	case "https", "wss":
		u.Scheme = "wss"
	default:
		return nil, fmt.Errorf("invalid protocol: %s", u.Scheme)
	}

	gSub, err := graphql.NewSubscription(ctx, u,
		graphql.AddVarNamesToQuery(query, playbookIDs),
		creator,
		append(options, graphql.SubscriptionWithVars(map[string]interface{}{playbookIDs: playbooks}))...)
	if err != nil {
		return nil, err
	}
	sub := &subscription{sub: gSub, msgs: gSub.Messages()}
	return sub, nil
}

func (playbookService *playbookSvc) PlaybookInstanceCreated(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error) {
	return playbookService.createSub(ctx, playbooks, playbookInstanceCreateSubQuery, func() interface{} {
		return &playbookInstanceCreatedEvent{}
	}, options...)
}

func (playbookService *playbookSvc) PlaybookInstanceDeleted(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error) {
	return playbookService.createSub(ctx, playbooks, playbookInstanceDeleteSubQuery, func() interface{} {
		return &playbookInstanceDeletedEvent{}
	}, options...)
}

func (playbookService *playbookSvc) PlaybookInstanceUpdated(ctx context.Context, playbooks common.IDs, options ...graphql.SubscriptionOption) (Subscription, error) {
	return playbookService.createSub(ctx, playbooks, playbookInstanceUpdateSubQuery, func() interface{} {
		return &playbookInstanceUpdatedEvent{}
	}, options...)
}
