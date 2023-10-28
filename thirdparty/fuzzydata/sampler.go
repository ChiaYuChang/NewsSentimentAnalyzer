package main

import (
	"math/rand"
	"sort"
)

type Sampler[T any] struct {
	x []T
	p []float64
	c []float64
}

func NewSampler[T any](x []T, weight []float64) Sampler[T] {
	if weight == nil {
		weight = make([]float64, len(x))
		for i := range weight {
			weight[i] = 1.0 / float64(len(x))
		}
	}
	c := make([]float64, len(x))
	for i, w := range weight {
		c[i] = w
		if i > 0 {
			c[i] += c[i-1]
		}
	}

	return Sampler[T]{x, weight, c}
}

func (s Sampler[T]) Get() T {
	r := rand.Float64()
	return s.x[sort.SearchFloat64s(s.c, r)]
}

func (s Sampler[T]) GetN(n int) []T {
	rs := make([]T, n)
	for i := 0; i < n; i++ {
		rs[i] = s.Get()
	}
	return rs
}
