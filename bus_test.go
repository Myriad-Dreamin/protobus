package protobus

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"
)

type testQuery = emptypb.Empty

func TestEventPublish(t *testing.T) {
	bus := ProvideBus()

	var invoked bool

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		invoked = true
		return nil
	})

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}
	if !invoked {
		t.Fatal("event listener not invoked")
	}
}

func TestEventPublish_NoRegisteredListener(t *testing.T) {
	bus := ProvideBus()

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}
}

func TestEventCtxPublishCtx(t *testing.T) {
	bus := ProvideBus()

	var invoked bool

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		invoked = true
		return nil
	})

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}

	if !invoked {
		t.Fatal("event listener not invoked")
	}
}

func TestEventPublishCtx_NoRegisteredListener(t *testing.T) {
	bus := ProvideBus()

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}
}

func TestEventPublishCtx(t *testing.T) {
	bus := ProvideBus()

	var invoked bool

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		invoked = true
		return nil
	})

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}

	if !invoked {
		t.Fatal("event listener not invoked")
	}
}

func TestEventCtxPublish(t *testing.T) {
	bus := ProvideBus()

	var invoked bool

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		invoked = true
		return nil
	})

	err := bus.Publish(context.Background(), &testQuery{})
	if err != nil {
		t.Fatalf("unable to publish event: %v", err)
	}

	if !invoked {
		t.Fatal("event listener not invoked")
	}
}

func TestErrorPropagate(t *testing.T) {
	bus := ProvideBus()

	var internalErr = errors.New("test error")

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		return internalErr
	})

	err := bus.Publish(context.Background(), &testQuery{})
	if err == nil {
		t.Fatalf("unable to publish event: %v", err)
	}
	if !errors.Is(err, internalErr) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func BenchmarkProtoEventCtxPublish(b *testing.B) {
	bus := ProvideBus()

	AddEventListener(bus, func(ctx context.Context, query *testQuery) error {
		return nil
	})

	for i := 0; i < b.N; i++ {
		bus.Publish(context.Background(), &testQuery{})
	}
}

func ExampleAddEventListener() {
	bus := ProvideBus()
	var ErrEmpty = errors.New("empty")

	AddEventListener(bus, func(ctx context.Context, query *emptypb.Empty) error {
		if query == nil {
			return ErrEmpty
		}
		fmt.Printf("event received: %T", query)
		return nil
	})
	if err := bus.Publish(context.Background(), (*emptypb.Empty)(nil)); err != nil {
		fmt.Println("error received:", err)
		fmt.Println("error is ErrEmpty?", errors.Is(err, ErrEmpty))
	}
	bus.Publish(context.Background(), &emptypb.Empty{})
	// Output:
	// error received: empty
	// error is ErrEmpty? true
	// event received: *emptypb.Empty
}

func ExampleInProcBus_Publish() {
	bus := ProvideBus()

	AddEventListener(bus, func(ctx context.Context, query *emptypb.Empty) error {
		fmt.Printf("event received: %T", query)
		return nil
	})

	bus.Publish(context.Background(), &emptypb.Empty{})
	// Output: event received: *emptypb.Empty
}
