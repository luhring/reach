package cache

import "testing"

func TestCache(t *testing.T) {
	t.Run("when empty", func(t *testing.T) {
		c := New()

		t.Run("Get", func(t *testing.T) {
			t.Run("returns nil", func(t *testing.T) {
				if result := c.Get("some-key"); result != nil {
					t.Errorf("result: %v", result)
				}
			})
		})
	})

	t.Run("with several items", func(t *testing.T) {
		items := []struct {
			key   string
			value interface{}
		}{
			{
				key:   "key1",
				value: "thing1",
			},
			{
				key:   "key2",
				value: "thing2",
			},
		}

		c := New()
		for _, item := range items {
			c.Put(item.key, item.value)
		}

		t.Run("returns correct value for key", func(t *testing.T) {
			for _, item := range items {
				if result := c.Get(item.key); result != item.value {
					t.Errorf("item: %v, result: %v", item, result)
				}
			}
		})
	})
}
