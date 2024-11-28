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
	out, err := p.fn()
	if err != nil {
		_, errJob := job.New(context.Background(), err)
		p.err <- errJob
		return
	}

	jobs := make([]job.Job[Out], len(out))
	spans := make([]trace.Span, len(out))
	for i, o := range out {
		ctx := context.Background()
		if co, ok := any(o).(ContextualOutput); ok {
			ctx = co.Context()
		}
		rootCtx, j := job.New(ctx, o)
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

func (p *producer[Out]) setErrorQueue(err job.Queue[error]) {
	p.err = err
}

type ContextualOutput interface {
	Context() context.Context
}
