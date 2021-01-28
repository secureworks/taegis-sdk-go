package graphql_test

import (
	"context"

	"github.com/99designs/gqlgen/example/chat"
)

type resolver struct{}

func (r *resolver) Mutation() chat.MutationResolver {
	return nil
}

func (r *resolver) Query() chat.QueryResolver {
	return nil
}

func (r *resolver) Subscription() chat.SubscriptionResolver {
	return &subscriptionResolver{r}
}

func newChatConfig() chat.Config {
	return chat.Config{
		Resolvers: &resolver{},
	}
}

type subscriptionResolver struct{ *resolver }

func (r *subscriptionResolver) MessageAdded(_ context.Context, _ string) (<-chan *chat.Message, error) {
	ch := make(chan *chat.Message, 1)
	ch <- &chat.Message{Text: "hello!"}
	return ch, nil
}
