package eventbus

import "context"

type Subscriber[T any] func(ctx context.Context, message T)

type Broker[T any] struct {
	subscribers []Subscriber[T]
}

func NewBroker[T any]() *Broker[T] {
	return &Broker[T]{
		subscribers: nil,
	}
}

func (b *Broker[T]) Subscribe(subscriber Subscriber[T]) {
	b.subscribers = append(b.subscribers, subscriber)
}

func (b *Broker[T]) Publish(ctx context.Context, message T) {
	for _, subscriber := range b.subscribers {
		subscriber(ctx, message)
	}
}
