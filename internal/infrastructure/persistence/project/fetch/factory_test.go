package fetch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchFactory(t *testing.T) {

	var factory *Factory
	var fetcher *LocalStorage

	factory = NewFactory()

	fetcher = NewLocalStorage(nil, nil)
	factory.Register("fetcher", fetcher)

	t.Run("Get fetcher", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get fetcher")
		f := factory.Get("fetcher")
		assert.NotNil(t, f)
		assert.Equal(t, fetcher, f)
	})

	t.Run("Get fetcher not registered", func(t *testing.T) {
		t.Parallel()
		t.Log("Testing get fetcher not registered")
		f := factory.Get("fetcher_not_registered")
		assert.Nil(t, f)
	})
}
