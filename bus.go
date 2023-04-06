package protobus

import (
	"context"
	"errors"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// HandlerFunc defines a handler function interface.
type HandlerFunc func(ctx context.Context, msg Msg) error
type HandlerFuncT[T Msg] func(ctx context.Context, msg T) error

// Msg defines a message interface.
type Msg = proto.Message
type eventKey = protoreflect.FullName

// Bus type defines the bus interface structure.
type Bus interface {
	Publish(ctx context.Context, msg Msg) error
	AddEventListener(handler HandlerFunc)
}

// InProcBus defines the bus structure.
type InProcBus struct {
	listeners map[eventKey][]HandlerFunc
}

func ProvideBus() *InProcBus {
	return &InProcBus{
		listeners: make(map[eventKey][]HandlerFunc),
	}
}

// Publish function publish a message to the bus listener.
func (b *InProcBus) Publish(ctx context.Context, msg Msg) error {
	var msgName = msg.ProtoReflect().Descriptor().FullName()
	var errs []error

	if listeners, exists := b.listeners[msgName]; exists {
		for _, listenerHandler := range listeners {
			if err := listenerHandler(ctx, msg); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}

func AddEventListener[T Msg](b *InProcBus, handler HandlerFuncT[T]) {
	eventName := (*new(T)).ProtoReflect().Descriptor().FullName()
	_, exists := b.listeners[eventName]
	if !exists {
		b.listeners[eventName] = make([]HandlerFunc, 0)
	}
	b.listeners[eventName] = append(b.listeners[eventName], func(ctx context.Context, msg Msg) error {
		return handler(ctx, msg.(T))
	})
}
