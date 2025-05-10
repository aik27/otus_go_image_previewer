package lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

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
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("Capacity exceeded", func(t *testing.T) {
		c := NewCache(5)
		for i := 0; i < 10; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		for i := 0; i < 5; i++ {
			_, ok := c.Get(Key(strconv.Itoa(i)))
			assert.False(t, ok, "Element with key \"%d\" still exist.", i)
		}
	})

	t.Run("Evict less used elements", func(t *testing.T) {
		c := NewCache(5)
		for i := 0; i < 5; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		_, ok := c.Get("0")
		assert.True(t, ok, "Element with key \"0\" doesn't exist.")

		_, ok = c.Get("1")
		assert.True(t, ok, "Element with key \"1\" doesn't exist.")

		c.Set("2", "new value")
		c.Set("6", 6)
		c.Set("7", 7)

		_, ok = c.Get("3")
		assert.False(t, ok, "Element with key \"3\" still exist.")

		_, ok = c.Get("4")
		assert.False(t, ok, "Element with key \"4\" still exist.")
	})

	t.Run("Clear cache", func(t *testing.T) {
		c := NewCache(5)
		for i := 0; i < 5; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		c.Clear()

		for i := 0; i < 5; i++ {
			_, ok := c.Get(Key(strconv.Itoa(i)))
			assert.False(t, ok, "Element with key \"%d\" still exist.", i)
		}
	})
}

func TestCacheMultithreading(t *testing.T) { //nolint
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000))))
		}
	}()

	wg.Wait()
}
