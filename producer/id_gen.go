// taken from https://github.com/open-telemetry/opentelemetry-go/blob/sdk/v1.30.0/sdk/trace/id_generator.go#L17
package producer

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	for {
		_, _ = gen.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	gen.Lock()
	defer gen.Unlock()
	tid := trace.TraceID{}
	sid := trace.SpanID{}
	for {
		_, _ = gen.randSource.Read(tid[:])
		if tid.IsValid() {
			break
		}
	}
	for {
		_, _ = gen.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return tid, sid
}

func defaultIDGenerator() *randomIDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}
