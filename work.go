package work

import (
	"github.com/imperfect-fourth/work/consumer"
	"github.com/imperfect-fourth/work/job"
	"github.com/imperfect-fourth/work/producer"
	"github.com/imperfect-fourth/work/transformer"
)

func NewProducer[Out any](name string, fn func() ([]Out, error), opts ...producer.Option) (producer.Producer, job.Queue[Out], chan error) {
	return producer.New(name, fn, opts...)
}

func NewTransformer[In any, Out any](name string, in job.Queue[In], fn func(In) (Out, error), opts ...transformer.Option) (transformer.Transformer, job.Queue[Out], chan error) {
	return transformer.New(name, in, fn, opts...)
}

func NewConsumer[In any](name string, in job.Queue[In], fn func(In) error, opts ...consumer.Option) (consumer.Consumer, chan error) {
	return consumer.New(name, in, fn, opts...)
}
