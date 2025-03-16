package work

import (
	"context"

	"github.com/imperfect-fourth/work/consumer"
	"github.com/imperfect-fourth/work/producer"
	"github.com/imperfect-fourth/work/transformer"
)

func NewProducer[Out any](name string, fn func() ([]Out, error)) producer.Producer[Out] {
	return producer.New(name, fn)
}

func NewTransformer[In any, Out any](name string, fn func(context.Context, In) (Out, error)) transformer.Transformer[In, Out] {
	return transformer.New(name, fn)
}

func NewConsumer[In any](name string, fn func(context.Context, In) error) consumer.Consumer[In] {
	return consumer.New(name, fn)
}
