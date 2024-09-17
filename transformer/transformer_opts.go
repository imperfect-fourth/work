package transformer

type TransformerOpt func(Transformer)

func WithParallelism(p int) TransformerOpt {
	return func(t Transformer) {
		t.withParallelism(p)
	}
}

func WithQueueSize(s int) TransformerOpt {
	return func(t Transformer) {
		t.withQueueSize(s)
	}
}
