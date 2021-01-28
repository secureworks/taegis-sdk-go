package mocks

import (
	"context"
	"errors"

	"github.com/secureworks/tdr-sdk-go/orchestration/playbooks"
)

type subscription struct {
	instances chan interface{}
}

func (s *subscription) Close() error {
	return nil
}
func (s *subscription) Next(ctx context.Context) (*playbooks.PlaybookInstance, error) {
	val, ok := <-s.instances
	if !ok {
		return nil, errors.New("instance channel closed")
	}

	switch v := val.(type) {
	case error:
		return nil, v
	case *playbooks.PlaybookInstance:
		return v, nil
	}
	return nil, errors.New("mock subscription only supports *PlaybooksInstance or error types")
}
