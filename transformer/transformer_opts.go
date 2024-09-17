package transformer

type TransformerOpt func(Transformer)

func WithParallelism(p int) TransformerOpt {
	return func(t Transformer) {
		t.setParallelism(p)
	}
}

func WithQueueSize(s int) TransformerOpt {
	return func(t Transformer) {
		t.setQueueSize(s)
	}
}
