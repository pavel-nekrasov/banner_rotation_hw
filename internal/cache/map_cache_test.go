package cache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMapCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewMapCache[string, int]()

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewMapCache[string, int]()

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.False(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Equal(t, 0, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewMapCache[string, int]()

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		require.Equal(t, 3, c.Len())

		c.Clear()

		require.Equal(t, 0, c.Len())

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get("bbb")
		require.False(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Equal(t, 0, val)
	})

	t.Run("value rewrite logic", func(t *testing.T) {
		c := NewMapCache[string, int]()

		c.Set("aaa", 100)
		c.Set("bbb", 200)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		ok = c.Set("bbb", 400)
		require.False(t, ok)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 400, val)
	})
}

func TestMapCacheMultithreading(_ *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	c := NewMapCache[string, int]()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 1_000_000; i++ {
			c.Get(strconv.Itoa(rnd.Intn(1000)))
		}
	}()

	wg.Wait()
}
