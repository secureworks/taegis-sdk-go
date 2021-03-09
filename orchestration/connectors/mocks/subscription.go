package mocks

import (
	"context"
	"io"

	"github.com/secureworks/taegis-sdk-go/orchestration/connectors"
)

type subscription struct {
	err        error
	connectors chan *connectors.Connector
}

func (s *subscription) Close() error {
	close(s.connectors)
	return nil
}
func (s *subscription) Next(_ context.Context) (*connectors.Connector, error) {
	c, ok := <-s.connectors
	if ok {
		return c, nil
	}
	return nil, io.EOF
}
