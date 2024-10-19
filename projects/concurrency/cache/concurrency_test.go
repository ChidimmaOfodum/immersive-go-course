package cache

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestConcurrency(t *testing.T) {
	input := map[string]any{
		"name":  "John",
		"age":   30,
		"hobby": "swimming",
	}

	t.Run("concurrent reads", func(t *testing.T) {
		c := NewCache[string, any](3)

		for key, value := range input {
			c.Put(key, value)
		}

		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				c.Get("hobby")
				wg.Done()
			}()
		}
		wg.Wait()

		got := c.GetStatistics()
		expect := Statistics{
			hitRate:             100.0,
			unReadEntries:       2,
			averageReads:        33.33,
			totalReadsAndWrites: 103,
		}
		require.Equal(t, expect, got)
	})

	t.Run("concurrent writes", func(t *testing.T) {
		c := NewCache[string, any](3)

		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				value := fmt.Sprintf("Pilot%d", i)
				c.Put("occupation", value)
				wg.Done()
			}()
		}
		wg.Wait()

		got := c.GetStatistics()
		expect := Statistics{
			hitRate:             0,
			unReadEntries:       1,
			averageReads:        0,
			totalReadsAndWrites: 100,
		}
		require.Equal(t, expect, got)

	})

	t.Run("concurrent reads and writes", func(t *testing.T) {
		c := NewCache[string, string](3)
		c.Put("occupation", "Pilot")

		var wg sync.WaitGroup

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				value := fmt.Sprintf("Pilot%d", i)
				c.Put("occupation", value)
				wg.Done()
			}()

			wg.Add(1)
			go func() {
				value, ok := c.Get("occupation")
				if !ok {
					t.Errorf("expected %v but got %v", true, ok)
				}
				if value == nil {
					t.Errorf("expect value to not be nil")
				} else {
					stringifiedValue := *value
					got := strings.Split(stringifiedValue, "")[0:5]
					if strings.Join(got, "") != "Pilot" {
						t.Errorf("expected %s but got %s", "Pilot", got)
					}
				}
				wg.Done()
			}()
		}

		wg.Wait()

		got := c.GetStatistics()
		expect := Statistics{
			hitRate:             100.0,
			unReadEntries:       0,
			averageReads:        100,
			totalReadsAndWrites: 201,
		}
		require.Equal(t, expect, got)

	})
}
