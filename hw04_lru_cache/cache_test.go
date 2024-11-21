package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

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

	t.Run("purge logic", func(_ *testing.T) {
		c := NewCache(3)
		c.Set("elem_1", 100)
		c.Set("elem_2", 200)
		c.Set("elem_3", 300)
		c.Set("elem_4", 400)

		_, existElem1 := c.Get("elem_1")
		require.Equal(t, false, existElem1, "failed simple push")

		c.Clear()
		_, existElem2 := c.Get("elem_2")
		_, existElem3 := c.Get("elem_3")
		_, existElem4 := c.Get("elem_4")
		require.Equal(t, false, existElem2 && existElem3 && existElem4, "failed clear cache")

		c.Set("elem_1", 100)
		c.Set("elem_2", 200)
		c.Set("elem_3", 300)
		c.Get("elem_1")
		c.Get("elem_1")
		c.Get("elem_2")
		c.Get("elem_2")
		c.Set("elem_1", 1100)
		c.Set("elem_4", 400)

		_, existElem3 = c.Get("elem_3")
		require.Equal(t, false, existElem3, "failed complex push: exist element")
		_, existElem1 = c.Get("elem_1")
		_, existElem2 = c.Get("elem_2")
		_, existElem4 = c.Get("elem_4")
		require.Equal(t, true, existElem1 && existElem2 && existElem4, "failed complex push: not exist element")
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
