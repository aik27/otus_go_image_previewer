package lrucache

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10, nil)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5, nil)

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
		c := NewCache(5, nil)
		for i := 0; i < 10; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		for i := 0; i < 5; i++ {
			_, ok := c.Get(Key(strconv.Itoa(i)))
			assert.False(t, ok, "Element with key \"%d\" still exist.", i)
		}
	})

	t.Run("Evict less used elements", func(t *testing.T) {
		c := NewCache(5, nil)
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
		c := NewCache(5, nil)
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
	c := NewCache(10, nil)
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
			c.Get(Key(strconv.Itoa(rand.Intn(1_000)))) //nolint:gosec
		}
	}()

	wg.Wait()
}

func TestCacheWithEvents(t *testing.T) {
	t.Run("OnEvictedEvent", func(t *testing.T) {
		c := NewCache(5, OnEvictedEvent)
		for i := 0; i < 10; i++ {
			savePath := fmt.Sprintf("%s/%d.txt", os.TempDir(), i)
			file, err := os.Create(savePath)
			assert.NoError(t, err, "File \"%s\" should be created", savePath)

			_, err = file.WriteString("Hello World")
			assert.NoError(t, err, "File \"%s\" should be written", savePath)

			cached := ImageItem{
				FilePath:    savePath,
				Width:       1,
				Height:      1,
				OriginalURL: "test",
			}
			c.Set(Key(strconv.Itoa(i)), cached)
		}

		for i := 0; i < 5; i++ {
			_, ok := c.Get(Key(strconv.Itoa(i)))
			assert.False(t, ok, "Element with key \"%d\" still exist.", i)

			savePath := fmt.Sprintf("%s/%d.jpg", os.TempDir(), i)
			_, err := os.Open(savePath)
			assert.Error(t, err, "File \"%s\" should be deleted", savePath)
		}
	})
}
