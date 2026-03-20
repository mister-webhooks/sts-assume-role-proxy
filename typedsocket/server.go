package typedsocket

import (
	"context"
	"fmt"
	"log"
	"net"
)

type TypedServer[Q any, R any] struct {
	listener net.Listener
}

func NewTypedServer[Q any, R any](
	mkListener func() (net.Listener, error),
) (*TypedServer[Q, R], error) {
	listener, err := mkListener()

	if err != nil {
		return nil, err
	}

	return &TypedServer[Q, R]{
		listener: listener,
	}, nil
}

func (ts *TypedServer[Q, R]) Close() error {
	return ts.listener.Close()
}

func (ts *TypedServer[Q, R]) Serve(ctx context.Context, handler func(context.Context, *TypedConnection[R, Q]) error) error {
	for {
		conn, err := ts.listener.Accept()
		if err != nil {
			return fmt.Errorf("accept error: %w", err)
		}

		tc := NewTypedConnection[R, Q](conn)

		go func() {
			defer tc.Close()
			err := handler(ctx, tc)

			if err != nil {
				log.Printf("error running handler: %s", err)
			}
		}()
	}
}
