package producer

import (
	"context"
	"fmt"
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
	setErrorQueue(job.Queue[error])
}

func New[Out any](name string, fn func() ([]Out, error), opts ...Option) (Producer, job.Queue[Out], job.Queue[error]) {
	p := &producer[Out]{
		name: name,
		fn:   fn,
		out:  make(chan job.Job[Out]),
		err:  make(chan job.Job[error]),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, p.out, p.err
}

type producer[Out any] struct {
	name string

	fn  func() ([]Out, error)
	out job.Queue[Out]
	err job.Queue[error]

	cooldown time.Duration
}

func (p producer[Out]) ProduceOnce() {
	ctx, span := otel.Tracer(p.name).Start(context.Background(), "produce job")
	out, err := p.fn()
	span.End()
	if err != nil {
		_, errJob := job.New(ctx, err)
		p.err <- errJob
	}
	if len(out) > 0 {
		fmt.Println(span.SpanContext().TraceID())
	}

	jobs := make([]job.Job[Out], len(out))
	spans := make([]trace.Span, len(out))
	for i, o := range out {
		rootCtx, j := job.New(context.Background(), o)
		_, jobspan := otel.Tracer(p.name).Start(rootCtx, "start queue wait")
		fmt.Println(jobspan.SpanContext().TraceID())
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

func (p *producer[Out]) setErrorQueue(err job.Queue[error]) {
	p.err = err
}
