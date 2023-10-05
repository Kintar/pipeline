package pipeline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
)

func produceInts(limit int) Source[int] {
	current := 0
	return func() (result int, more bool, err error) {
		result = current
		current++
		more = current < limit
		return
	}
}

func TestProduceInts(t *testing.T) {
	s := produceInts(10)
	for i := 0; i < 9; i++ {
		r, m, e := s()
		assert.Equal(t, i, r)
		assert.True(t, m)
		assert.Nil(t, e)
	}
	r, m, e := s()
	assert.Equal(t, 9, r)
	assert.False(t, m)
	assert.Nil(t, e)
}

func TestSourceProducer(t *testing.T) {
	s := produceInts(10)
	p := sourceProducer(s)
	c, j := p(context.Background(), 0)
	assert.NotNil(t, c)
	assert.NotNil(t, j)
	var g errgroup.Group
	g.Go(j)
	count := 0
	for _ = range c {
		count++
	}
	assert.Equal(t, 10, count)
}
