package producer

import (
	"context"
	"time"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
		name: name,
		fn:   fn,
		out:  make(chan job.Job[Out]),
		err:  make(chan error),
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
}

func (p producer[Out]) ProduceOnce() {
	_, span := otel.Tracer(p.name).Start(context.Background(), "produce job")
	out, err := p.fn()
	span.End()
	if err != nil {
		p.err <- err
	}

	jobs := make([]job.Job[Out], len(out))
	spans := make([]trace.Span, len(out))
	for i, o := range out {
		rootCtx, j := job.New(context.Background(), o)
		_, jobspan := otel.Tracer(p.name).Start(rootCtx, "start queue wait")
		jobs[i] = j
		spans[i] = jobspan
	}

	for i, j := range jobs {
		p.out <- j
		spans[i].End()
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
