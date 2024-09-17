package work

import (
	"github.com/imperfect-fourth/work/consumer"
	"github.com/imperfect-fourth/work/creator"
	"github.com/imperfect-fourth/work/transformer"
)

func NewCreator[Out any](fn func() ([]Out, error), opts ...creator.Option) (creator.Creator, chan Out, chan error) {
	return creator.New(fn, opts...)
}

func NewTransformer[In any, Out any](in chan In, fn func(In) (Out, error), opts ...transformer.Option) (transformer.Transformer, chan Out, chan error) {
	return transformer.New(in, fn, opts...)
}

func NewConsumer[In any](in chan In, fn func(In) error, opts ...consumer.Option) (consumer.Consumer, chan error) {
	return consumer.New(in, fn, opts...)
}
