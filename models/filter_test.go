package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Passes(t *testing.T) {
	t.Run("correct timestamps pass the filter", func(t *testing.T) {
		definition := FilterDefinition{
			Granularity: 2 * time.Millisecond,
			Names:       []string{"bla", "blup"},
		}
		filter := NewFilterCollection(definition)
		assert.False(t, filter.Passes("foo", &Numerical{Ts: 1000000}))
		assert.True(t, filter.Passes("bla", &Numerical{Ts: 1000000}))
		assert.False(t, filter.Passes("bla", &Numerical{Ts: 2000000}))
		assert.True(t, filter.Passes("bla", &Numerical{Ts: 3000000}))
		assert.False(t, filter.Passes("bla", &Numerical{Ts: 4000000}))
	})
}

func Test_TimestampFilter_Passes(t *testing.T) {
	t.Run("correct timestamps pass the filter", func(t *testing.T) {
		filter := &TimestampFilter{Granularity: 2 * time.Millisecond}
		assert.True(t, filter.Passes(&Numerical{Ts: 1000000}))
		assert.False(t, filter.Passes(&Numerical{Ts: 2000000}))
		assert.True(t, filter.Passes(&Numerical{Ts: 3000000}))
		assert.False(t, filter.Passes(&Numerical{Ts: 4000000}))
	})
}
