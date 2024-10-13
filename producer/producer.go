package producer

import (
	"context"
	"time"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Producer interface {
	ProduceOnce()
	Produce()
	Work()

	setCooldown(time.Duration)
	setErrorChan(chan error)
}

func New[Out any](name string, fn func() ([]Out, error), opts ...Option) (Producer, chan job.Job[Out], chan error) {
	p := &producer[Out]{
		name:  name,
		fn:    fn,
		out:   make(chan job.Job[Out]),
		err:   make(chan error),
		idgen: defaultIDGenerator(),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, p.out, p.err
}

type producer[Out any] struct {
	name string

	fn  func() ([]Out, error)
	out chan job.Job[Out]
	err chan error

	cooldown time.Duration

	idgen sdktrace.IDGenerator
}

func (p producer[Out]) ProduceOnce() {
	_, span := otel.Tracer(p.name).Start(context.Background(), "produce job")
	out, err := p.fn()
	span.End()
	if err != nil {
		p.err <- err
	}
	jobs := make([]job.Job[Out], len(out))
	for i, o := range out {
		j := job.New(context.Background(), o)
		j.StartSpan(otel.Tracer(p.name), "queue wait")
		jobs[i] = j
	}

	for _, j := range jobs {
		p.out <- j
		j.EndSpan()
	}
}

func (p producer[Out]) Produce() {
	for {
		p.ProduceOnce()
		time.Sleep(p.cooldown)
	}
}

func (p producer[Out]) Work() {
	p.Produce()
}

func (p *producer[Out]) setCooldown(t time.Duration) {
	p.cooldown = t
}

func (p *producer[Out]) setErrorChan(err chan error) {
	p.err = err
}
