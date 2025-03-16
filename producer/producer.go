package producer

import (
	"context"
	"time"

	"github.com/imperfect-fourth/work/job"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Producer[Out any] interface {
	ProduceOnce()
	Produce()
	Work()
	Output() job.Queue[Out]
	Error() job.Queue[error]

	WithCooldown(time.Duration) Producer[Out]
	WithErrorQueue(job.Queue[error]) Producer[Out]
}

func New[Out any](name string, fn func() ([]Out, error)) Producer[Out] {
	p := &producer[Out]{
		name: name,
		fn:   fn,
		out:  make(chan job.Job[Out]),
		err:  make(chan job.Job[error]),
	}
	return p
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

func (p producer[Out]) Output() job.Queue[Out] {
	return p.out
}

func (p producer[Out]) Error() job.Queue[error] {
	return p.err
}

func (p *producer[Out]) WithCooldown(t time.Duration) Producer[Out] {
	p.cooldown = t
	return p
}

func (p *producer[Out]) WithErrorQueue(err job.Queue[error]) Producer[Out] {
	p.err = err
	return p
}

type ContextualOutput interface {
	Context() context.Context
}
